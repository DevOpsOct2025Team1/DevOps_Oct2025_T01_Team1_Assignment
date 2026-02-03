package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndValidate(t *testing.T) {
	manager := NewJWTManager("testsecret", time.Hour)

	token, err := manager.Generate("u1", "alice", "USER")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if claims.UserID != "u1" {
		t.Fatalf("expected UserID u1, got %s", claims.UserID)
	}

	if claims.Username != "alice" {
		t.Fatalf("expected username alice, got %s", claims.Username)
	}

	if claims.Role != "USER" {
		t.Fatalf("expected role USER, got %s", claims.Role)
	}
}

func TestValidate_InvalidToken(t *testing.T) {
	manager := NewJWTManager("testsecret", time.Hour)

	_, err := manager.Validate("not.a.real.token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
