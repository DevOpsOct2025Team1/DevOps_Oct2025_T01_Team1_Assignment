//go:build integration

package store

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestDB(t *testing.T) (*UserStore, func()) {
	t.Helper()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database("test_user_service")
	store := NewUserStore(db)

	cleanup := func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}

	return store, cleanup
}

func TestListUsers_NoFilters(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, err := store.CreateUser(ctx, &User{
		Username:       "admin1",
		HashedPassword: "hash1",
		Role:           "admin",
	})
	if err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	_, err = store.CreateUser(ctx, &User{
		Username:       "user1",
		HashedPassword: "hash2",
		Role:           "user",
	})
	if err != nil {
		t.Fatalf("failed to create regular user: %v", err)
	}

	users, err := store.ListUsers(ctx, "", "")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestListUsers_RoleFilter(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "admin1", HashedPassword: "hash1", Role: "admin"})
	_, _ = store.CreateUser(ctx, &User{Username: "admin2", HashedPassword: "hash2", Role: "admin"})
	_, _ = store.CreateUser(ctx, &User{Username: "user1", HashedPassword: "hash3", Role: "user"})

	users, err := store.ListUsers(ctx, "admin", "")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 admin users, got %d", len(users))
	}

	for _, user := range users {
		if user.Role != "admin" {
			t.Errorf("expected role admin, got %s", user.Role)
		}
	}
}

func TestListUsers_UsernameFilter(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "john_admin", HashedPassword: "hash1", Role: "admin"})
	_, _ = store.CreateUser(ctx, &User{Username: "john_user", HashedPassword: "hash2", Role: "user"})
	_, _ = store.CreateUser(ctx, &User{Username: "alice", HashedPassword: "hash3", Role: "user"})

	users, err := store.ListUsers(ctx, "", "john")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users matching 'john', got %d", len(users))
	}
}

func TestListUsers_CombinedFilters(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "john_admin", HashedPassword: "hash1", Role: "admin"})
	_, _ = store.CreateUser(ctx, &User{Username: "john_user", HashedPassword: "hash2", Role: "user"})
	_, _ = store.CreateUser(ctx, &User{Username: "alice_admin", HashedPassword: "hash3", Role: "admin"})

	users, err := store.ListUsers(ctx, "admin", "john")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user (admin named john), got %d", len(users))
	}

	if len(users) > 0 && users[0].Username != "john_admin" {
		t.Errorf("expected john_admin, got %s", users[0].Username)
	}
}

func TestListUsers_NoMatches(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "alice", HashedPassword: "hash1", Role: "user"})

	users, err := store.ListUsers(ctx, "", "nonexistent")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	id, err := store.CreateUser(ctx, &User{
		Username:       "testuser",
		HashedPassword: "password",
		Role:           "user",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := store.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}

	if user.Id != id {
		t.Errorf("expected id %s, got %s", id, user.Id)
	}
}

func TestDeleteUserByID(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	id, err := store.CreateUser(ctx, &User{
		Username:       "testuser_delete",
		HashedPassword: "password",
		Role:           "user",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	err = store.DeleteUserByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	_, err = store.GetUserByID(ctx, id)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
