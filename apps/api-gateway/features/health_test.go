package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/server"
	authv1 "github.com/provsalt/DOP_P01_Team1/common/auth/v1"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
)

type healthTestContext struct {
	server       *httptest.Server
	response     *http.Response
	responseBody map[string]interface{}
	responseTime time.Duration
}

func newHealthTestContext() *healthTestContext {
	mockAuthClient := &mockAuthClient{}
	mockUserClient := &mockUserClient{}

	srv := server.New(mockAuthClient, mockUserClient)
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

func (m *mockAuthClient) Login(_ context.Context, _ *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAuthClient) ValidateToken(_ context.Context, _ *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAuthClient) Close() error {
	return nil
}

type mockUserClient struct{}

func (m *mockUserClient) GetUser(_ context.Context, _ *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserClient) DeleteAccount(_ context.Context, _ *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserClient) Close() error {
	return nil
}

func (h *healthTestContext) iSendAGETRequestTo(endpoint string) error {
	start := time.Now()

	resp, err := http.Get(h.server.URL + endpoint)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}

	h.responseTime = time.Since(start)
	h.response = resp

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	h.responseBody = body

	return nil
}

func (h *healthTestContext) theResponseStatusCodeShouldBe(expectedCode int) error {
	if h.response.StatusCode != expectedCode {
		return fmt.Errorf("expected status code %d, got %d", expectedCode, h.response.StatusCode)
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
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"health.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
