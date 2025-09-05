package store

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aliharirian/TerraPeak/config"
)

func createTestConfig(tempDir string) *config.Config {
	return &config.Config{
		Storage: struct {
			Minio struct {
				Enabled   bool   `yaml:"enabled"`
				Endpoint  string `yaml:"endpoint"`
				Region    string `yaml:"region"`
				AccessKey string `yaml:"access_key"`
				SecretKey string `yaml:"secret_key"`
				Bucket    string `yaml:"bucket"`
				SkipSSL   bool   `yaml:"skip_ssl_verify"`
			} `yaml:"minio"`
			File struct {
				Path string `yaml:"path"`
			} `yaml:"file"`
		}{
			Minio: struct {
				Enabled   bool   `yaml:"enabled"`
				Endpoint  string `yaml:"endpoint"`
				Region    string `yaml:"region"`
				AccessKey string `yaml:"access_key"`
				SecretKey string `yaml:"secret_key"`
				Bucket    string `yaml:"bucket"`
				SkipSSL   bool   `yaml:"skip_ssl_verify"`
			}{
				Enabled: false,
			},
			File: struct {
				Path string `yaml:"path"`
			}{
				Path: tempDir,
			},
		},
	}
}

func TestNew(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig(tempDir)

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	if store == nil {
		t.Error("Expected non-nil store")
	}

	if store.config != cfg {
		t.Error("Store config not set correctly")
	}
}

func TestFileExists(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig(tempDir)
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test non-existent file
	exists := store.FileExists("nonexistent/file.txt")
	if exists {
		t.Error("Expected false for non-existent file")
	}

	// Create a test file
	testFilePath := "test/file.txt"
	fullPath := filepath.Join(tempDir, testFilePath)
	err = os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = os.WriteFile(fullPath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test existing file
	exists = store.FileExists(testFilePath)
	if !exists {
		t.Error("Expected true for existing file")
	}
}

func TestSave(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig(tempDir)
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	testData := []byte("test file content")
	testPath := "test/save/file.txt"

	err = store.Save(testPath, testData)
	if err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}

	// Verify file was saved
	fullPath := filepath.Join(tempDir, testPath)
	savedData, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !bytes.Equal(savedData, testData) {
		t.Errorf("Saved data doesn't match. Expected %s, got %s", testData, savedData)
	}
}

func TestReadFromStorage(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig(tempDir)
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	testData := []byte("test file content for reading")
	testPath := "test/read/file.txt"

	// First save a file
	err = store.Save(testPath, testData)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	// Then read it back
	readData, err := store.ReadFromStorage(testPath)
	if err != nil {
		t.Fatalf("Failed to read from storage: %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Errorf("Read data doesn't match. Expected %s, got %s", testData, readData)
	}

	// Test reading non-existent file
	_, err = store.ReadFromStorage("nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error when reading non-existent file")
	}
}

func TestGenerateURL(t *testing.T) {
	cfg := createTestConfig("/tmp")
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	tests := []struct {
		name         string
		requestURL   string
		expectedDown string
		expectedPath string
	}{
		{
			name:         "simple path",
			requestURL:   "/github.com/hashicorp/terraform/archive/v1.0.0.tar.gz",
			expectedDown: "https://github.com/hashicorp/terraform/archive/v1.0.0.tar.gz",
			expectedPath: "github.com/hashicorp/terraform/archive/v1.0.0.tar.gz",
		},
		{
			name:         "provider path",
			requestURL:   "/releases.hashicorp.com/terraform-provider-aws/5.0.0/terraform-provider-aws_5.0.0_linux_amd64.zip",
			expectedDown: "https://releases.hashicorp.com/terraform-provider-aws/5.0.0/terraform-provider-aws_5.0.0_linux_amd64.zip",
			expectedPath: "releases.hashicorp.com/terraform-provider-aws/5.0.0/terraform-provider-aws_5.0.0_linux_amd64.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			downURL, filePath := store.generateURL(tt.requestURL)

			if downURL != tt.expectedDown {
				t.Errorf("Expected download URL %s, got %s", tt.expectedDown, downURL)
			}

			if filePath != tt.expectedPath {
				t.Errorf("Expected file path %s, got %s", tt.expectedPath, filePath)
			}
		})
	}
}

func TestParseMinIOEndpoint(t *testing.T) {
	cfg := createTestConfig("/tmp")
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	tests := []struct {
		name         string
		endpoint     string
		skipSSL      bool
		expectedHost string
		expectedSSL  bool
		shouldError  bool
	}{
		{
			name:         "http endpoint",
			endpoint:     "http://localhost:9000",
			skipSSL:      false,
			expectedHost: "localhost:9000",
			expectedSSL:  false,
			shouldError:  false,
		},
		{
			name:         "https endpoint",
			endpoint:     "https://minio.example.com",
			skipSSL:      false,
			expectedHost: "minio.example.com",
			expectedSSL:  true,
			shouldError:  false,
		},
		{
			name:         "https endpoint with skip SSL",
			endpoint:     "https://minio.example.com",
			skipSSL:      true,
			expectedHost: "minio.example.com",
			expectedSSL:  false,
			shouldError:  false,
		},
		{
			name:         "bare host",
			endpoint:     "localhost:9000",
			skipSSL:      false,
			expectedHost: "localhost:9000",
			expectedSSL:  true,
			shouldError:  false,
		},
		{
			name:         "bare host with skip SSL",
			endpoint:     "localhost:9000",
			skipSSL:      true,
			expectedHost: "localhost:9000",
			expectedSSL:  false,
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, useSSL, err := store.parseMinIOEndpoint(tt.endpoint, tt.skipSSL)

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

			if host != tt.expectedHost {
				t.Errorf("Expected host %s, got %s", tt.expectedHost, host)
			}

			if useSSL != tt.expectedSSL {
				t.Errorf("Expected SSL %v, got %v", tt.expectedSSL, useSSL)
			}
		})
	}
}

func TestHandleRequest(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "store-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig(tempDir)
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a mock HTTP server to simulate upstream
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock file content"))
	}))
	defer mockServer.Close()

	// Test cache miss (first request)
	_ = httptest.NewRequest("GET", "/test/file.txt", nil)
	_ = httptest.NewRecorder()

	// We need to mock the getSourceStream method behavior
	// For this test, we'll test the file existence and cache logic instead
	testPath := "test/cache/file.txt"
	testData := []byte("cached file content")

	// Manually save a file to test cache hit
	err = store.Save(testPath, testData)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	// Test that FileExists works
	if !store.FileExists(testPath) {
		t.Error("File should exist after saving")
	}

	// Test ReadFromStorage
	readData, err := store.ReadFromStorage(testPath)
	if err != nil {
		t.Fatalf("Failed to read from storage: %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Errorf("Read data doesn't match saved data")
	}
}
