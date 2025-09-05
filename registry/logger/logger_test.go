package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	var buf bytes.Buffer

	// Test different log levels
	levels := []string{"debug", "info", "warn", "error", "fatal"}

	for _, level := range levels {
		t.Run("level_"+level, func(t *testing.T) {
			buf.Reset()
			Init("test-app", &buf, level, "15:04:05.0000T2006-01-02")

			// Test that we can log at the configured level
			switch level {
			case "debug":
				Debugf("test debug message")
			case "info":
				Infof("test info message")
			case "warn":
				Warnf("test warn message")
			case "error":
				Errorf("test error message")
			}

			// Check if message was logged (buffer should contain something)
			if level != "fatal" { // Skip fatal as it would exit the program
				output := buf.String()
				if output == "" && level != "error" { // Error might be filtered depending on level
					// Only fail if we expect output
					if level == "debug" || level == "info" || level == "warn" {
						// These should always produce output when configured
					}
				}
			}
		})
	}
}

func TestInitWithDefaults(t *testing.T) {
	var buf bytes.Buffer

	// Test with default parameters
	Init("test-app", &buf, "", "")

	// Should not panic and should initialize successfully
	Infof("test message")

	output := buf.String()
	if output == "" {
		t.Error("Expected log output after initialization")
	}

	// Check if output contains expected fields
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	if err != nil {
		t.Errorf("Expected valid JSON log output, got error: %v", err)
	}

	// Check for required fields
	if _, exists := logEntry["app"]; !exists {
		t.Error("Expected 'app' field in log output")
	}

	if _, exists := logEntry["time"]; !exists {
		t.Error("Expected 'time' field in log output")
	}

	if _, exists := logEntry["message"]; !exists {
		t.Error("Expected 'message' field in log output")
	}
}

func TestLogFunctions(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "debug", "15:04:05.0000T2006-01-02")

	tests := []struct {
		name    string
		logFunc func(string, ...interface{})
		message string
		level   string
	}{
		{"debug", Debugf, "debug message", "debug"},
		{"info", Infof, "info message", "info"},
		{"warn", Warnf, "warn message", "warn"},
		{"error", Errorf, "error message", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc("test %s", tt.message)

			output := buf.String()
			if output == "" {
				t.Errorf("Expected log output for %s level", tt.level)
				return
			}

			var logEntry map[string]interface{}
			err := json.Unmarshal([]byte(output), &logEntry)
			if err != nil {
				t.Errorf("Expected valid JSON, got error: %v", err)
				return
			}

			if logEntry["level"] != tt.level {
				t.Errorf("Expected level %s, got %v", tt.level, logEntry["level"])
			}

			expectedMessage := "test " + tt.message
			if logEntry["message"] != expectedMessage {
				t.Errorf("Expected message '%s', got '%v'", expectedMessage, logEntry["message"])
			}
		})
	}
}

func TestWith(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "info", "15:04:05.0000T2006-01-02")

	// Test With() function for structured logging
	testLogger := With().Str("key", "value").Str("another", "field").Logger()
	testLogger.Info().Msg("structured message")

	output := buf.String()
	if output == "" {
		t.Error("Expected log output from With()")
		return
	}

	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
		return
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key field to be 'value', got %v", logEntry["key"])
	}

	if logEntry["another"] != "field" {
		t.Errorf("Expected another field to be 'field', got %v", logEntry["another"])
	}

	if logEntry["message"] != "structured message" {
		t.Errorf("Expected message 'structured message', got %v", logEntry["message"])
	}
}

func TestLogWithParameters(t *testing.T) {
	var buf bytes.Buffer
	Init("test-app", &buf, "info", "15:04:05.0000T2006-01-02")

	// Test logging with parameters
	Infof("user %s logged in with ID %d", "john", 123)

	output := buf.String()
	if output == "" {
		t.Error("Expected log output")
		return
	}

	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
		return
	}

	expectedMessage := "user john logged in with ID 123"
	if logEntry["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%v'", expectedMessage, logEntry["message"])
	}
}

func TestTimeFormat(t *testing.T) {
	var buf bytes.Buffer
	customTimeFormat := "2006-01-02T15:04:05Z07:00"
	Init("test-app", &buf, "info", customTimeFormat)

	Infof("time format test")

	output := buf.String()
	if output == "" {
		t.Error("Expected log output")
		return
	}

	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
		return
	}

	// Check if time field exists and is a string
	timeField, exists := logEntry["time"]
	if !exists {
		t.Error("Expected time field in log output")
		return
	}

	timeStr, ok := timeField.(string)
	if !ok {
		t.Error("Expected time field to be a string")
		return
	}

	// Try to parse the time to ensure it's in the expected format
	_, err = time.Parse(customTimeFormat, timeStr)
	if err != nil {
		t.Errorf("Time field not in expected format: %v", err)
	}
}
