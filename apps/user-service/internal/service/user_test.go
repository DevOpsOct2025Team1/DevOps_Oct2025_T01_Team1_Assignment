package service

import (
	"context"
	"errors"
	"testing"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser_Validation(t *testing.T) {
	srv := NewUserServiceServer(nil)

	_, err := srv.CreateUser(context.Background(), &userv1.CreateUserRequest{HashedPassword: "x"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}

	_, err = srv.CreateUser(context.Background(), &userv1.CreateUserRequest{Username: "u"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestGetUser_Validation(t *testing.T) {
	srv := NewUserServiceServer(nil)
	_, err := srv.GetUser(context.Background(), &userv1.GetUserRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestGetUserByUsername_Validation(t *testing.T) {
	srv := NewUserServiceServer(nil)
	_, err := srv.GetUserByUsername(context.Background(), &userv1.GetUserByUsernameRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestVerifyPassword_Validation(t *testing.T) {
	srv := NewUserServiceServer(nil)

	_, err := srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{Password: "p"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}

	_, err = srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{Username: "u"})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestDeleteUser_Validation(t *testing.T) {
	srv := NewUserServiceServer(nil)
	_, err := srv.DeleteUser(context.Background(), &userv1.DeleteUserByIdRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestRoleConversions(t *testing.T) {
	if roleToString(userv1.Role_ROLE_ADMIN) != "admin" {
		t.Fatalf("expected admin")
	}
	if roleToString(userv1.Role_ROLE_USER) != "user" {
		t.Fatalf("expected user")
	}
	if stringToRole("admin") != userv1.Role_ROLE_ADMIN {
		t.Fatalf("expected ROLE_ADMIN")
	}
	if stringToRole("user") != userv1.Role_ROLE_USER {
		t.Fatalf("expected ROLE_USER")
	}
	if stringToRole("unknown") != userv1.Role_ROLE_USER {
		t.Fatalf("expected unknown to map to ROLE_USER")
	}
}

type mockUserStore struct {
	listUsersFunc func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error)
}

func (m *mockUserStore) CreateUser(ctx context.Context, user *store.User) (string, error) {
	return "", nil
}

func (m *mockUserStore) GetUserByID(ctx context.Context, id string) (*store.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	return nil, nil
}

func (m *mockUserStore) DeleteUserByID(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserStore) ListUsers(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc(ctx, roleFilter, usernameFilter)
	}
	return nil, nil
}

func TestListUsers_Success_NoFilters(t *testing.T) {
	mockStore := &mockUserStore{
		listUsersFunc: func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error) {
			if roleFilter != "" {
				t.Errorf("expected empty roleFilter, got %q", roleFilter)
			}
			if usernameFilter != "" {
				t.Errorf("expected empty usernameFilter, got %q", usernameFilter)
			}

			return []*store.User{
				{Id: "1", Username: "user1", Role: "user", HashedPassword: "hash1"},
				{Id: "2", Username: "admin1", Role: "admin", HashedPassword: "hash2"},
			}, nil
		},
	}

	srv := NewUserServiceServer(mockStore)
	resp, err := srv.ListUsers(context.Background(), &userv1.ListUsersRequest{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(resp.Users))
	}

	if resp.Users[0].Id != "1" || resp.Users[0].Username != "user1" || resp.Users[0].Role != userv1.Role_ROLE_USER {
		t.Errorf("first user mismatch: got %+v", resp.Users[0])
	}

	if resp.Users[1].Id != "2" || resp.Users[1].Username != "admin1" || resp.Users[1].Role != userv1.Role_ROLE_ADMIN {
		t.Errorf("second user mismatch: got %+v", resp.Users[1])
	}

}

func TestListUsers_InvalidRole(t *testing.T) {
	mockStore := &mockUserStore{}
	srv := NewUserServiceServer(mockStore)

	resp, err := srv.ListUsers(context.Background(), &userv1.ListUsersRequest{
		Role: userv1.Role(999),
	})

	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestListUsers_WithRoleFilter(t *testing.T) {
	mockStore := &mockUserStore{
		listUsersFunc: func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error) {
			if roleFilter != "admin" {
				t.Errorf("expected roleFilter 'admin', got %q", roleFilter)
			}
			if usernameFilter != "" {
				t.Errorf("expected empty usernameFilter, got %q", usernameFilter)
			}

			return []*store.User{
				{Id: "2", Username: "admin1", Role: "admin", HashedPassword: "hash2"},
			}, nil
		},
	}

	srv := NewUserServiceServer(mockStore)
	resp, err := srv.ListUsers(context.Background(), &userv1.ListUsersRequest{
		Role: userv1.Role_ROLE_ADMIN,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(resp.Users))
	}

	if resp.Users[0].Role != userv1.Role_ROLE_ADMIN {
		t.Errorf("expected ROLE_ADMIN, got %v", resp.Users[0].Role)
	}
}

func TestListUsers_WithUsernameFilter(t *testing.T) {
	mockStore := &mockUserStore{
		listUsersFunc: func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error) {
			if roleFilter != "" {
				t.Errorf("expected empty roleFilter, got %q", roleFilter)
			}
			if usernameFilter != "john" {
				t.Errorf("expected usernameFilter 'john', got %q", usernameFilter)
			}

			return []*store.User{
				{Id: "3", Username: "john_doe", Role: "user", HashedPassword: "hash3"},
			}, nil
		},
	}

	srv := NewUserServiceServer(mockStore)
	resp, err := srv.ListUsers(context.Background(), &userv1.ListUsersRequest{
		UsernameFilter: "john",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(resp.Users))
	}

	if resp.Users[0].Username != "john_doe" {
		t.Errorf("expected username 'john_doe', got %q", resp.Users[0].Username)
	}
}

func TestListUsers_StoreError(t *testing.T) {
	mockStore := &mockUserStore{
		listUsersFunc: func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	srv := NewUserServiceServer(mockStore)
	resp, err := srv.ListUsers(context.Background(), &userv1.ListUsersRequest{})

	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
	}

	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal error, got %v", status.Code(err))
	}
}
