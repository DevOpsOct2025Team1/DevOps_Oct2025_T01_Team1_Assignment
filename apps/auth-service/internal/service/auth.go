package service

import (
	"context"

	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/client"
	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/jwt"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	authv1.UnimplementedAuthServiceServer
	userClient client.UserClient
	jwtManager *jwt.Manager
}

func NewAuthServiceServer(userClient client.UserClient, jwtManager *jwt.Manager) *AuthServiceServer {
	return &AuthServiceServer{
		userClient: userClient,
		jwtManager: jwtManager,
	}
}

func (s *AuthServiceServer) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	user, err := s.userClient.CreateUser(ctx, req.Username, string(hashedPassword), userv1.Role_ROLE_USER)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtManager.Generate(user.Id, user.Username, user.Role.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &authv1.SignUpResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, valid, err := s.userClient.VerifyPassword(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := s.jwtManager.Generate(user.Id, user.Username, user.Role.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &authv1.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := s.jwtManager.Validate(req.Token)
	if err != nil {
		return &authv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}
	return &authv1.ValidateTokenResponse{
		Valid: true,
		User: &userv1.User{
			Id:       claims.UserID,
			Username: claims.Username,
			Role:     stringToRole(claims.Role),
		},
	}, nil
}

func stringToRole(role string) userv1.Role {
	switch role {
	case "ROLE_ADMIN":
		return userv1.Role_ROLE_ADMIN
	case "ROLE_USER":
		return userv1.Role_ROLE_USER
	default:
		return userv1.Role_ROLE_USER
	}
}
