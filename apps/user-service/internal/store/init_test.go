package store

import (
	"context"
	"testing"
)

func TestEnsureDefaultAdmin_EmptyCredentials(t *testing.T) {
	s := &UserStore{}

	if err := s.EnsureDefaultAdmin(context.Background(), "", "pass"); err == nil {
		t.Fatalf("expected error for empty username")
	}
	if err := s.EnsureDefaultAdmin(context.Background(), "user", ""); err == nil {
		t.Fatalf("expected error for empty password")
	}
}
