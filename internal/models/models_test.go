package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/IFA-01/messenger/internal/repository/queries"
)

func TestUserJSONHasNoPasswordHash(t *testing.T) {
	user := DatabaseUserToUser(queries.User{
		ID:           1,
		Username:     "ifa",
		Email:        "ifa@example.com",
		PasswordHash: "$2a$10$secret",
	})
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}
	body := string(data)
	if strings.Contains(body, "password") || strings.Contains(body, "$2a$") {
		t.Errorf("user JSON leaks the password hash: %s", body)
	}
}
