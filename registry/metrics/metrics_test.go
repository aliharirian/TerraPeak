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
	// Since Metrics() doesn't take parameters, we test that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Metrics() panicked: %v", r)
		}
	}()

	Metrics()

	// Test passed if we reach here without panic
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
