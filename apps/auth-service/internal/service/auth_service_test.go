package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/provsalt/DOP_P01_Team1/auth-service/internal/jwt"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

// --------------------
// Mock User Client
// --------------------

type mockUserClient struct{}

func (m *mockUserClient) CreateUser(ctx context.Context, username, password string, role userv1.Role) (*userv1.User, error) {
	if username == "duplicate" {
		return nil, errors.New("user already exists")
	}
	return &userv1.User{
		Id:       "u1",
		Username: username,
		Role:     role,
	}, nil
}

func (m *mockUserClient) VerifyPassword(ctx context.Context, username, password string) (*userv1.User, bool, error) {
	if username != "testuser" || password != "password123" {
		return nil, false, errors.New("invalid")
	}

	return &userv1.User{
		Id:       "u1",
		Username: "testuser",
		Role:     userv1.Role_ROLE_USER,
	}, true, nil
}

func (m *mockUserClient) Close() error {
	return nil
}

// --------------------
// Helpers
// --------------------

func setupAuthService() *AuthServiceServer {
	jwtManager := jwt.NewJWTManager("testsecret", time.Hour)
	return NewAuthServiceServer(&mockUserClient{}, jwtManager)
}

// --------------------
// SignUp Tests
// --------------------

func TestSignUp_Success(t *testing.T) {
	svc := setupAuthService()

	resp, err := svc.SignUp(context.Background(), &authv1.SignUpRequest{
		Username: "alice",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.User == nil {
		t.Fatal("expected user, got nil")
	}
	if resp.Token == "" {
		t.Fatal("expected token, got empty")
	}
}

func TestSignUp_MissingUsername(t *testing.T) {
	svc := setupAuthService()

	_, err := svc.SignUp(context.Background(), &authv1.SignUpRequest{
		Password: "secret",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSignUp_MissingPassword(t *testing.T) {
	svc := setupAuthService()

	_, err := svc.SignUp(context.Background(), &authv1.SignUpRequest{
		Username: "bob",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

// --------------------
// Login Tests
// --------------------

func TestLogin_Success(t *testing.T) {
	svc := setupAuthService()

	resp, err := svc.Login(context.Background(), &authv1.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	svc := setupAuthService()

	_, err := svc.Login(context.Background(), &authv1.LoginRequest{
		Username: "wrong",
		Password: "wrong",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := setupAuthService()

	_, err := svc.Login(context.Background(), &authv1.LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
	})

	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}


// --------------------
// Validate Token Tests
// --------------------

func TestValidateToken_Valid(t *testing.T) {
	svc := setupAuthService()

	loginResp, _ := svc.Login(context.Background(), &authv1.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	resp, err := svc.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: loginResp.Token,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Valid {
		t.Fatal("expected token to be valid")
	}
	if resp.User == nil {
		t.Fatal("expected user")
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := setupAuthService()

	resp, err := svc.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: "badtoken",
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}
	if resp.Valid {
		t.Fatal("expected invalid token")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	svc := setupAuthService()

	_, err := svc.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: "",
	})

	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	jwtManager := jwt.NewJWTManager("testsecret", -1*time.Hour) // expired
	svc := NewAuthServiceServer(&mockUserClient{}, jwtManager)

	token, _ := jwtManager.Generate("u1", "user", "ROLE_USER")

	resp, err := svc.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}
	if resp.Valid {
		t.Fatal("expected token to be invalid")
	}
}
