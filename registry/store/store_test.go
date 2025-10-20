package store

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/aliharirian/TerraPeak/config"
)

func createTestConfig(tempDir string) *config.Config {
	return &config.Config{
		Storage: struct {
			S3 struct {
				Enabled   bool   `yaml:"enabled"`
				Endpoint  string `yaml:"endpoint"`
				Region    string `yaml:"region"`
				AccessKey string `yaml:"access_key"`
				SecretKey string `yaml:"secret_key"`
				Bucket    string `yaml:"bucket"`
				SkipSSL   bool   `yaml:"skip_ssl_verify"`
			} `yaml:"s3"`
			File struct {
				Path string `yaml:"path"`
			} `yaml:"file"`
		}{
			S3: struct {
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
		t.Fatal("Expected non-nil store")
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

// TestGenerateURL removed - generateURL method was moved to cache package

// TestParseS3Endpoint removed - parseS3Endpoint is now in s3 package

func TestStoreIntegration(t *testing.T) {
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

	// Test complete workflow: save -> exists -> read
	testPath := "test/integration/file.txt"
	testData := []byte("integration test content")

	// Save file
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
