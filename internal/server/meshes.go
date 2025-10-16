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

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/meshdb"
)

// CreateMeshRequest represents a request to create a mesh
type CreateMeshRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateMeshRequest represents a request to update a mesh
type UpdateMeshRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// handleListMeshes handles listing meshes for the current user
func (s *Server) handleListMeshes(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshes, err := s.DB().ListMeshesByUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list meshes")
		return
	}

	writeJSON(w, http.StatusOK, meshes)
}

// handleCreateMesh handles creating a new mesh
func (s *Server) handleCreateMesh(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateMeshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	mesh, err := s.DB().CreateMesh(r.Context(), meshdb.CreateMeshParams{
		OwnerID:     user.ID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create mesh")
		return
	}

	writeJSON(w, http.StatusCreated, mesh)
}

// handleGetMesh handles getting a single mesh
func (s *Server) handleGetMesh(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := chi.URLParam(r, "meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	mesh, err := s.DB().GetMeshByID(r.Context(), meshID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get mesh")
		return
	}

	writeJSON(w, http.StatusOK, mesh)
}

// handleUpdateMesh handles updating a mesh
func (s *Server) handleUpdateMesh(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := chi.URLParam(r, "meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	var req UpdateMeshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing mesh to check ownership
	mesh, err := s.DB().GetMeshByID(r.Context(), meshID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get mesh")
		return
	}

	// Check if user owns the mesh (basic check - will need mesh_access later)
	if mesh.OwnerID != user.ID {
		writeError(w, http.StatusForbidden, "You don't have permission to update this mesh")
		return
	}

	params := meshdb.UpdateMeshParams{
		ID:          meshID,
		Name:        req.Name,
		Description: req.Description,
	}

	updatedMesh, err := s.DB().UpdateMesh(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update mesh")
		return
	}

	writeJSON(w, http.StatusOK, updatedMesh)
}

// handleDeleteMesh handles deleting a mesh
func (s *Server) handleDeleteMesh(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := chi.URLParam(r, "meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	// Get existing mesh to check ownership
	mesh, err := s.DB().GetMeshByID(r.Context(), meshID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get mesh")
		return
	}

	// Check if user owns the mesh (basic check - will need mesh_access later)
	if mesh.OwnerID != user.ID {
		writeError(w, http.StatusForbidden, "You don't have permission to delete this mesh")
		return
	}

	if err := s.DB().DeleteMesh(r.Context(), meshID);	 err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete mesh")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Mesh deleted successfully",
	})
}
