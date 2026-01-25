package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

// Compile-time check that mock implements the real interface
var _ handlers.AuthServiceClient = (*mockAuthService)(nil)

type mockAuthService struct {
	validateTokenFn func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error)
}

func (m *mockAuthService) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	return nil, errors.New("not used")
}
func (m *mockAuthService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return nil, errors.New("not used")
}
func (m *mockAuthService) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if m.validateTokenFn != nil {
		return m.validateTokenFn(ctx, req)
	}
	return nil, errors.New("not implemented")
}
func (m *mockAuthService) Close() error { return nil }

func setupProtectedRoute(authSvc handlers.AuthServiceClient, roles []userv1.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// A protected route that only passes if middleware calls c.Next()
	r.GET("/admin", ValidateRole(authSvc, roles), func(c *gin.Context) {
		// also check middleware set "user" in context
		if _, exists := c.Get("user"); !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not set in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	return r
}

func TestValidateRole_MissingAuthorizationHeader_Returns401(t *testing.T) {
	authSvc := &mockAuthService{}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidateRole_InvalidAuthType_Returns401(t *testing.T) {
	authSvc := &mockAuthService{}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Basic abc")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidateRole_InvalidHeaderFormat_Returns401(t *testing.T) {
	authSvc := &mockAuthService{}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer") // split length != 2
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidateRole_ValidateTokenError_Returns500(t *testing.T) {
	authSvc := &mockAuthService{
		validateTokenFn: func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
			return nil, errors.New("auth service down")
		},
	}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestValidateRole_WrongRole_Returns401(t *testing.T) {
	authSvc := &mockAuthService{
		validateTokenFn: func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
			return &authv1.ValidateTokenResponse{
				User: &userv1.User{
					Id:       "u1",
					Username: "normal",
					Role:     userv1.Role_ROLE_USER,
				},
			}, nil
		},
	}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer user-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidateRole_AdminRole_Allows200AndSetsUser(t *testing.T) {
	authSvc := &mockAuthService{
		validateTokenFn: func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
			// Assert token passed correctly into ValidateTokenRequest
			if req.Token != "admin-token" {
				return nil, errors.New("wrong token passed to auth service")
			}
			return &authv1.ValidateTokenResponse{
				User: &userv1.User{
					Id:       "a1",
					Username: "admin",
					Role:     userv1.Role_ROLE_ADMIN,
				},
			}, nil
		},
	}
	r := setupProtectedRoute(authSvc, []userv1.Role{userv1.Role_ROLE_ADMIN})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
