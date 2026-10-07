package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// These tests cover validation paths that return before any database call,
// so the handlers can be constructed with a nil *queries.Queries.

func withUserID(r *http.Request, userID int64) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "userID", userID))
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	return body.Error
}

func TestRespondWithJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	respondWithJSON(rec, http.StatusCreated, map[string]string{"status": "ok"})

	if rec.Code != http.StatusCreated {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("got Content-Type %q, want application/json", ct)
	}
	if body := strings.TrimSpace(rec.Body.String()); body != `{"status":"ok"}` {
		t.Errorf("got body %s", body)
	}
}

func TestRespondWithError(t *testing.T) {
	rec := httptest.NewRecorder()

	respondWithError(rec, http.StatusBadRequest, "boom")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if msg := decodeError(t, rec); msg != "boom" {
		t.Errorf("got error %q, want %q", msg, "boom")
	}
}

func TestHandlersRejectInvalidJSON(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		auth    bool
	}{
		{"create user", HandlerCreateUser(nil), false},
		{"login", HandleLogin(nil, "secret"), false},
		{"create chat", HandleCreateChat(nil), true},
		{"create message", HandleCreateMessage(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{not json"))
			if tt.auth {
				req = withUserID(req, 1)
			}
			rec := httptest.NewRecorder()

			tt.handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestCreateMessageRequiresContent(t *testing.T) {
	body := strings.NewReader(`{"chat_id": 1, "content": ""}`)
	req := withUserID(httptest.NewRequest(http.MethodPost, "/v1/messages", body), 1)
	rec := httptest.NewRecorder()

	HandleCreateMessage(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if msg := decodeError(t, rec); msg != "content is required" {
		t.Errorf("got error %q", msg)
	}
}

func TestListMessagesRejectsInvalidChatID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/v1/messages/{chatID}", func(w http.ResponseWriter, req *http.Request) {
		HandleListMessages(nil).ServeHTTP(w, withUserID(req, 1))
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/messages/abc", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
