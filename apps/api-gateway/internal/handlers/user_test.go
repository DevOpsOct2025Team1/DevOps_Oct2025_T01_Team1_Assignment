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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUserClient struct {
	getUserFunc       func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	deleteAccountFunc func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error)
	listUsersFunc     func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error)
}

func (m *mockUserClient) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) DeleteAccount(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	if m.deleteAccountFunc != nil {
		return m.deleteAccountFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) Close() error {
	return nil
}

func setupUserTestRouter(handler *UserHandler, currentUser *userv1.User) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		if currentUser != nil {
			c.Set("user", currentUser)
		}
		c.Next()
	})
	router.DELETE("/api/admin/delete_user", handler.DeleteUser)
	router.GET("/api/admin/list_users", handler.ListUsers)
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
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		getUserFunc: func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
			return &userv1.GetUserResponse{
				User: &userv1.User{
					Id:       req.Id,
					Username: "regular-user",
					Role:     userv1.Role_ROLE_USER,
				},
			}, nil
		},
		deleteAccountFunc: func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
			return &userv1.DeleteUserByIdResponse{
				Success: true,
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

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
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_EmptyId(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "",
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_InvalidJSON(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("DELETE", "/api/admin/delete_user", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_UserNotFound(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		getUserFunc: func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
			return nil, status.Error(codes.NotFound, "user not found")
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "nonexistent-id",
	})

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteAccount_SelfDeletion(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "admin",
	})

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "cannot delete your own account" {
		t.Errorf("expected error 'cannot delete your own account', got %v", response["error"])
	}
}

func TestDeleteAccount_AdminDeletion(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		getUserFunc: func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
			return &userv1.GetUserResponse{
				User: &userv1.User{
					Id:       "admin2",
					Username: "admin2",
					Role:     userv1.Role_ROLE_ADMIN,
				},
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "admin2",
	})

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response["error"] != "cannot delete admin accounts" {
		t.Errorf("expected error 'cannot delete admin accounts', got %v", response["error"])
	}
}

func TestDeleteAccount_NoUserInContext(t *testing.T) {
	mock := &mockUserClient{}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, nil)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{
		"id": "user",
	})

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListUsers_Success_NoFilters(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			if req.Role != userv1.Role_ROLE_UNSPECIFIED || req.UsernameFilter != "" {
				t.Errorf("expected no filters, got role=%v, username=%s", req.Role, req.UsernameFilter)
			}
			return &userv1.ListUsersResponse{
				Users: []*userv1.User{
					{Id: "1", Username: "admin1", Role: userv1.Role_ROLE_ADMIN},
					{Id: "2", Username: "user1", Role: userv1.Role_ROLE_USER},
				},
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	users, ok := response["users"].([]interface{})
	if !ok {
		t.Fatalf("expected users array, got %T", response["users"])
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestListUsers_WithRoleFilter(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			if req.Role != userv1.Role_ROLE_ADMIN {
				t.Errorf("expected role=ROLE_ADMIN, got %v", req.Role)
			}
			return &userv1.ListUsersResponse{
				Users: []*userv1.User{
					{Id: "1", Username: "admin1", Role: userv1.Role_ROLE_ADMIN},
				},
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users?role=admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListUsers_WithUsernameFilter(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			if req.UsernameFilter != "john" {
				t.Errorf("expected username=john, got %s", req.UsernameFilter)
			}
			return &userv1.ListUsersResponse{
				Users: []*userv1.User{
					{Id: "1", Username: "john_admin", Role: userv1.Role_ROLE_ADMIN},
					{Id: "2", Username: "john_user", Role: userv1.Role_ROLE_USER},
				},
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users?username=john", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListUsers_WithBothFilters(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			if req.Role != userv1.Role_ROLE_USER || req.UsernameFilter != "test" {
				t.Errorf("expected role=ROLE_USER and username=test, got role=%v, username=%s", req.Role, req.UsernameFilter)
			}
			return &userv1.ListUsersResponse{
				Users: []*userv1.User{
					{Id: "1", Username: "test_user", Role: userv1.Role_ROLE_USER},
				},
			}, nil
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users?role=user&username=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestListUsers_InvalidRole(t *testing.T) {
	currentUser := &userv1.User{
		Id:       "admin",
		Username: "admin",
		Role:     userv1.Role_ROLE_ADMIN,
	}

	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			return nil, status.Error(codes.InvalidArgument, "role must be 'admin' or 'user'")
		},
	}

	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users?role=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_GetUserInvalidArgument(t *testing.T) {
	currentUser := &userv1.User{Id: "admin", Username: "admin", Role: userv1.Role_ROLE_ADMIN}
	mock := &mockUserClient{
		getUserFunc: func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		},
	}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{"id": "bad-id"})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccount_DeleteAccountNotFound(t *testing.T) {
	currentUser := &userv1.User{Id: "admin", Username: "admin", Role: userv1.Role_ROLE_ADMIN}
	mock := &mockUserClient{
		getUserFunc: func(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
			return &userv1.GetUserResponse{
				User: &userv1.User{Id: req.Id, Username: "user", Role: userv1.Role_ROLE_USER},
			}, nil
		},
		deleteAccountFunc: func(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
			return nil, status.Error(codes.NotFound, "not found")
		},
	}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	w := makeUserRequest(t, router, "DELETE", "/api/admin/delete_user", map[string]string{"id": "some-id"})
	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListUsers_ServiceError(t *testing.T) {
	currentUser := &userv1.User{Id: "admin", Username: "admin", Role: userv1.Role_ROLE_ADMIN}
	mock := &mockUserClient{
		listUsersFunc: func(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
			return nil, status.Error(codes.Internal, "db error")
		},
	}
	handler := NewUserHandler(mock)
	router := setupUserTestRouter(handler, currentUser)

	req, _ := http.NewRequest("GET", "/api/admin/list_users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
