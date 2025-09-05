package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestZerologAdapter_NewLogEntry(t *testing.T) {
	adapter := &ZerologAdapter{}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "192.168.1.1:12345"

	logEntry := adapter.NewLogEntry(req)

	if logEntry == nil {
		t.Error("Expected non-nil log entry")
		return
	}

	chiEntry, ok := logEntry.(*chiLogEntry)
	if !ok {
		t.Error("Expected chiLogEntry type")
		return
	}

	if chiEntry.req != req {
		t.Error("Expected request to be stored in log entry")
	}
}

func TestChiLogEntry_Write(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "info", "2006-01-02T15:04:05Z07:00")

	req := httptest.NewRequest("GET", "/v1/providers/hashicorp/aws/versions", nil)
	req.Header.Set("User-Agent", "terraform/1.6.0")
	req.RemoteAddr = "192.168.1.100:12345"

	adapter := &ZerologAdapter{}
	logEntry := adapter.NewLogEntry(req)

	chiEntry := logEntry.(*chiLogEntry)

	// Test successful request
	elapsed := 150 * time.Millisecond
	chiEntry.Write(200, 1024, http.Header{}, elapsed, nil)

	output := buf.String()
	if output == "" {
		t.Error("Expected log output")
		return
	}

	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
		return
	}

	// Check required fields
	if logData["method"] != "GET" {
		t.Errorf("Expected method GET, got %v", logData["method"])
	}

	if logData["path"] != "/v1/providers/hashicorp/aws/versions" {
		t.Errorf("Expected path /v1/providers/hashicorp/aws/versions, got %v", logData["path"])
	}

	if logData["status"] != float64(200) {
		t.Errorf("Expected status 200, got %v", logData["status"])
	}

	if logData["bytes"] != float64(1024) {
		t.Errorf("Expected bytes 1024, got %v", logData["bytes"])
	}

	if logData["remote_addr"] != "192.168.1.100:12345" {
		t.Errorf("Expected remote_addr 192.168.1.100:12345, got %v", logData["remote_addr"])
	}

	if logData["user_agent"] != "terraform/1.6.0" {
		t.Errorf("Expected user_agent terraform/1.6.0, got %v", logData["user_agent"])
	}

	// Check message format
	message, ok := logData["message"].(string)
	if !ok {
		t.Error("Expected message to be a string")
		return
	}

	expectedElements := []string{
		"HTTP GET",
		"/v1/providers/hashicorp/aws/versions",
		"200",
		"OK",
		"1024 bytes",
	}

	for _, element := range expectedElements {
		if !strings.Contains(message, element) {
			t.Errorf("Expected message to contain '%s', got: %s", element, message)
		}
	}
}

func TestChiLogEntry_WriteWithDifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		status     int
		statusText string
	}{
		{200, "OK"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
		{502, "Bad Gateway"},
	}

	for _, tc := range testCases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			var buf bytes.Buffer
			Init("test-app", &buf, "info", "2006-01-02T15:04:05Z07:00")

			req := httptest.NewRequest("GET", "/test", nil)
			adapter := &ZerologAdapter{}
			logEntry := adapter.NewLogEntry(req)

			chiEntry := logEntry.(*chiLogEntry)
			chiEntry.Write(tc.status, 512, http.Header{}, 100*time.Millisecond, nil)

			output := buf.String()
			if output == "" {
				t.Error("Expected log output")
				return
			}

			var logData map[string]interface{}
			err := json.Unmarshal([]byte(output), &logData)
			if err != nil {
				t.Errorf("Expected valid JSON, got error: %v", err)
				return
			}

			if logData["status"] != float64(tc.status) {
				t.Errorf("Expected status %d, got %v", tc.status, logData["status"])
			}

			message := logData["message"].(string)
			if !strings.Contains(message, tc.statusText) {
				t.Errorf("Expected message to contain '%s', got: %s", tc.statusText, message)
			}
		})
	}
}

func TestChiLogEntry_Panic(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "error", "2006-01-02T15:04:05Z07:00")

	req := httptest.NewRequest("GET", "/test", nil)
	adapter := &ZerologAdapter{}
	logEntry := adapter.NewLogEntry(req)

	chiEntry := logEntry.(*chiLogEntry)

	panicValue := "test panic"
	stackTrace := []byte("goroutine 1 [running]:\ntest stack trace")

	chiEntry.Panic(panicValue, stackTrace)

	output := buf.String()
	if output == "" {
		t.Error("Expected log output for panic")
		return
	}

	var logData map[string]interface{}
	err := json.Unmarshal([]byte(output), &logData)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
		return
	}

	if logData["level"] != "error" {
		t.Errorf("Expected error level for panic, got %v", logData["level"])
	}

	if logData["panic"] != panicValue {
		t.Errorf("Expected panic value '%s', got %v", panicValue, logData["panic"])
	}

	if logData["message"] != "panic recovered" {
		t.Errorf("Expected message 'panic recovered', got %v", logData["message"])
	}
}

func TestForRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")

	logger := ForRequest(req)

	// Test that ForRequest returns a non-zero logger (basic validation)
	// We just check that the function doesn't panic and returns something
	if logger.GetLevel() > 10 { // zerolog levels are typically -1 to 7
		t.Error("Expected valid logger level from ForRequest")
	}
}

func TestIntegrationWithChiMiddleware(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "info", "2006-01-02T15:04:05Z07:00")

	// This test verifies that the adapter works with chi middleware
	adapter := &ZerologAdapter{}

	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("User-Agent", "test-client/1.0")
	req.RemoteAddr = "10.0.0.1:54321"

	logEntry := adapter.NewLogEntry(req)

	// Simulate middleware logging
	logEntry.Write(201, 256, http.Header{"Content-Type": []string{"application/json"}}, 75*time.Millisecond, nil)

	output := buf.String()
	if output == "" {
		t.Error("Expected log output from middleware integration")
		return
	}

	// Verify the log contains expected request information
	if !strings.Contains(output, "POST") {
		t.Error("Expected log to contain POST method")
	}

	if !strings.Contains(output, "/api/test") {
		t.Error("Expected log to contain request path")
	}

	if !strings.Contains(output, "201") {
		t.Error("Expected log to contain status code")
	}

	if !strings.Contains(output, "256 bytes") {
		t.Error("Expected log to contain response size")
	}
}
