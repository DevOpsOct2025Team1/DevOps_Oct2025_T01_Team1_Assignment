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
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

type mockUserClient struct {
	deleteAccountFunc func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error)
}

func (m *mockUserClient) DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	if m.deleteAccountFunc != nil {
		return m.deleteAccountFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) Close() error {
	return nil
}

func setupUserTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/api/admin/delete_user", handler.DeleteAccount)
	return router
}

func makeUserRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		if err != nil {
			t.Fatalf("failed to encode request body: %v", err)
		}
	}

	req, _ := http.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestDeleteAccount_Success(t *testing.T) {
	mock := &mockUserClient{
		deleteAccountFunc: func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
			return &userv1.DeleteUserByIdResponse{
				Success: true,
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "69654eb7a1135a809430d0b7",
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["success"] != true {
		t.Errorf("expected success true, got %v", response["success"])
	}
}

func TestDeleteAccount_MissingId(t *testing.T) {
	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_EmptyId(t *testing.T) {
	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_InvalidJSON(t *testing.T) {
	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	req, _ := http.NewRequest("DELETE", "/api/admin/delete_user", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_ServiceError(t *testing.T) {
	mock := &mockUserClient{
		deleteAccountFunc: func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "nonexistent-id",
	})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteAccount_ServiceReturnsFalse(t *testing.T) {
	mock := &mockUserClient{
		deleteAccountFunc: func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
			return &userv1.DeleteUserByIdResponse{
				Success: false,
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "69654eb7a1135a809430d0b7",
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["success"] != false {
		t.Errorf("expected success false, got %v", response["success"])
	}
}
