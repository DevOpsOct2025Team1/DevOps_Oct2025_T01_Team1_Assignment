package service

import (
	"context"
	"errors"
	"testing"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/store"
	"golang.org/x/crypto/bcrypt"
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
	createUserFunc        func(ctx context.Context, user *store.User) (string, error)
	getUserByIDFunc       func(ctx context.Context, id string) (*store.User, error)
	getUserByUsernameFunc func(ctx context.Context, username string) (*store.User, error)
	deleteUserByIDFunc    func(ctx context.Context, id string) error
	listUsersFunc         func(ctx context.Context, roleFilter string, usernameFilter string) ([]*store.User, error)
}

func (m *mockUserStore) CreateUser(ctx context.Context, user *store.User) (string, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	return "", nil
}

func (m *mockUserStore) GetUserByID(ctx context.Context, id string) (*store.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserStore) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	if m.getUserByUsernameFunc != nil {
		return m.getUserByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (m *mockUserStore) DeleteUserByID(ctx context.Context, id string) error {
	if m.deleteUserByIDFunc != nil {
		return m.deleteUserByIDFunc(ctx, id)
	}
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

func TestCreateUser_Success(t *testing.T) {
	mockStore := &mockUserStore{
		createUserFunc: func(ctx context.Context, user *store.User) (string, error) {
			return "abc123", nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.CreateUser(context.Background(), &userv1.CreateUserRequest{
		Username:       "newuser",
		HashedPassword: "hashed",
		Role:           userv1.Role_ROLE_USER,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.User.Id != "abc123" {
		t.Errorf("expected id abc123, got %s", resp.User.Id)
	}
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	mockStore := &mockUserStore{
		createUserFunc: func(ctx context.Context, user *store.User) (string, error) {
			return "", store.ErrUserExists
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.CreateUser(context.Background(), &userv1.CreateUserRequest{
		Username:       "existing",
		HashedPassword: "hashed",
	})
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected AlreadyExists, got %v", status.Code(err))
	}
}

func TestCreateUser_InternalError(t *testing.T) {
	mockStore := &mockUserStore{
		createUserFunc: func(ctx context.Context, user *store.User) (string, error) {
			return "", errors.New("db down")
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.CreateUser(context.Background(), &userv1.CreateUserRequest{
		Username:       "user",
		HashedPassword: "hashed",
	})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %v", status.Code(err))
	}
}

func TestCreateUser_DefaultRole(t *testing.T) {
	var capturedRole string
	mockStore := &mockUserStore{
		createUserFunc: func(ctx context.Context, user *store.User) (string, error) {
			capturedRole = user.Role
			return "id1", nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.CreateUser(context.Background(), &userv1.CreateUserRequest{
		Username:       "user",
		HashedPassword: "hashed",
		Role:           userv1.Role_ROLE_UNSPECIFIED,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedRole != "user" {
		t.Errorf("expected role 'user', got %q", capturedRole)
	}
	if resp.User.Role != userv1.Role_ROLE_USER {
		t.Errorf("expected ROLE_USER in response, got %v", resp.User.Role)
	}
}

func TestGetUser_Success(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByIDFunc: func(ctx context.Context, id string) (*store.User, error) {
			return &store.User{Id: id, Username: "found", Role: "admin"}, nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.GetUser(context.Background(), &userv1.GetUserRequest{Id: "u1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Username != "found" {
		t.Errorf("expected 'found', got %q", resp.User.Username)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByIDFunc: func(ctx context.Context, id string) (*store.User, error) {
			return nil, store.ErrUserNotFound
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.GetUser(context.Background(), &userv1.GetUserRequest{Id: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", status.Code(err))
	}
}

func TestGetUser_InternalError(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByIDFunc: func(ctx context.Context, id string) (*store.User, error) {
			return nil, errors.New("db error")
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.GetUser(context.Background(), &userv1.GetUserRequest{Id: "u1"})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %v", status.Code(err))
	}
}

func TestGetUserByUsername_Success(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return &store.User{Id: "u1", Username: username, Role: "user"}, nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.GetUserByUsername(context.Background(), &userv1.GetUserByUsernameRequest{Username: "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Username != "alice" {
		t.Errorf("expected alice, got %s", resp.User.Username)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return nil, store.ErrUserNotFound
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.GetUserByUsername(context.Background(), &userv1.GetUserByUsernameRequest{Username: "nobody"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", status.Code(err))
	}
}

func TestGetUserByUsername_InternalError(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return nil, errors.New("db error")
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.GetUserByUsername(context.Background(), &userv1.GetUserByUsernameRequest{Username: "alice"})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %v", status.Code(err))
	}
}

func TestVerifyPassword_Valid(t *testing.T) {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return &store.User{Id: "u1", Username: username, HashedPassword: string(hashedPw), Role: "user"}, nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{
		Username: "alice", Password: "correct",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Valid {
		t.Fatal("expected valid=true")
	}
}

func TestVerifyPassword_Invalid(t *testing.T) {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return &store.User{Id: "u1", Username: username, HashedPassword: string(hashedPw), Role: "user"}, nil
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{
		Username: "alice", Password: "wrong",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Valid {
		t.Fatal("expected valid=false")
	}
}

func TestVerifyPassword_UserNotFound(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return nil, store.ErrUserNotFound
		},
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{
		Username: "nobody", Password: "any",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Valid {
		t.Fatal("expected valid=false")
	}
}

func TestVerifyPassword_InternalError(t *testing.T) {
	mockStore := &mockUserStore{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*store.User, error) {
			return nil, errors.New("db error")
		},
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.VerifyPassword(context.Background(), &userv1.VerifyPasswordRequest{
		Username: "alice", Password: "any",
	})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %v", status.Code(err))
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mockStore := &mockUserStore{
		deleteUserByIDFunc: func(ctx context.Context, id string) error { return nil },
	}
	srv := NewUserServiceServer(mockStore)
	resp, err := srv.DeleteUser(context.Background(), &userv1.DeleteUserByIdRequest{Id: "u1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockStore := &mockUserStore{
		deleteUserByIDFunc: func(ctx context.Context, id string) error { return store.ErrUserNotFound },
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.DeleteUser(context.Background(), &userv1.DeleteUserByIdRequest{Id: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", status.Code(err))
	}
}

func TestDeleteUser_InternalError(t *testing.T) {
	mockStore := &mockUserStore{
		deleteUserByIDFunc: func(ctx context.Context, id string) error { return errors.New("db error") },
	}
	srv := NewUserServiceServer(mockStore)
	_, err := srv.DeleteUser(context.Background(), &userv1.DeleteUserByIdRequest{Id: "u1"})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected Internal, got %v", status.Code(err))
	}
}
