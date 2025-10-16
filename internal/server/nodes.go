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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/skandragon/meshmgr/meshdb"
)

// CreateNodeRequest represents a request to create a node
type CreateNodeRequest struct {
	HardwareID    string  `json:"hardware_id"`
	Name          string  `json:"name"`
	LongName      string  `json:"long_name"`
	Role          *string `json:"role,omitempty"`
	PublicKey     *string `json:"public_key,omitempty"`
	PrivateKey    *string `json:"private_key,omitempty"`
	Status        *string `json:"status,omitempty"`
	Unmessageable *bool   `json:"unmessageable,omitempty"`
}

// UpdateNodeRequest represents a request to update a node
type UpdateNodeRequest struct {
	Name           *string `json:"name,omitempty"`
	LongName       *string `json:"long_name,omitempty"`
	Role           *string `json:"role,omitempty"`
	PublicKey      *string `json:"public_key,omitempty"`
	PrivateKey     *string `json:"private_key,omitempty"`
	Status         *string `json:"status,omitempty"`
	Unmessageable  *bool   `json:"unmessageable,omitempty"`
	PendingChanges *bool   `json:"pending_changes,omitempty"`
}

// UpdateNodeStatusRequest represents a request to update node status
type UpdateNodeStatusRequest struct {
	Status string `json:"status"`
}

// handleListNodes handles listing nodes for a mesh
func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
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

	nodes, err := s.DB().ListNodesByMesh(r.Context(), meshID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list nodes")
		return
	}

	writeJSON(w, http.StatusOK, nodes)
}

// handleGetNode handles getting a single node
func (s *Server) handleGetNode(w http.ResponseWriter, r *http.Request) {
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

	nodeIDStr := r.PathValue("nodeID")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid node ID")
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

	node, err := s.DB().GetNode(r.Context(), nodeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Node not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get node")
		return
	}

	// Verify the node belongs to the specified mesh
	if node.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Node not found")
		return
	}

	writeJSON(w, http.StatusOK, node)
}

// handleCreateNode handles creating a new node
func (s *Server) handleCreateNode(w http.ResponseWriter, r *http.Request) {
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

	var req CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.HardwareID == "" {
		writeError(w, http.StatusBadRequest, "Hardware ID is required")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if req.LongName == "" {
		writeError(w, http.StatusBadRequest, "Long name is required")
		return
	}

	// Create the node
	unmessageable := false
	if req.Unmessageable != nil {
		unmessageable = *req.Unmessageable
	}

	node, err := s.DB().CreateNode(r.Context(), meshdb.CreateNodeParams{
		MeshID:        meshID,
		HardwareID:    req.HardwareID,
		Name:          req.Name,
		LongName:      req.LongName,
		Role:          req.Role,
		PublicKey:     req.PublicKey,
		PrivateKey:    req.PrivateKey,
		Status:        req.Status,
		Unmessageable: unmessageable,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create node")
		return
	}

	writeJSON(w, http.StatusCreated, node)
}

// handleUpdateNode handles updating a node
func (s *Server) handleUpdateNode(w http.ResponseWriter, r *http.Request) {
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

	nodeIDStr := r.PathValue("nodeID")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid node ID")
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

	// Verify the node exists and belongs to this mesh
	node, err := s.DB().GetNode(r.Context(), nodeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Node not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get node")
		return
	}

	if node.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Node not found")
		return
	}

	var req UpdateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update the node
	var unmessageable pgtype.Bool
	if req.Unmessageable != nil {
		unmessageable = pgtype.Bool{Bool: *req.Unmessageable, Valid: true}
	}

	var pendingChanges pgtype.Bool
	if req.PendingChanges != nil {
		pendingChanges = pgtype.Bool{Bool: *req.PendingChanges, Valid: true}
	}

	updatedNode, err := s.DB().UpdateNode(r.Context(), meshdb.UpdateNodeParams{
		ID:             nodeID,
		Name:           req.Name,
		LongName:       req.LongName,
		Role:           req.Role,
		PublicKey:      req.PublicKey,
		PrivateKey:     req.PrivateKey,
		Status:         req.Status,
		Unmessageable:  unmessageable,
		PendingChanges: pendingChanges,
		LastSeen:       nil, // Don't update last_seen on manual updates
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update node")
		return
	}

	writeJSON(w, http.StatusOK, updatedNode)
}

// handleUpdateNodeStatus handles updating a node's status
func (s *Server) handleUpdateNodeStatus(w http.ResponseWriter, r *http.Request) {
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

	nodeIDStr := r.PathValue("nodeID")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid node ID")
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

	// Verify the node exists and belongs to this mesh
	node, err := s.DB().GetNode(r.Context(), nodeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Node not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get node")
		return
	}

	if node.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Node not found")
		return
	}

	var req UpdateNodeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status == "" {
		writeError(w, http.StatusBadRequest, "Status is required")
		return
	}

	// Update the node status (this also updates last_seen)
	updatedNode, err := s.DB().UpdateNodeStatus(r.Context(), meshdb.UpdateNodeStatusParams{
		ID:     nodeID,
		Status: &req.Status,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update node status")
		return
	}

	writeJSON(w, http.StatusOK, updatedNode)
}

// handleDeleteNode handles deleting a node
func (s *Server) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
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

	nodeIDStr := r.PathValue("nodeID")
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid node ID")
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

	// Verify the node exists and belongs to this mesh
	node, err := s.DB().GetNode(r.Context(), nodeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			writeError(w, http.StatusNotFound, "Node not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get node")
		return
	}

	if node.MeshID != meshID {
		writeError(w, http.StatusNotFound, "Node not found")
		return
	}

	// Delete the node
	if err := s.DB().DeleteNode(r.Context(), nodeID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete node")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Node deleted successfully",
	})
}
