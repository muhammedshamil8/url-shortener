package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if hashedPassword == password {
		t.Fatalf("password was not hashed")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}
}

func TestGenerateToken(t *testing.T) {
	userID := 1
	email := "shamil@example.com"
	secret := "test-secret"
	expiry := "15m"

	token, err := GenerateToken(userID, email, secret, expiry)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}
	if token == "" {
		t.Fatalf("JWT token is empty")
	}
}

func TestValidateToken(t *testing.T) {
	userID := 1
	email := "shamil@example.com"
	secret := "test-secret"
	expiry := "15m"

	token, err := GenerateToken(userID, email, secret, expiry)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	decodedClaims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	decodedUserID := int(decodedClaims["user_id"].(float64))
	if decodedUserID != userID {
		t.Fatalf("user_id does not match: got %v, want %v", decodedUserID, userID)
	}

	if decodedClaims["email"] != email {
		t.Fatalf("email does not match: got %v, want %v", decodedClaims["email"], email)
	}
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
	userID := 1
	email := "shamil@example.com"
	secret := "test-secret"
	expiry := "15m"

	token, err := GenerateToken(userID, email, secret, expiry)
	if err != nil {
		t.Fatalf("failed to generate JWT token: %v", err)
	}

	_, err = ValidateToken(token, "invalid-secret")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}
