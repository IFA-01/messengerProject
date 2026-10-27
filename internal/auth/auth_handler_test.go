package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

const testSecret = "test-secret"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "correct-horse" {
		t.Fatal("hash must not equal the plain password")
	}

	if err := CheckPassword("correct-horse", hash); err != nil {
		t.Errorf("CheckPassword with correct password returned error: %v", err)
	}
	if err := CheckPassword("wrong-password", hash); err == nil {
		t.Error("CheckPassword with wrong password returned nil error")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	if _, err := HashPassword(""); err == nil {
		t.Error("expected error for empty password")
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	if err := CheckPassword("", "some-hash"); err == nil {
		t.Error("expected error for empty password")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	token, err := GenerateToken(42, testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	userID, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if userID != 42 {
		t.Errorf("got userID %d, want 42", userID)
	}
}

func TestValidateTokenRejectsInvalid(t *testing.T) {
	validToken, err := GenerateToken(1, testSecret)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	noneToken, err := jwt.NewWithClaims(jwt.SigningMethodNone, &Claims{UserID: 1}).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to build unsigned token: %v", err)
	}

	tests := []struct {
		name   string
		token  string
		secret string
	}{
		{"wrong secret", validToken, "another-secret"},
		{"garbage", "not-a-jwt", testSecret},
		{"empty", "", testSecret},
		{"unsigned token", noneToken, testSecret},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateToken(tt.token, tt.secret); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
