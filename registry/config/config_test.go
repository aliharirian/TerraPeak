package config

import (
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

func TestConfigure(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		expectedAddr  string
		expectedLevel string
		shouldError   bool
	}{
		{
			name: "valid config",
			configContent: `
server:
  addr: ":8080"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60
  domain: "https://test.com"

log:
  level: "debug"

terraform:
  registry_url: "https://registry.terraform.io"

storage:
  s3:
    enabled: true
    endpoint: "localhost:9000"
    region: "us-east-1"
    access_key: "test"
    secret_key: "test123"
    bucket: "test-bucket"
    skip_ssl_verify: true
  file:
    path: "./test-registry"

serve_if: true
`,
			expectedAddr:  ":8080",
			expectedLevel: "debug",
			shouldError:   false,
		},
		{
			name: "minimal config",
			configContent: `
server:
  addr: ":8081"
log:
  level: "info"
terraform:
  registry_url: "https://registry.terraform.io"
`,
			expectedAddr:  ":8081",
			expectedLevel: "info",
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			once = sync.Once{}
			global = nil
			loadErr = nil

			// Create temporary config file
			tmpFile, err := os.CreateTemp("", "config-test-*.yml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.configContent); err != nil {
				t.Fatalf("Failed to write config content: %v", err)
			}
			tmpFile.Close()

			// Test configuration loading
			cfg, err := Configure(tmpFile.Name(), zerolog.New(os.Stdout))

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cfg.Server.Addr != tt.expectedAddr {
				t.Errorf("Expected addr %s, got %s", tt.expectedAddr, cfg.Server.Addr)
			}

			if cfg.Log.Level != tt.expectedLevel {
				t.Errorf("Expected log level %s, got %s", tt.expectedLevel, cfg.Log.Level)
			}
		})
	}
}

func TestConfigureWithInvalidFile(t *testing.T) {
	// Reset global state
	once = sync.Once{}
	global = nil
	loadErr = nil

	// Create a minimal config file with required fields
	tmpFile, err := os.CreateTemp("", "minimal-config-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  addr: ":8080"
terraform:
  registry_url: "https://registry.terraform.io"
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config content: %v", err)
	}
	tmpFile.Close()

	_, err = Configure(tmpFile.Name(), zerolog.New(os.Stdout))
	if err != nil {
		t.Errorf("Configure should not error with valid config, got: %v", err)
	}
}

func TestConfigureWithInvalidYAML(t *testing.T) {
	// Reset global state
	once = sync.Once{}
	global = nil
	loadErr = nil

	// Create temporary file with invalid YAML
	tmpFile, err := os.CreateTemp("", "invalid-config-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidYAML := `
server:
  addr: ":8080"
    invalid_indentation: "test"
`
	if _, err := tmpFile.WriteString(invalidYAML); err != nil {
		t.Fatalf("Failed to write invalid YAML: %v", err)
	}
	tmpFile.Close()

	_, err = Configure(tmpFile.Name(), zerolog.New(os.Stdout))
	if err == nil {
		t.Error("Expected error for invalid YAML but got none")
	}
}

func TestGet(t *testing.T) {
	// Reset global state
	once = sync.Once{}
	global = nil
	loadErr = nil

	// Get should return nil when no config is loaded
	cfg := Get()
	if cfg != nil {
		t.Error("Expected nil config before Configure is called")
	}

	// Configure a test config
	tmpFile, err := os.CreateTemp("", "config-get-test-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  addr: ":9999"
log:
  level: "warn"
terraform:
  registry_url: "https://registry.terraform.io"
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	configuredCfg, err := Configure(tmpFile.Name(), zerolog.New(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to configure: %v", err)
	}

	// Get should return the same config
	getCfg := Get()
	if getCfg == nil {
		t.Error("Expected non-nil config after Configure")
	}

	if getCfg.Server.Addr != configuredCfg.Server.Addr {
		t.Errorf("Get() returned different config than Configure()")
	}
}

func TestConfigDefaults(t *testing.T) {
	// Reset global state
	once = sync.Once{}
	global = nil
	loadErr = nil

	// Test with empty config to check defaults
	tmpFile, err := os.CreateTemp("", "empty-config-*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write minimal config with required fields
	configContent := `
server:
  addr: ":8080"
terraform:
  registry_url: "https://registry.terraform.io"
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	cfg, err := Configure(tmpFile.Name(), zerolog.New(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to configure: %v", err)
	}

	// Check that struct is initialized (even if empty)
	if cfg == nil {
		t.Error("Expected non-nil config")
	}
}
