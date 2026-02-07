package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/server"
)

func TestHealthEndpoint_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{Environment: "test"}
	srv := server.New(nil, nil, cfg)

	ts := httptest.NewServer(srv.Router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", resp.StatusCode, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if body.Status != "ok" {
		t.Fatalf("unexpected body.status: got %q want %q", body.Status, "ok")
	}
}
