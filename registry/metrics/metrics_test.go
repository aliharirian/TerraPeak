package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()

	Health(w)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty health response")
	}

	// Check if response indicates healthy status
	if !contains(body, "ok") && !contains(body, "healthy") && !contains(body, "OK") {
		t.Errorf("Expected health response to indicate healthy status, got: %s", body)
	}
}

func TestMetrics(t *testing.T) {
	// Call the HTTP handler and ensure it responds with 200 and non-empty body
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Metrics() panicked: %v", r)
		}
	}()

	Metrics(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 from metrics handler, got %d", resp.StatusCode)
	}

	if w.Body.Len() == 0 {
		t.Error("Expected non-empty metrics body")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
