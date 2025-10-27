// Copyright (C) 2025 Michael Graff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/internal/auth"
	"github.com/skandragon/meshmgr/meshdb"
)

type contextKey string

const userContextKey contextKey = "user"

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	User  *meshdb.User `json:"user"`
}

// handleRegister handles user registration
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "Email, password, and display name are required")
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password, s.config.Auth.BCryptCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	user, err := s.DB().CreateUser(r.Context(), meshdb.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			writeError(w, http.StatusConflict, "Email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.config.Auth.JWTSecret, s.config.Auth.JWTExpiration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Create session
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate session token")
		return
	}

	_, err = s.DB().CreateSession(r.Context(), meshdb.CreateSessionParams{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(s.config.Auth.JWTExpiration),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// handleLogin handles user login
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := s.DB().GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, s.config.Auth.JWTSecret, s.config.Auth.JWTExpiration)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Create session
	sessionToken, err := auth.GenerateRandomToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate session token")
		return
	}

	_, err = s.DB().CreateSession(r.Context(), meshdb.CreateSessionParams{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(s.config.Auth.JWTExpiration),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// handleLogout handles user logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	// Extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}

	token := parts[1]

	// Validate token to get user ID
	claims, err := auth.ValidateToken(token, s.config.Auth.JWTSecret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Delete all sessions for the user
	if err := s.DB().DeleteUserSessions(r.Context(), claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// handleMe returns the current user's information
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, "Missing authorization header")
		return
	}

	// Extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}

	token := parts[1]

	// Validate token
	claims, err := auth.ValidateToken(token, s.config.Auth.JWTSecret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Get user
	user, err := s.DB().GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// authMiddleware is a middleware that validates JWT tokens or API keys
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Try JWT validation first
		claims, err := auth.ValidateToken(token, s.config.Auth.JWTSecret)
		if err == nil {
			// JWT is valid, get user
			user, err := s.DB().GetUserByID(r.Context(), claims.UserID)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "User not found")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), userContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// JWT validation failed, try API key
		user, err := s.validateAPIKey(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Invalid token or API key")
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateAPIKey validates an API key and returns the associated user
func (s *Server) validateAPIKey(ctx context.Context, key string) (*meshdb.User, error) {
	// Hash the provided key using SHA256
	keyHash := auth.HashAPIKey(key)

	// Look up the API key by hash
	apiKey, err := s.DB().GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, pgx.ErrNoRows
	}

	// Update last used timestamp (async, ignore errors)
	go func() {
		_ = s.DB().UpdateAPIKeyLastUsed(context.Background(), apiKey.ID)
	}()

	// Get user
	user, err := s.DB().GetUserByID(ctx, apiKey.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// getUserFromContext retrieves the user from the request context
func getUserFromContext(ctx context.Context) *meshdb.User {
	user, ok := ctx.Value(userContextKey).(*meshdb.User)
	if !ok {
		return nil
	}
	return user
}
