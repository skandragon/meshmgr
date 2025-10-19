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

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		email      string
		secret     string
		expiration time.Duration
		wantErr    bool
	}{
		{
			name:       "valid token generation",
			userID:     1,
			email:      "test@example.com",
			secret:     "test-secret-key",
			expiration: 24 * time.Hour,
			wantErr:    false,
		},
		{
			name:       "valid token with short expiration",
			userID:     999,
			email:      "user@test.com",
			secret:     "another-secret",
			expiration: 1 * time.Minute,
			wantErr:    false,
		},
		{
			name:       "empty secret",
			userID:     1,
			email:      "test@example.com",
			secret:     "",
			expiration: 24 * time.Hour,
			wantErr:    false, // JWT allows empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.secret, tt.expiration)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify the token can be validated
			claims, err := ValidateToken(token, tt.secret)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, tt.email, claims.Email)
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := int64(123)
	email := "test@example.com"

	validToken, err := GenerateToken(userID, email, secret, 24*time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		secret  string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			secret:  secret,
			wantErr: false,
		},
		{
			name:    "wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.valid.jwt",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "completely invalid token",
			token:   "invalid",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secret)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, claims)
			assert.Equal(t, userID, claims.UserID)
			assert.Equal(t, email, claims.Email)
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test-secret"
	userID := int64(1)
	email := "test@example.com"

	// Create an expired token (already expired)
	expiredToken, err := GenerateToken(userID, email, secret, -1*time.Hour)
	require.NoError(t, err)

	// Try to validate expired token
	claims, err := ValidateToken(expiredToken, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "expired")
}

func TestTokenClaims(t *testing.T) {
	secret := "test-secret"
	userID := int64(42)
	email := "user@example.com"
	expiration := 2 * time.Hour

	token, err := GenerateToken(userID, email, secret, expiration)
	require.NoError(t, err)

	claims, err := ValidateToken(token, secret)
	require.NoError(t, err)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)

	// Check that timestamps are reasonable
	now := time.Now()
	assert.True(t, claims.IssuedAt.Before(now.Add(time.Second)))
	assert.True(t, claims.NotBefore.Before(now.Add(time.Second)))
	assert.True(t, claims.ExpiresAt.After(now))
}

func TestGenerateRandomToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "generate random token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token1, err1 := GenerateRandomToken()
			require.NoError(t, err1)
			assert.NotEmpty(t, token1)

			token2, err2 := GenerateRandomToken()
			require.NoError(t, err2)
			assert.NotEmpty(t, token2)

			// Tokens should be unique
			assert.NotEqual(t, token1, token2)

			// Tokens should be base64 encoded
			assert.Greater(t, len(token1), 32)
		})
	}
}

func TestTokenUniqueness(t *testing.T) {
	secret := "test-secret"
	userID := int64(1)
	email := "test@example.com"
	expiration := 24 * time.Hour

	// Generate a token and verify it can be validated
	token1, err := GenerateToken(userID, email, secret, expiration)
	require.NoError(t, err)

	// Wait to ensure different timestamp (JWT uses seconds)
	time.Sleep(1100 * time.Millisecond)

	token2, err := GenerateToken(userID, email, secret, expiration)
	require.NoError(t, err)

	// Tokens generated at different times should be different
	assert.NotEqual(t, token1, token2, "tokens generated at different times should be different")

	// But both should validate correctly
	claims1, err1 := ValidateToken(token1, secret)
	require.NoError(t, err1)
	assert.Equal(t, userID, claims1.UserID)

	claims2, err2 := ValidateToken(token2, secret)
	require.NoError(t, err2)
	assert.Equal(t, userID, claims2.UserID)
}
