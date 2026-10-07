package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IFA-01/messenger/internal/auth"
)

const testSecret = "test-secret"

func TestAuthMiddlewareRejectsBadRequests(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"missing header", ""},
		{"wrong scheme", "Basic abc"},
		{"no token", "Bearer"},
		{"invalid token", "Bearer not-a-jwt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
			})

			req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()

			AuthMiddleWare(testSecret)(next).ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
			}
			if nextCalled {
				t.Error("next handler must not be called")
			}
		})
	}
}

func TestAuthMiddlewarePassesUserID(t *testing.T) {
	token, err := auth.GenerateToken(7, testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	var gotUserID int64
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = r.Context().Value("userID").(int64)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	AuthMiddleWare(testSecret)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if gotUserID != 7 {
		t.Errorf("got userID %d, want 7", gotUserID)
	}
}
