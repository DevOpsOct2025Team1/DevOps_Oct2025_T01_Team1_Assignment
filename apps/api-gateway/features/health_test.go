package features

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/server"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	filev1 "github.com/provsalt/DOP_P01_Team1/common/file/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

type healthTestContext struct {
	server           *httptest.Server
	response         *http.Response
	responseBody     map[string]interface{}
	responseRaw      []byte
	responseTime     time.Duration
	authToken        string
	customAuthHeader string
}

func newHealthTestContext() *healthTestContext {
	mockAuthClient := &mockAuthClient{}
	mockUserClient := &mockUserClient{}
	mockFileClient := &mockFileClient{}

	cfg := &config.Config{Environment: "test"}
	srv := server.New(mockAuthClient, mockUserClient, mockFileClient, cfg)
	testServer := httptest.NewServer(srv.Router)

	return &healthTestContext{
		server: testServer,
	}
}

// mock clients are needed since health endpoint doesn't use them but server needs them for init
type mockAuthClient struct{}

func (m *mockAuthClient) SignUp(_ context.Context, _ *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAuthClient) Login(_ context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Username == "testuser" && req.Password == "password123" {
		return &authv1.LoginResponse{
			User:  &userv1.User{Id: "u1", Username: "testuser", Role: userv1.Role_ROLE_USER},
			Token: "user-token",
		}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (m *mockAuthClient) ValidateToken(_ context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	switch req.Token {
	case "admin-token":
		return &authv1.ValidateTokenResponse{
			Valid: true,
			User:  &userv1.User{Id: "a1", Username: "admin", Role: userv1.Role_ROLE_ADMIN},
		}, nil
	case "user-token":
		return &authv1.ValidateTokenResponse{
			Valid: true,
			User:  &userv1.User{Id: "u1", Username: "user", Role: userv1.Role_ROLE_USER},
		}, nil
	default:
		// simulate failed validation: invalid token => Valid=false and no User
		return &authv1.ValidateTokenResponse{
			Valid: false,
			User:  nil,
		}, nil
	}
}

func (m *mockAuthClient) Close() error {
	return nil
}

type mockUserClient struct{}

func (m *mockUserClient) GetUser(_ context.Context, _ *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// For BDD scenarios we just need a non-admin target user to exist.
	return &userv1.GetUserResponse{
		User: &userv1.User{Id: "u123", Username: "target", Role: userv1.Role_ROLE_USER},
	}, nil
}

func (m *mockUserClient) DeleteAccount(_ context.Context, _ *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	return &userv1.DeleteUserByIdResponse{Success: true}, nil
}	

func (m *mockUserClient) ListUsers(_ context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	return &userv1.ListUsersResponse{
		Users: []*userv1.User{
			{Id: "a1", Username: "admin", Role: userv1.Role_ROLE_ADMIN},
			{Id: "u1", Username: "user", Role: userv1.Role_ROLE_USER},
		},
	}, nil
}

func (m *mockUserClient) Close() error {
	return nil
}

type mockFileClient struct{}

func (m *mockFileClient) ListFiles(ctx context.Context, req *filev1.ListFilesRequest) (*filev1.ListFilesResponse, error) {
	return &filev1.ListFilesResponse{Files: []*filev1.File{}}, nil
}

func (m *mockFileClient) GetFile(ctx context.Context, req *filev1.GetFileRequest) (*filev1.FileResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockFileClient) DeleteFile(ctx context.Context, req *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
	return &filev1.DeleteFileResponse{Success: false}, nil
}

func (m *mockFileClient) UploadFile(ctx context.Context) (filev1.FileService_UploadFileClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockFileClient) DownloadFile(ctx context.Context, req *filev1.DownloadFileRequest) (filev1.FileService_DownloadFileClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockFileClient) Close() error {
	return nil
}

func (h *healthTestContext) iSendAGETRequestTo(endpoint string) error {
	start := time.Now()

	req, err := http.NewRequest("GET", h.server.URL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	if h.customAuthHeader != "" {
		req.Header.Set("Authorization", h.customAuthHeader)
	} else if h.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.authToken)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}

	h.responseTime = time.Since(start)

	if h.response != nil && h.response.Body != nil {
		_ = h.response.Body.Close()
	}
	h.response = resp

	raw, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(raw))
	h.responseRaw = raw

	var body map[string]interface{}
	_ = json.Unmarshal(raw, &body)
	h.responseBody = body

	return nil
}

func (h *healthTestContext) theResponseStatusCodeShouldBe(expectedCode int) error {
	if h.response.StatusCode != expectedCode {
		return fmt.Errorf(
			"expected status code %d, got %d. body=%s",
			expectedCode, h.response.StatusCode, string(h.responseRaw),
		)
	}
	return nil
}

func (h *healthTestContext) theResponseShouldContainWithValue(key, expectedValue string) error {
	value, ok := h.responseBody[key]
	if !ok {
		return fmt.Errorf("response does not contain key %q", key)
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("value for key %q is not a string", key)
	}

	if strValue != expectedValue {
		return fmt.Errorf("expected %q to be %q, got %q", key, expectedValue, strValue)
	}

	return nil
}

func (h *healthTestContext) theResponseTimeShouldBeLessThanMilliseconds(maxMs int) error {
	maxDuration := time.Duration(maxMs) * time.Millisecond
	if h.responseTime > maxDuration {
		return fmt.Errorf("response time %v exceeded maximum %v", h.responseTime, maxDuration)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	h := newHealthTestContext()

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		h = newHealthTestContext()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if h.server != nil {
			h.server.Close()
		}
		if h.response != nil {
			h.response.Body.Close()
		}
		return ctx, nil
	})

	ctx.Step(`^I send a GET request to "([^"]*)"$`, h.iSendAGETRequestTo)
	ctx.Step(`^the response status code should be (\d+)$`, h.theResponseStatusCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, h.theResponseShouldContainWithValue)
	ctx.Step(`^the response time should be less than (\d+) milliseconds$`, h.theResponseTimeShouldBeLessThanMilliseconds)
	ctx.Step(`^I send a POST request to "([^"]*)" with json:$`, h.iSendAPOSTRequestToWithJSON)
	ctx.Step(`^I send a DELETE request to "([^"]*)" with json:$`, h.iSendADELETERequestToWithJSON)
	ctx.Step(`^I am authenticated as "([^"]*)"$`, h.iAmAuthenticatedAs)
	ctx.Step(`^I set headers:$`, h.iSetHeaders)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"health.feature", "auth.feature", "admin.feature", "security.feature", "list_users.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
