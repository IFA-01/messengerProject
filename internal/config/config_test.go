package config

import (
	"errors"
	"testing"
)

const validSecret = "0123456789abcdef0123456789abcdef"

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr error
	}{
		{
			name: "valid config",
			env:  map[string]string{"DATABASE_URL": "testdb", "JWT_SECRET": validSecret},
		},
		{
			name:    "missing DATABASE_URL",
			env:     map[string]string{"DATABASE_URL": "", "JWT_SECRET": validSecret},
			wantErr: ErrMissingDatabaseURL,
		},
		{
			name:    "missing JWT_SECRET",
			env:     map[string]string{"DATABASE_URL": "testdb", "JWT_SECRET": ""},
			wantErr: ErrWeakJWTSecret,
		},
		{
			name:    "short JWT_SECRET",
			env:     map[string]string{"DATABASE_URL": "testdb", "JWT_SECRET": "testsecret"},
			wantErr: ErrWeakJWTSecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVER_PORT", "")
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.DatabaseURL != "testdb" {
				t.Errorf("got DATABASE_URL %q, want testdb", cfg.DatabaseURL)
			}
			if cfg.JWTSecret != validSecret {
				t.Errorf("got JWT_SECRET %q, want %q", cfg.JWTSecret, validSecret)
			}
			if cfg.ServerPort != "8080" {
				t.Errorf("got SERVER_PORT %q, want 8080", cfg.ServerPort)
			}
		})
	}
}
