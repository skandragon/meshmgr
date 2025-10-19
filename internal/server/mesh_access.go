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

// TODO: Uncomment these imports when implementing mesh access control endpoints
/*
import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/meshdb"
)
*/

// GrantAccessRequest represents a request to grant mesh access
type GrantAccessRequest struct {
	UserEmail   string `json:"user_email"`
	AccessLevel string `json:"access_level"`
}

// UpdateAccessRequest represents a request to update access level
type UpdateAccessRequest struct {
	AccessLevel string `json:"access_level"`
}

// TODO: Implement mesh access control endpoints when needed
// The following functions are commented out to avoid unused function warnings
// They provide a complete implementation for mesh access control features

/*
// handleListMeshAccess handles listing all users with access to a mesh
func (s *Server) handleListMeshAccess(w http.ResponseWriter, r *http.Request) {
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

	access, err := s.DB().ListMeshAccessByMesh(r.Context(), meshID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list mesh access")
		return
	}

	writeJSON(w, http.StatusOK, access)
}

// handleGrantMeshAccess handles granting access to a user
func (s *Server) handleGrantMeshAccess(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is the mesh owner
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelOwner); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	var req GrantAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserEmail == "" {
		writeError(w, http.StatusBadRequest, "User email is required")
		return
	}

	if req.AccessLevel == "" {
		writeError(w, http.StatusBadRequest, "Access level is required")
		return
	}

	// Validate access level
	if req.AccessLevel != "admin" && req.AccessLevel != "viewer" {
		writeError(w, http.StatusBadRequest, "Access level must be 'admin' or 'viewer'")
		return
	}

	// Find user by email
	targetUser, err := s.DB().GetUserByEmail(r.Context(), req.UserEmail)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to find user")
		return
	}

	// Grant access
	grantedBy := user.ID
	access, err := s.DB().GrantMeshAccess(r.Context(), meshdb.GrantMeshAccessParams{
		MeshID:      meshID,
		UserID:      targetUser.ID,
		AccessLevel: req.AccessLevel,
		GrantedBy:   &grantedBy,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to grant access")
		return
	}

	writeJSON(w, http.StatusCreated, access)
}

// handleUpdateMeshAccess handles updating a user's access level
func (s *Server) handleUpdateMeshAccess(w http.ResponseWriter, r *http.Request) {
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

	targetUserIDStr := r.PathValue("userID")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Check if user is the mesh owner
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelOwner); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	var req UpdateAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AccessLevel == "" {
		writeError(w, http.StatusBadRequest, "Access level is required")
		return
	}

	// Validate access level
	if req.AccessLevel != "admin" && req.AccessLevel != "viewer" {
		writeError(w, http.StatusBadRequest, "Access level must be 'admin' or 'viewer'")
		return
	}

	// Update access
	access, err := s.DB().UpdateMeshAccess(r.Context(), meshdb.UpdateMeshAccessParams{
		MeshID:      meshID,
		UserID:      targetUserID,
		AccessLevel: req.AccessLevel,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Access record not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update access")
		return
	}

	writeJSON(w, http.StatusOK, access)
}

// handleRevokeMeshAccess handles revoking a user's access
func (s *Server) handleRevokeMeshAccess(w http.ResponseWriter, r *http.Request) {
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

	targetUserIDStr := r.PathValue("userID")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Check if user is the mesh owner
	if _, err := s.requireMeshAccess(r.Context(), user.ID, meshID, AccessLevelOwner); err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Mesh not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to check permissions")
		return
	}

	// Revoke access
	if err := s.DB().RevokeMeshAccess(r.Context(), meshdb.RevokeMeshAccessParams{
		MeshID: meshID,
		UserID: targetUserID,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to revoke access")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Access revoked successfully",
	})
}
*/
