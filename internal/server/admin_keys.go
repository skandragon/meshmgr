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

	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/meshdb"
)

// CreateAdminKeyRequest represents a request to add an admin key
type CreateAdminKeyRequest struct {
	PublicKey string  `json:"public_key"`
	KeyName   *string `json:"key_name,omitempty"`
}

// handleListAdminKeys handles listing admin keys for a mesh
func (s *Server) handleListAdminKeys(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := r.PathValue("meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	// Check if user has at least viewer access
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelViewer); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	keys, err := s.DB().ListAdminKeysByMesh(r.Context(), meshID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list admin keys")
		return
	}

	writeJSON(w, http.StatusOK, keys)
}

// handleGetAdminKey handles getting a single admin key
func (s *Server) handleGetAdminKey(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := r.PathValue("meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	keyIDStr := r.PathValue("keyID")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	// Check if user has at least viewer access
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelViewer); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	key, err := s.DB().GetAdminKey(r.Context(), keyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Admin key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get admin key")
		return
	}

	// Verify the key belongs to the specified mesh
	if key.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Admin key not found")
		return
	}

	writeJSON(w, http.StatusOK, key)
}

// handleCreateAdminKey handles adding a new admin key
func (s *Server) handleCreateAdminKey(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := r.PathValue("meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	// Check if user has at least admin access
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelAdmin); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	var req CreateAdminKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PublicKey == "" {
		writeError(w, http.StatusBadRequest, "Public key is required")
		return
	}

	// Check if mesh already has 3 admin keys
	count, err := s.DB().CountAdminKeysByMesh(r.Context(), meshID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check admin key count")
		return
	}

	if count >= 3 {
		writeError(w, http.StatusBadRequest, "Mesh already has maximum of 3 admin keys")
		return
	}

	// Decode base64 public key to bytes
	publicKeyBytes := []byte(req.PublicKey)

	// Create the admin key
	key, err := s.DB().CreateAdminKey(r.Context(), meshdb.CreateAdminKeyParams{
		MeshID:    meshID,
		PublicKey: publicKeyBytes,
		KeyName:   req.KeyName,
		AddedBy:   user.ID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create admin key")
		return
	}

	writeJSON(w, http.StatusCreated, key)
}

// handleDeleteAdminKey handles deleting an admin key
func (s *Server) handleDeleteAdminKey(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	meshIDStr := r.PathValue("meshID")
	meshID, err := strconv.ParseInt(meshIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid mesh ID")
		return
	}

	keyIDStr := r.PathValue("keyID")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid key ID")
		return
	}

	// Check if user has at least admin access
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelAdmin); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	// Verify the key exists and belongs to this mesh
	key, err := s.DB().GetAdminKey(r.Context(), keyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Admin key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get admin key")
		return
	}

	if key.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Admin key not found")
		return
	}

	// Delete the key
	if err := s.DB().DeleteAdminKey(r.Context(), keyID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete admin key")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Admin key deleted successfully",
	})
}
