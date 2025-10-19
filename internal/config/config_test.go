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

package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Save original env vars and restore after test
	originalEnv := map[string]string{
		"JWT_SECRET":    os.Getenv("JWT_SECRET"),
		"SERVER_HOST":   os.Getenv("SERVER_HOST"),
		"SERVER_PORT":   os.Getenv("SERVER_PORT"),
		"DB_HOST":       os.Getenv("DB_HOST"),
		"DB_PORT":       os.Getenv("DB_PORT"),
		"DB_USER":       os.Getenv("DB_USER"),
		"DB_PASSWORD":   os.Getenv("DB_PASSWORD"),
		"DB_NAME":       os.Getenv("DB_NAME"),
		"DB_SSLMODE":    os.Getenv("DB_SSLMODE"),
		"JWT_EXPIRATION": os.Getenv("JWT_EXPIRATION"),
		"BCRYPT_COST":   os.Getenv("BCRYPT_COST"),
	}
	defer func() {
		for k, v := range originalEnv {
			if v == "" {
				require.NoError(t, os.Unsetenv(k))
			} else {
				require.NoError(t, os.Setenv(k, v))
			}
		}
	}()

	tests := []struct {
		name          string
		setupEnv      func()
		wantErr       bool
		checkConfig   func(*testing.T, *Config)
	}{
		{
			name: "default values with JWT secret set",
			setupEnv: func() {
				os.Clearenv()
				require.NoError(t, os.Setenv("JWT_SECRET", "test-secret"))
			},
			wantErr: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "0.0.0.0", cfg.Server.Host)
				assert.Equal(t, 8080, cfg.Server.Port)
				assert.Equal(t, "localhost", cfg.Database.Host)
				assert.Equal(t, 5432, cfg.Database.Port)
				assert.Equal(t, "meshmgr", cfg.Database.User)
				assert.Equal(t, "", cfg.Database.Password)
				assert.Equal(t, "meshmgr", cfg.Database.DBName)
				assert.Equal(t, "disable", cfg.Database.SSLMode)
				assert.Equal(t, "test-secret", cfg.Auth.JWTSecret)
				assert.Equal(t, 7*24*time.Hour, cfg.Auth.JWTExpiration)
				assert.Equal(t, 12, cfg.Auth.BCryptCost)
			},
		},
		{
			name: "custom values",
			setupEnv: func() {
				os.Clearenv()
				require.NoError(t, os.Setenv("JWT_SECRET", "custom-secret"))
				require.NoError(t, os.Setenv("SERVER_HOST", "127.0.0.1"))
				require.NoError(t, os.Setenv("SERVER_PORT", "9000"))
				require.NoError(t, os.Setenv("DB_HOST", "dbhost"))
				require.NoError(t, os.Setenv("DB_PORT", "3306"))
				require.NoError(t, os.Setenv("DB_USER", "customuser"))
				require.NoError(t, os.Setenv("DB_PASSWORD", "custompass"))
				require.NoError(t, os.Setenv("DB_NAME", "customdb"))
				require.NoError(t, os.Setenv("DB_SSLMODE", "require"))
				require.NoError(t, os.Setenv("JWT_EXPIRATION", "24h"))
				require.NoError(t, os.Setenv("BCRYPT_COST", "10"))
			},
			wantErr: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "127.0.0.1", cfg.Server.Host)
				assert.Equal(t, 9000, cfg.Server.Port)
				assert.Equal(t, "dbhost", cfg.Database.Host)
				assert.Equal(t, 3306, cfg.Database.Port)
				assert.Equal(t, "customuser", cfg.Database.User)
				assert.Equal(t, "custompass", cfg.Database.Password)
				assert.Equal(t, "customdb", cfg.Database.DBName)
				assert.Equal(t, "require", cfg.Database.SSLMode)
				assert.Equal(t, "custom-secret", cfg.Auth.JWTSecret)
				assert.Equal(t, 24*time.Hour, cfg.Auth.JWTExpiration)
				assert.Equal(t, 10, cfg.Auth.BCryptCost)
			},
		},
		{
			name: "missing JWT secret",
			setupEnv: func() {
				os.Clearenv()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.checkConfig != nil {
				tt.checkConfig(t, cfg)
			}
		})
	}
}

func TestConnectionString(t *testing.T) {
	tests := []struct {
		name string
		cfg  DatabaseConfig
		want string
	}{
		{
			name: "with password",
			cfg: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "without password",
			cfg: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=testuser dbname=testdb sslmode=disable",
		},
		{
			name: "with SSL",
			cfg: DatabaseConfig{
				Host:     "dbserver",
				Port:     5433,
				User:     "admin",
				Password: "secret",
				DBName:   "production",
				SSLMode:  "require",
			},
			want: "host=dbserver port=5433 user=admin password=secret dbname=production sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.ConnectionString()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "env var set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "env var not set",
			key:          "NONEXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				require.NoError(t, os.Setenv(tt.key, tt.envValue))
				defer func() { require.NoError(t, os.Unsetenv(tt.key)) }()
			}

			got := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		want         int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT",
			defaultValue: 123,
			envValue:     "456",
			want:         456,
		},
		{
			name:         "invalid integer",
			key:          "TEST_INT",
			defaultValue: 123,
			envValue:     "not-a-number",
			want:         123,
		},
		{
			name:         "env var not set",
			key:          "NONEXISTENT_INT",
			defaultValue: 999,
			envValue:     "",
			want:         999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				require.NoError(t, os.Setenv(tt.key, tt.envValue))
				defer func() { require.NoError(t, os.Unsetenv(tt.key)) }()
			}

			got := getEnvInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		want         time.Duration
	}{
		{
			name:         "valid duration",
			key:          "TEST_DURATION",
			defaultValue: 1 * time.Hour,
			envValue:     "30m",
			want:         30 * time.Minute,
		},
		{
			name:         "invalid duration",
			key:          "TEST_DURATION",
			defaultValue: 1 * time.Hour,
			envValue:     "not-a-duration",
			want:         1 * time.Hour,
		},
		{
			name:         "env var not set",
			key:          "NONEXISTENT_DURATION",
			defaultValue: 24 * time.Hour,
			envValue:     "",
			want:         24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				require.NoError(t, os.Setenv(tt.key, tt.envValue))
				defer func() { require.NoError(t, os.Unsetenv(tt.key)) }()
			}

			got := getEnvDuration(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}
