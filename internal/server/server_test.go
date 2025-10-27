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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/skandragon/meshmgr/internal/config"
	"github.com/skandragon/meshmgr/meshdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer wraps Server and the gnomock container for testing
type testServer struct {
	server    *Server
	container *gnomock.Container
	db        *pgxpool.Pool
}

// setupTestServer creates a test server with a fresh PostgreSQL instance
func setupTestServer(t *testing.T) *testServer {
	t.Helper()

	// Start PostgreSQL container
	p := postgres.Preset(
		postgres.WithUser("testuser", "testpass"),
		postgres.WithDatabase("testdb"),
	)

	container, err := gnomock.Start(p)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, gnomock.Stop(container))
	})

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=testuser password=testpass dbname=testdb sslmode=disable",
		container.Host, container.DefaultPort(),
	)

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	// Run migrations
	runMigrations(t, pool)

	// Create config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0, // Let the test framework choose a port
		},
		Database: config.DatabaseConfig{
			Host:     container.Host,
			Port:     container.DefaultPort(),
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret-key",
			JWTExpiration: 24 * time.Hour,
			BCryptCost:    4, // Use low cost for faster tests
		},
	}

	// Create server
	srv := &Server{
		config: cfg,
		db:     pool,
		mux:    http.NewServeMux(),
	}
	srv.setupRoutes()

	return &testServer{
		server:    srv,
		container: container,
		db:        pool,
	}
}

// runMigrations applies database migrations
func runMigrations(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	migrations := []string{
		"1760580803_initial_schema.up.sql",
		"1760587537_add_nodes.up.sql",
		"1760605927_add_mesh_lora_config.up.sql",
		"1760605928_add_node_state_tracking.up.sql",
		"1760605929_add_node_admin_keys_mapping.up.sql",
		"1760631305_update_preset_names.up.sql",
		"1760632000_update_frequency_slot_range.up.sql",
		"1760640000_add_device_config_storage.up.sql",
		"1760650000_add_api_keys.up.sql",
	}

	for _, migration := range migrations {
		migrationPath := filepath.Join("..", "..", "meshdb", "migrations", migration)
		migrationSQL, err := os.ReadFile(migrationPath)
		require.NoError(t, err)

		_, err = pool.Exec(context.Background(), string(migrationSQL))
		require.NoError(t, err)
	}
}

// makeRequest is a helper to make HTTP requests to the test server
func (ts *testServer) makeRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	ts.server.mux.ServeHTTP(rr, req)

	return rr
}

func TestHealthEndpoint(t *testing.T) {
	ts := setupTestServer(t)

	rr := ts.makeRequest(t, "GET", "/health", nil, "")

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestRegisterEndpoint(t *testing.T) {
	ts := setupTestServer(t)

	tests := []struct {
		name       string
		request    RegisterRequest
		wantStatus int
		checkResp  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			request: RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123",
				DisplayName: "Test User",
			},
			wantStatus: http.StatusCreated,
			checkResp: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var resp AuthResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, "test@example.com", resp.User.Email)
				assert.Equal(t, "Test User", resp.User.DisplayName)
			},
		},
		{
			name: "missing email",
			request: RegisterRequest{
				Password:    "password123",
				DisplayName: "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			request: RegisterRequest{
				Email:       "test2@example.com",
				DisplayName: "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing display name",
			request: RegisterRequest{
				Email:    "test3@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := ts.makeRequest(t, "POST", "/api/auth/register", tt.request, "")

			assert.Equal(t, tt.wantStatus, rr.Code)

			if tt.checkResp != nil {
				tt.checkResp(t, rr)
			}
		})
	}
}

func TestLoginEndpoint(t *testing.T) {
	ts := setupTestServer(t)

	// First, register a user
	registerReq := RegisterRequest{
		Email:       "login@example.com",
		Password:    "mypassword",
		DisplayName: "Login User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)

	tests := []struct {
		name       string
		request    LoginRequest
		wantStatus int
		checkResp  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			request: LoginRequest{
				Email:    "login@example.com",
				Password: "mypassword",
			},
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var resp AuthResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, "login@example.com", resp.User.Email)
			},
		},
		{
			name: "wrong password",
			request: LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "nonexistent user",
			request: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			request: LoginRequest{
				Password: "password",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := ts.makeRequest(t, "POST", "/api/auth/login", tt.request, "")

			assert.Equal(t, tt.wantStatus, rr.Code)

			if tt.checkResp != nil {
				tt.checkResp(t, rr)
			}
		})
	}
}

