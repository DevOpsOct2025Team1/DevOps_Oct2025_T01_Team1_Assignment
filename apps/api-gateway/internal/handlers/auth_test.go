package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

type mockAuthClient struct {
	signUpFunc        func(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error)
	loginFunc         func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	validateTokenFunc func(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error)
}

func (m *mockAuthClient) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	if m.signUpFunc != nil {
		return m.signUpFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthClient) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthClient) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthClient) Close() error {
	return nil
}

func setupTestRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/signup", handler.SignUp)
	router.POST("/api/login", handler.Login)
	return router
}

func makeRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		if err != nil {
			t.Errorf("failed to encode request body: %v", err)
		}
	}

	req, _ := http.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSignUp_Success(t *testing.T) {
	mock := &mockAuthClient{
		signUpFunc: func(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
			return &authv1.SignUpResponse{
				User: &userv1.User{
					Id:       "69654eb7a1135a809430d0b7",
					Username: req.Username,
					Role:     userv1.Role_ROLE_USER,
				},
				Token: "jwt-token-123",
			}, nil
		},
	}

	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/signup", map[string]string{
		"username": "testing",
		"password": "password123",
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["token"] != "jwt-token-123" {
		t.Errorf("expected token 'jwt-token-123', got %v", response["token"])
	}

	user := response["user"].(map[string]interface{})
	if user["username"] != "testing" {
		t.Errorf("expected username 'testing', got %v", user["username"])
	}
}

func TestSignUp_MissingUsername(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/signup", map[string]string{
		"password": "password123",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSignUp_MissingPassword(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/signup", map[string]string{
		"username": "testing",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSignUp_EmptyBody(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/signup", map[string]string{})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSignUp_InvalidJSON(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSignUp_AuthServiceError(t *testing.T) {
	mock := &mockAuthClient{
		signUpFunc: func(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
			return nil, errors.New("user already exists")
		},
	}

	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/signup", map[string]string{
		"username": "testing",
		"password": "password123",
	})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "user already exists" {
		t.Errorf("expected error 'user already exists', got %v", response["error"])
	}
}

func TestLogin_Success(t *testing.T) {
	mock := &mockAuthClient{
		loginFunc: func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
			return &authv1.LoginResponse{
				User: &userv1.User{
					Id:       "69654eb7a1135a809430d0b7",
					Username: req.Username,
					Role:     userv1.Role_ROLE_USER,
				},
				Token: "jwt-token-456",
			}, nil
		},
	}

	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/login", map[string]string{
		"username": "testing",
		"password": "password123",
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	if response["token"] != "jwt-token-456" {
		t.Errorf("expected token 'jwt-token-456', got %v", response["token"])
	}
}

func TestLogin_MissingUsername(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/login", map[string]string{
		"password": "password123",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogin_MissingPassword(t *testing.T) {
	mock := &mockAuthClient{}
	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/login", map[string]string{
		"username": "testing",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &mockAuthClient{
		loginFunc: func(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
			return nil, errors.New("invalid credentials")
		},
	}

	handler := NewAuthHandler(mock)
	router := setupTestRouter(handler)

	w := makeRequest(t, router, "POST", "/api/login", map[string]string{
		"username": "testing",
		"password": "wrongpassword",
	})

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
