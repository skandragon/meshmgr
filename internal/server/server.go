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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skandragon/meshmgr/internal/config"
	"github.com/skandragon/meshmgr/meshdb"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	db     *pgxpool.Pool
	router *chi.Mux
}

// New creates a new Server instance
func New(cfg *config.Config) (*Server, error) {
	// Create database connection pool
	pool, err := pgxpool.New(context.Background(), cfg.Database.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	s := &Server{
		config: cfg,
		db:     pool,
		router: chi.NewRouter(),
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

// setupMiddleware configures middleware for the router
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", s.handleHealth)

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.handleRegister)
			r.Post("/login", s.handleLogin)
			r.Post("/logout", s.handleLogout)
			r.Get("/me", s.handleMe)
		})

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware)

			// Mesh routes
			r.Route("/meshes", func(r chi.Router) {
				r.Get("/", s.handleListMeshes)
				r.Post("/", s.handleCreateMesh)
				r.Route("/{meshID}", func(r chi.Router) {
					r.Get("/", s.handleGetMesh)
					r.Put("/", s.handleUpdateMesh)
					r.Delete("/", s.handleDeleteMesh)
				})
			})
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// Close closes the server and database connections
func (s *Server) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// DB returns the database queries interface
func (s *Server) DB() *meshdb.Queries {
	return meshdb.New(s.db)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
