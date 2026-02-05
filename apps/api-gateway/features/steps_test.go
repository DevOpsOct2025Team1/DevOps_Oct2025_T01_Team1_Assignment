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

// And I set headers:
// """
// Authorization: Basic abc123
// X-Trace-Id: 123
// """
func (h *healthTestContext) iSetHeaders(headers string) error {
	// reset any previously set custom header
	h.customAuthHeader = ""

	// Very small parser: lines in the form `Key: Value`
	for _, line := range strings.Split(headers, "\n") {
		line = strings.TrimSpace(line)
		if line == "" { // skip blanks
			continue
		}
		// Only support Authorization explicitly for now
		if strings.HasPrefix(strings.ToLower(line), "authorization:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				h.customAuthHeader = strings.TrimSpace(parts[1])
			}
		}
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

	if h.customAuthHeader != "" {
		req.Header.Set("Authorization", h.customAuthHeader)
	} else if h.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.authToken)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	h.responseTime = time.Since(start)

	// Close any previous response body before overwriting h.response to avoid leaks
	if h.response != nil && h.response.Body != nil {
		_ = h.response.Body.Close()
	}
	h.response = resp

	raw, _ := io.ReadAll(resp.Body)
	// Close the original body now that it has been fully read
	_ = resp.Body.Close()
	// Replace the body with a new reader over the captured bytes so it remains usable
	resp.Body = io.NopCloser(bytes.NewReader(raw))
	h.responseRaw = raw

	var m map[string]interface{}
	_ = json.Unmarshal(raw, &m)
	h.responseBody = m

	return nil
}
