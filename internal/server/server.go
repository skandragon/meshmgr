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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skandragon/meshmgr/internal/config"
	"github.com/skandragon/meshmgr/meshdb"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	db     *pgxpool.Pool
	mux    *http.ServeMux
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
		mux:    http.NewServeMux(),
	}

	s.setupRoutes()

	return s, nil
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Health check
	s.mux.HandleFunc("GET /health", s.handleHealth)

	// Auth routes (public)
	s.mux.HandleFunc("POST /api/auth/register", s.handleRegister)
	s.mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	s.mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	s.mux.HandleFunc("GET /api/auth/me", s.handleMe)

	// LoRa configuration (public)
	s.mux.HandleFunc("GET /api/lora-config", s.handleGetLoRaConfig)

	// Mesh routes (protected)
	s.mux.HandleFunc("GET /api/meshes", s.withAuth(s.handleListMeshes))
	s.mux.HandleFunc("POST /api/meshes", s.withAuth(s.handleCreateMesh))
	s.mux.HandleFunc("GET /api/meshes/{meshID}", s.withAuth(s.handleGetMesh))
	s.mux.HandleFunc("PUT /api/meshes/{meshID}", s.withAuth(s.handleUpdateMesh))
	s.mux.HandleFunc("DELETE /api/meshes/{meshID}", s.withAuth(s.handleDeleteMesh))

	// Admin keys routes (protected)
	s.mux.HandleFunc("GET /api/meshes/{meshID}/admin-keys", s.withAuth(s.handleListAdminKeys))
	s.mux.HandleFunc("POST /api/meshes/{meshID}/admin-keys", s.withAuth(s.handleCreateAdminKey))
	s.mux.HandleFunc("GET /api/meshes/{meshID}/admin-keys/{keyID}", s.withAuth(s.handleGetAdminKey))
	s.mux.HandleFunc("DELETE /api/meshes/{meshID}/admin-keys/{keyID}", s.withAuth(s.handleDeleteAdminKey))

	// Nodes routes (protected)
	s.mux.HandleFunc("GET /api/meshes/{meshID}/nodes", s.withAuth(s.handleListNodes))
	s.mux.HandleFunc("POST /api/meshes/{meshID}/nodes", s.withAuth(s.handleCreateNode))
	s.mux.HandleFunc("GET /api/meshes/{meshID}/nodes/{nodeID}", s.withAuth(s.handleGetNode))
	s.mux.HandleFunc("PUT /api/meshes/{meshID}/nodes/{nodeID}", s.withAuth(s.handleUpdateNode))
	s.mux.HandleFunc("PATCH /api/meshes/{meshID}/nodes/{nodeID}/status", s.withAuth(s.handleUpdateNodeStatus))
	s.mux.HandleFunc("DELETE /api/meshes/{meshID}/nodes/{nodeID}", s.withAuth(s.handleDeleteNode))
}

// withAuth wraps a handler with authentication middleware
func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.authMiddleware(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s", addr)

	// Apply middleware to the entire mux
	handler := chain(s.mux, corsMiddleware, loggingMiddleware, recovererMiddleware)

	return http.ListenAndServe(addr, handler)
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
