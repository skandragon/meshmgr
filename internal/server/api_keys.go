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
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/internal/auth"
	"github.com/skandragon/meshmgr/meshdb"
)

// CreateAPIKeyRequest represents a request to create an API key
type CreateAPIKeyRequest struct {
	KeyName   string `json:"key_name"`
	ExpiresIn *int64 `json:"expires_in,omitempty"` // seconds from now, or null for no expiration
}

// CreateAPIKeyResponse represents the response when creating an API key
type CreateAPIKeyResponse struct {
	APIKey string           `json:"api_key"` // Plain key, only shown once
	Key    meshdb.UserApiKey `json:"key"`     // Key metadata (no hash)
}

// handleListAPIKeys handles listing API keys for the current user
func (s *Server) handleListAPIKeys(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	keys, err := s.DB().ListAPIKeysByUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list API keys")
		return
	}

	writeJSON(w, http.StatusOK, keys)
}

// handleCreateAPIKey handles creating a new API key
func (s *Server) handleCreateAPIKey(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.KeyName == "" {
		writeError(w, http.StatusBadRequest, "Key name is required")
		return
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn != nil {
		expirationTime := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Second)
		expiresAt = &expirationTime
	}

	// Create API key in database with placeholder hash to get the ID
	key, err := s.DB().CreateAPIKey(r.Context(), meshdb.CreateAPIKeyParams{
		UserID:  user.ID,
		KeyHash: "placeholder", // Will be updated below
		KeyName: req.KeyName,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	// Generate API key with the ID
	plainKey, keyHash, err := auth.GenerateAPIKey(key.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate API key")
		return
	}

	// Update the key with the actual hash
	key, err = s.DB().UpdateAPIKeyHash(r.Context(), meshdb.UpdateAPIKeyHashParams{
		ID:      key.ID,
		KeyHash: keyHash,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	writeJSON(w, http.StatusCreated, CreateAPIKeyResponse{
		APIKey: plainKey,
		Key:    key,
	})
}

// handleDeleteAPIKey handles deleting an API key
func (s *Server) handleDeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	keyIDStr := r.PathValue("keyID")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	// Verify the key belongs to the user
	key, err := s.DB().GetAPIKey(r.Context(), keyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "API key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get API key")
		return
	}

	if key.UserID != user.ID {
		writeError(w, http.StatusNotFound, "API key not found")
		return
	}

	// Delete the key
	if err := s.DB().DeleteAPIKey(r.Context(), keyID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete API key")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "API key deleted successfully",
	})
}
