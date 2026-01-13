package service

import (
	"context"
	"errors"

	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"github.com/provsalt/DOP_P01_Team1/user-service/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	store *store.UserStore

	userv1.UnimplementedUserServiceServer
}

func NewUserServiceServer(store *store.UserStore) *UserServiceServer {
	return &UserServiceServer{
		store: store,
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.HashedPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "hashed_password is required")
	}

	role := req.Role
	if role == userv1.Role_ROLE_UNSPECIFIED {
		role = userv1.Role_ROLE_USER
	}

	user := &store.User{
		Username:       req.Username,
		HashedPassword: req.HashedPassword,
		Role:           roleToString(role),
	}

	id, err := s.store.CreateUser(user)
	if err != nil {
		if errors.Is(err, store.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	user.Id = id

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:       user.Id,
			Username: user.Username,
			Role:     role,
		},
	}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.store.GetUserByID(req.Id)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:       user.Id,
			Username: user.Username,
			Role:     stringToRole(user.Role),
		},
	}, nil
}

func (s *UserServiceServer) GetUserByUsername(ctx context.Context, req *userv1.GetUserByUsernameRequest) (*userv1.GetUserByUsernameResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	user, err := s.store.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &userv1.GetUserByUsernameResponse{
		User: &userv1.User{
			Id:       user.Id,
			Username: user.Username,
			Role:     stringToRole(user.Role),
		},
		HashedPassword: user.HashedPassword,
	}, nil
}

func roleToString(role userv1.Role) string {
	switch role {
	case userv1.Role_ROLE_ADMIN:
		return "admin"
	case userv1.Role_ROLE_USER:
		return "user"
	default:
		return "user"
	}
}

func stringToRole(role string) userv1.Role {
	switch role {
	case "admin":
		return userv1.Role_ROLE_ADMIN
	case "user":
		return userv1.Role_ROLE_USER
	default:
		return userv1.Role_ROLE_USER
	}
}
