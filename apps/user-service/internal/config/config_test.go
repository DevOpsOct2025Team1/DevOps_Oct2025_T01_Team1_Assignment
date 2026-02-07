package config

import (
	"os"
	"testing"
)

func unsetenv(t *testing.T, key string) {
	t.Helper()
	prev, wasSet := os.LookupEnv(key)
	_ = os.Unsetenv(key)
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv(key, prev)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func TestLoad_RequiresMongoURI(t *testing.T) {
	unsetenv(t, "PORT")
	unsetenv(t, "MONGODB_URI")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error when MONGODB_URI is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	unsetenv(t, "PORT")
	unsetenv(t, "MONGODB_DATABASE")
	unsetenv(t, "ENVIRONMENT")

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
