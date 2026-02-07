package config

import (
	"os"
	"testing"
)

func TestLoad_RequiresMongoURI(t *testing.T) {
	t.Setenv("PORT", "")
	os.Unsetenv("MONGODB_URI")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error when MONGODB_URI is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Unsetenv("PORT")
	os.Unsetenv("MONGODB_DATABASE")
	os.Unsetenv("ENVIRONMENT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.MongoDBDatabase != "user_service" {
		t.Fatalf("expected default db user_service, got %q", cfg.MongoDBDatabase)
	}
	if cfg.Environment != "development" {
		t.Fatalf("expected default environment development, got %q", cfg.Environment)
	}
}
