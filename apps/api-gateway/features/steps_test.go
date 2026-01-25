package features

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// ---- Extra BDD Steps ----

// When I send a POST request to "/path" with json: """..."""
func (h *healthTestContext) iSendAPOSTRequestToWithJSON(endpoint, body string) error {
	return h.sendJSON("POST", endpoint, body)
}

// When I send a DELETE request to "/path" with json: """..."""
func (h *healthTestContext) iSendADELETERequestToWithJSON(endpoint, body string) error {
	return h.sendJSON("DELETE", endpoint, body)
}

// Given I am authenticated as "admin" / "user"
func (h *healthTestContext) iAmAuthenticatedAs(role string) error {
	switch role {
	case "admin":
		h.authToken = "admin-token"
	case "user":
		h.authToken = "user-token"
	default:
		h.authToken = ""
	}
	return nil
}

// Shared request helper
func (h *healthTestContext) sendJSON(method, endpoint, body string) error {
	start := time.Now()

	var reader io.Reader
	if strings.TrimSpace(body) != "" {
		reader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, h.server.URL+endpoint, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if h.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.authToken)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	h.responseTime = time.Since(start)
	h.response = resp

	raw, _ := io.ReadAll(resp.Body)
	h.responseRaw = raw

	var m map[string]interface{}
	_ = json.Unmarshal(raw, &m)
	h.responseBody = m

	return nil
}
