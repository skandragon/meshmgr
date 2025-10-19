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

	"github.com/jackc/pgx/v5"
	"github.com/skandragon/meshmgr/meshdb"
)

// AccessLevel represents the permission level
type AccessLevel string

const (
	AccessLevelOwner  AccessLevel = "owner"
	AccessLevelAdmin  AccessLevel = "admin"
	AccessLevelViewer AccessLevel = "viewer"
)

// checkMeshAccess checks if a user has at least the specified access level for a mesh
// Returns the actual access level if user has access, empty string if not
func (s *Server) checkMeshAccess(ctx context.Context, userID, meshID int64, minLevel AccessLevel) (AccessLevel, error) {
	// First check if user is the mesh owner
	mesh, err := s.DB().GetMeshByID(ctx, meshID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil // Mesh doesn't exist
		}
		return "", err
	}

	if mesh.OwnerID == userID {
		return AccessLevelOwner, nil
	}

	// Check mesh_access table
	accessLevel, err := s.DB().CheckUserMeshAccess(ctx, meshdb.CheckUserMeshAccessParams{
		MeshID: meshID,
		UserID: userID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil // No access
		}
		return "", err
	}

	level := AccessLevel(accessLevel)

	// Check if user has sufficient access
	if hasAccess(level, minLevel) {
		return level, nil
	}

	return "", nil
}

// hasAccess checks if userLevel meets the minimum required level
func hasAccess(userLevel, minLevel AccessLevel) bool {
	levels := map[AccessLevel]int{
		AccessLevelViewer: 1,
		AccessLevelAdmin:  2,
		AccessLevelOwner:  3,
	}

	return levels[userLevel] >= levels[minLevel]
}

// requireMeshAccess is a helper that returns an error if user doesn't have access
func (s *Server) requireMeshAccess(ctx context.Context, userID, meshID int64, minLevel AccessLevel) (AccessLevel, error) {
	level, err := s.checkMeshAccess(ctx, userID, meshID, minLevel)
	if err != nil {
		return "", err
	}
	if level == "" {
		return "", pgx.ErrNoRows // Use this to signal "not found/no access"
	}
	return level, nil
}
