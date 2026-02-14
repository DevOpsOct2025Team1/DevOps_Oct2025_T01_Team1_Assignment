//go:build integration

package store

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestContainer(t *testing.T) (*UserStore, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:8",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Fatalf("failed to get endpoint: %v", err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	db := client.Database("test_store")
	store := NewUserStore(db)

	cleanup := func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
		_ = container.Terminate(context.Background())
	}

	return store, cleanup
}

func TestCreateUser_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()

	id, err := store.CreateUser(context.Background(), &User{
		Username: "alice", HashedPassword: "hash", Role: "user",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestCreateUser_Duplicate_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "alice", HashedPassword: "hash", Role: "user"})
	_, err := store.CreateUser(ctx, &User{Username: "alice", HashedPassword: "hash2", Role: "user"})
	if err != ErrUserExists {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestGetUserByID_NotFound_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()

	_, err := store.GetUserByID(context.Background(), "000000000000000000000000")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetUserByID_InvalidHex_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()

	_, err := store.GetUserByID(context.Background(), "not-a-hex-id")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound for invalid hex, got %v", err)
	}
}

func TestGetUserByUsername_NotFound_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()

	_, err := store.GetUserByUsername(context.Background(), "nonexistent")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteUserByID_NotFound_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()

	err := store.DeleteUserByID(context.Background(), "000000000000000000000000")
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestEnsureDefaultAdmin_CreatesAdmin_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	err := store.EnsureDefaultAdmin(ctx, "admin", "password123")
	if err != nil {
		t.Fatalf("EnsureDefaultAdmin failed: %v", err)
	}

	user, err := store.GetUserByUsername(ctx, "admin")
	if err != nil {
		t.Fatalf("expected admin to exist: %v", err)
	}
	if user.Role != "admin" {
		t.Errorf("expected role admin, got %s", user.Role)
	}
}

func TestEnsureDefaultAdmin_SkipsWhenUsersExist_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "existing", HashedPassword: "hash", Role: "user"})
	err := store.EnsureDefaultAdmin(ctx, "admin", "password123")
	if err != nil {
		t.Fatalf("EnsureDefaultAdmin failed: %v", err)
	}

	_, err = store.GetUserByUsername(ctx, "admin")
	if err != ErrUserNotFound {
		t.Fatal("expected admin NOT to be created")
	}
}

func TestGetUserByID_Success_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	id, err := store.CreateUser(ctx, &User{
		Username: "testuser", HashedPassword: "hash", Role: "user",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user, err := store.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}
	if user.Role != "user" {
		t.Errorf("expected role user, got %s", user.Role)
	}
}

func TestGetUserByUsername_Success_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	_, err := store.CreateUser(ctx, &User{
		Username: "alice", HashedPassword: "hash", Role: "admin",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user, err := store.GetUserByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("expected username alice, got %s", user.Username)
	}
	if user.Role != "admin" {
		t.Errorf("expected role admin, got %s", user.Role)
	}
}

func TestDeleteUserByID_Success_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	id, err := store.CreateUser(ctx, &User{
		Username: "todelete", HashedPassword: "hash", Role: "user",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	err = store.DeleteUserByID(ctx, id)
	if err != nil {
		t.Fatalf("DeleteUserByID failed: %v", err)
	}

	_, err = store.GetUserByID(ctx, id)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestListUsers_Success_Integration(t *testing.T) {
	store, cleanup := setupTestContainer(t)
	defer cleanup()
	ctx := context.Background()

	_, _ = store.CreateUser(ctx, &User{Username: "admin1", HashedPassword: "hash1", Role: "admin"})
	_, _ = store.CreateUser(ctx, &User{Username: "user1", HashedPassword: "hash2", Role: "user"})
	_, _ = store.CreateUser(ctx, &User{Username: "user2", HashedPassword: "hash3", Role: "user"})

	users, err := store.ListUsers(ctx, "", "")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}
