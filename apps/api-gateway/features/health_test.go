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
	"github.com/gin-gonic/gin"
)

type healthTestContext struct {
	server       *httptest.Server
	response     *http.Response
	responseBody map[string]interface{}
	responseTime time.Duration
}

func newHealthTestContext() *healthTestContext {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	server := httptest.NewServer(router)

	return &healthTestContext{
		server: server,
	}
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