func TestMeEndpoint(t *testing.T) {
	ts := setupTestServer(t)

	// Register and get token
	registerReq := RegisterRequest{
		Email:       "me@example.com",
		Password:    "password",
		DisplayName: "Me User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)

	var authResp AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &authResp)
	require.NoError(t, err)
	token := authResp.Token

	tests := []struct {
		name       string
		token      string
		wantStatus int
		checkResp  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "valid token",
			token:      token,
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var user meshdb.User
				err := json.Unmarshal(rr.Body.Bytes(), &user)
				require.NoError(t, err)
				assert.Equal(t, "me@example.com", user.Email)
				assert.Equal(t, "Me User", user.DisplayName)
			},
		},
		{
			name:       "no token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := ts.makeRequest(t, "GET", "/api/auth/me", nil, tt.token)

			assert.Equal(t, tt.wantStatus, rr.Code)

			if tt.checkResp != nil {
				tt.checkResp(t, rr)
			}
		})
	}
}

func TestMeshCRUD(t *testing.T) {
	ts := setupTestServer(t)

	// Register and get token
	registerReq := RegisterRequest{
		Email:       "mesh@example.com",
		Password:    "password",
		DisplayName: "Mesh User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)

	var authResp AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &authResp)
	require.NoError(t, err)
	token := authResp.Token

	// Test list meshes (should be empty)
	rr = ts.makeRequest(t, "GET", "/api/meshes", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test create mesh
	desc := "Test mesh description"
	createReq := CreateMeshRequest{
		Name:        "My Test Mesh",
		Description: &desc,
	}
	rr = ts.makeRequest(t, "POST", "/api/meshes", createReq, token)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var mesh meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	assert.Equal(t, "My Test Mesh", mesh.Name)
	assert.Equal(t, "Test mesh description", *mesh.Description)
	meshID := mesh.ID

	// Test get mesh
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test update mesh
	newDesc := "Updated description"
	updateReq := UpdateMeshRequest{
		Description: &newDesc,
	}
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d", meshID), updateReq, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", *mesh.Description)

	// Test list meshes (should have one)
	rr = ts.makeRequest(t, "GET", "/api/meshes", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	var meshes []meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &meshes)
	require.NoError(t, err)
	assert.Len(t, meshes, 1)

	// Test delete mesh
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test list meshes (should be empty again)
	rr = ts.makeRequest(t, "GET", "/api/meshes", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestMeshOwnership(t *testing.T) {
	ts := setupTestServer(t)
	var err error

	// Register two users
	user1Req := RegisterRequest{
		Email:       "user1@example.com",
		Password:    "password1",
		DisplayName: "User One",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", user1Req, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var auth1 AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &auth1)
	require.NoError(t, err)
	token1 := auth1.Token

	user2Req := RegisterRequest{
		Email:       "user2@example.com",
		Password:    "password2",
		DisplayName: "User Two",
	}
	rr = ts.makeRequest(t, "POST", "/api/auth/register", user2Req, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var auth2 AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &auth2)
	require.NoError(t, err)
	token2 := auth2.Token

	// User 1 creates a mesh
	desc := "User 1's mesh"
	createReq := CreateMeshRequest{
		Name:        "User 1 Mesh",
		Description: &desc,
	}
	rr = ts.makeRequest(t, "POST", "/api/meshes", createReq, token1)
	require.Equal(t, http.StatusCreated, rr.Code)
	var mesh meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	meshID := mesh.ID

	// User 2 should not be able to see User 1's mesh (no access)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// User 2 should not be able to update User 1's mesh (no access)
	updateDesc := "Hacked description"
	updateReq := UpdateMeshRequest{
		Description: &updateDesc,
	}
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d", meshID), updateReq, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// User 2 should not be able to delete User 1's mesh (no access)
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// User 1 should still be able to update
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d", meshID), updateReq, token1)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLogoutEndpoint(t *testing.T) {
	ts := setupTestServer(t)

	// Register user
	registerReq := RegisterRequest{
		Email:       "logout@example.com",
		Password:    "password",
		DisplayName: "Logout User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)

	var authResp AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &authResp)
	require.NoError(t, err)
	token := authResp.Token

	// Verify token works
	rr = ts.makeRequest(t, "GET", "/api/auth/me", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Logout
	rr = ts.makeRequest(t, "POST", "/api/auth/logout", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Token should still work (JWT tokens can't be invalidated server-side in this implementation)
	// This test documents current behavior - sessions are deleted but JWT still validates
	rr = ts.makeRequest(t, "GET", "/api/auth/me", nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestMeshAccessManagement is disabled because mesh access sharing feature was removed
func TestMeshAccessManagement_DISABLED(t *testing.T) {
	t.Skip("Mesh access feature removed")
	ts := setupTestServer(t)
	var err error

	// Register two users
	user1Req := RegisterRequest{
		Email:       "owner@example.com",
		Password:    "password1",
		DisplayName: "Owner User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", user1Req, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var auth1 AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &auth1)
	require.NoError(t, err)
	token1 := auth1.Token

	user2Req := RegisterRequest{
		Email:       "viewer@example.com",
		Password:    "password2",
		DisplayName: "Viewer User",
	}
	rr = ts.makeRequest(t, "POST", "/api/auth/register", user2Req, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var auth2 AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &auth2)
	require.NoError(t, err)
	token2 := auth2.Token

	// User 1 creates a mesh
	desc := "Test mesh"
	createReq := CreateMeshRequest{
		Name:        "Test Mesh",
		Description: &desc,
	}
	rr = ts.makeRequest(t, "POST", "/api/meshes", createReq, token1)
	require.Equal(t, http.StatusCreated, rr.Code)
	var mesh meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	meshID := mesh.ID

	// User 2 cannot access the mesh initially
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Owner grants viewer access to User 2
	grantReq := GrantAccessRequest{
		UserEmail:   "viewer@example.com",
		AccessLevel: "viewer",
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/access", meshID), grantReq, token1)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// User 2 can now view the mesh
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusOK, rr.Code)

	// User 2 cannot update the mesh (viewer level)
	updateReq := UpdateMeshRequest{Description: &desc}
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d", meshID), updateReq, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// List mesh access
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/access", meshID), nil, token1)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Update User 2's access to admin
	updateAccessReq := UpdateAccessRequest{AccessLevel: "admin"}
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d/access/%d", meshID, auth2.User.ID), updateAccessReq, token1)
	assert.Equal(t, http.StatusOK, rr.Code)

	// User 2 can now update the mesh (admin level)
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d", meshID), updateReq, token2)
	assert.Equal(t, http.StatusOK, rr.Code)

	// User 2 still cannot delete the mesh (requires owner)
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Revoke User 2's access
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d/access/%d", meshID, auth2.User.ID), nil, token1)
	assert.Equal(t, http.StatusOK, rr.Code)

	// User 2 cannot access the mesh anymore
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d", meshID), nil, token2)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAdminKeys(t *testing.T) {
	ts := setupTestServer(t)
	var err error

	// Register user and create mesh
	registerReq := RegisterRequest{
		Email:       "admin@example.com",
		Password:    "password",
		DisplayName: "Admin User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var authResp AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &authResp)
	require.NoError(t, err)
	token := authResp.Token

	desc := "Test mesh"
	createReq := CreateMeshRequest{
		Name:        "Test Mesh",
		Description: &desc,
	}
	rr = ts.makeRequest(t, "POST", "/api/meshes", createReq, token)
	require.Equal(t, http.StatusCreated, rr.Code)
	var mesh meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	meshID := mesh.ID

	// List admin keys (should be empty)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Add first admin key
	keyName1 := "Key 1"
	addKeyReq1 := CreateAdminKeyRequest{
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		KeyName:   &keyName1,
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), addKeyReq1, token)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var key1 meshdb.AdminKey
	err = json.Unmarshal(rr.Body.Bytes(), &key1)
	require.NoError(t, err)

	// Add second admin key
	keyName2 := "Key 2"
	addKeyReq2 := CreateAdminKeyRequest{
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD...",
		KeyName:   &keyName2,
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), addKeyReq2, token)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Add third admin key
	keyName3 := "Key 3"
	addKeyReq3 := CreateAdminKeyRequest{
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQE...",
		KeyName:   &keyName3,
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), addKeyReq3, token)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Try to add fourth key (should fail - max 3)
	keyName4 := "Key 4"
	addKeyReq4 := CreateAdminKeyRequest{
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQF...",
		KeyName:   &keyName4,
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), addKeyReq4, token)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// List admin keys (should have 3)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var keys []meshdb.AdminKey
	err = json.Unmarshal(rr.Body.Bytes(), &keys)
	require.NoError(t, err)
	assert.Len(t, keys, 3)

	// Get single admin key
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/admin-keys/%d", meshID, key1.ID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Delete an admin key
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d/admin-keys/%d", meshID, key1.ID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// List admin keys (should have 2)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &keys)
	require.NoError(t, err)
	assert.Len(t, keys, 2)

	// Can now add another key
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/admin-keys", meshID), addKeyReq4, token)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestNodeManagement(t *testing.T) {
	ts := setupTestServer(t)
	var err error

	// Register user and create mesh
	registerReq := RegisterRequest{
		Email:       "node@example.com",
		Password:    "password",
		DisplayName: "Node User",
	}
	rr := ts.makeRequest(t, "POST", "/api/auth/register", registerReq, "")
	require.Equal(t, http.StatusCreated, rr.Code)
	var authResp AuthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &authResp)
	require.NoError(t, err)
	token := authResp.Token

	desc := "Test mesh"
	createReq := CreateMeshRequest{
		Name:        "Test Mesh",
		Description: &desc,
	}
	rr = ts.makeRequest(t, "POST", "/api/meshes", createReq, token)
	require.Equal(t, http.StatusCreated, rr.Code)
	var mesh meshdb.Mesh
	err = json.Unmarshal(rr.Body.Bytes(), &mesh)
	require.NoError(t, err)
	meshID := mesh.ID

	// List nodes (should be empty)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/nodes", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Create a node
	role := "CLIENT"
	status := "online"
	createNodeReq := CreateNodeRequest{
		HardwareID: "!abc123",
		Name:       "node1",
		LongName:   "Node One",
		Role:       &role,
		Status:     &status,
	}
	rr = ts.makeRequest(t, "POST", fmt.Sprintf("/api/meshes/%d/nodes", meshID), createNodeReq, token)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var node meshdb.Node
	err = json.Unmarshal(rr.Body.Bytes(), &node)
	require.NoError(t, err)
	nodeID := node.ID
	assert.Equal(t, "!abc123", node.HardwareID)
	assert.Equal(t, "node1", node.Name)

	// Get the node
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/nodes/%d", meshID, nodeID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Update the node
	newName := "updated-node1"
	updateNodeReq := UpdateNodeRequest{
		Name: &newName,
	}
	rr = ts.makeRequest(t, "PUT", fmt.Sprintf("/api/meshes/%d/nodes/%d", meshID, nodeID), updateNodeReq, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &node)
	require.NoError(t, err)
	assert.Equal(t, "updated-node1", node.Name)
	assert.Equal(t, "Node One", node.LongName) // Should remain unchanged

	// Update node status
	updateStatusReq := UpdateNodeStatusRequest{
		Status: "offline",
	}
	rr = ts.makeRequest(t, "PATCH", fmt.Sprintf("/api/meshes/%d/nodes/%d/status", meshID, nodeID), updateStatusReq, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &node)
	require.NoError(t, err)
	assert.Equal(t, "offline", *node.Status)

	// List nodes (should have 1)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/nodes", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var nodes []meshdb.Node
	err = json.Unmarshal(rr.Body.Bytes(), &nodes)
	require.NoError(t, err)
	assert.Len(t, nodes, 1)

	// Delete the node
	rr = ts.makeRequest(t, "DELETE", fmt.Sprintf("/api/meshes/%d/nodes/%d", meshID, nodeID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)

	// List nodes (should be empty)
	rr = ts.makeRequest(t, "GET", fmt.Sprintf("/api/meshes/%d/nodes", meshID), nil, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &nodes)
	require.NoError(t, err)
	assert.Len(t, nodes, 0)
}
