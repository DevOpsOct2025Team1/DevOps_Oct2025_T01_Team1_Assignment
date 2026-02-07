package service

import (
	"context"
	"testing"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
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
