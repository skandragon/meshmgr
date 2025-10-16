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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		cost     int
		wantErr  bool
	}{
		{
			name:     "valid password with default cost",
			password: "mySecurePassword123!",
			cost:     12,
			wantErr:  false,
		},
		{
			name:     "valid password with minimum cost",
			password: "test123",
			cost:     4,
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			cost:     12,
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "long password",
			password: string(make([]byte, 72)), // bcrypt max is 72 bytes
			cost:     12,
			wantErr:  false,
		},
		{
			name:     "very long password",
			password: string(make([]byte, 100)), // longer than bcrypt max
			cost:     12,
			wantErr:  true, // bcrypt errors on passwords > 72 bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password, tt.cost)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash, "hash should not equal plaintext password")

			// Verify the hash can be used to check the password
			assert.True(t, CheckPassword(tt.password, hash))
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "correctPassword123"
	hash, err := HashPassword(password, 12)
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "not-a-valid-hash",
			want:     false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashPasswordDeterminism(t *testing.T) {
	password := "testPassword123"

	hash1, err1 := HashPassword(password, 12)
	require.NoError(t, err1)

	hash2, err2 := HashPassword(password, 12)
	require.NoError(t, err2)

	// Hashes should be different due to random salt
	assert.NotEqual(t, hash1, hash2, "hashes should differ due to random salt")

	// But both should validate the same password
	assert.True(t, CheckPassword(password, hash1))
	assert.True(t, CheckPassword(password, hash2))
}

func TestPasswordCaseSensitivity(t *testing.T) {
	password := "TestPassword123"
	hash, err := HashPassword(password, 12)
	require.NoError(t, err)

	// Case should matter
	assert.True(t, CheckPassword("TestPassword123", hash))
	assert.False(t, CheckPassword("testpassword123", hash))
	assert.False(t, CheckPassword("TESTPASSWORD123", hash))
}
