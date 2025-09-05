package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aliharirian/TerraPeak/config"
	"github.com/go-chi/chi/v5"
)

func createTestConfig() *config.Config {
	return &config.Config{
		Server: struct {
			Addr         string `yaml:"addr"`
			ReadTimeout  int    `yaml:"read_timeout"`
			WriteTimeout int    `yaml:"write_timeout"`
			IdleTimeout  int    `yaml:"idle_timeout"`
			Domain       string `yaml:"domain"`
		}{
			Domain: "https://test.example.com",
		},
		Terraform: struct {
			RegistryUrl string `yaml:"registry_url"`
		}{
			RegistryUrl: "https://registry.terraform.io",
		},
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
				Path: "./test-registry",
			},
		},
	}
}

func TestNew(t *testing.T) {
	cfg := createTestConfig()

	service, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	if service == nil {
		t.Error("Expected non-nil service")
	}

	if service.cfg != cfg {
		t.Error("Service config not set correctly")
	}

	if service.store == nil {
		t.Error("Expected non-nil store")
	}
}

func TestWellKnown(t *testing.T) {
	cfg := createTestConfig()
	service, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	req := httptest.NewRequest("GET", "/.well-known/terraform.json", nil)
	w := httptest.NewRecorder()

	service.WellKnown(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	expectedBody := `{"modules.v1": "/v1/modules/", "providers.v1": "/v1/providers/"}`
	if body != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, body)
	}
}

func TestHello(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Hello(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Check if response contains expected message
	if !contains(body, "Welcome") || !contains(body, "Terraform Registry") {
		t.Error("Expected response to contain welcome message")
	}
}

func TestRegisterRoutes(t *testing.T) {
	cfg := createTestConfig()
	service, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	router := chi.NewRouter()
	service.RegisterRoutes(router)

	// Test root endpoint
	testEndpoint(t, router, "GET", "/", http.StatusOK)

	// Test well-known endpoint
	testEndpoint(t, router, "GET", "/.well-known/terraform.json", http.StatusOK)

	// Test health endpoint
	testEndpoint(t, router, "GET", "/healthz", http.StatusOK)

	// Test metrics endpoint
	testEndpoint(t, router, "GET", "/metrics", http.StatusOK)
}

func testEndpoint(t *testing.T, router chi.Router, method, path string, expectedStatus int) {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != expectedStatus {
		t.Errorf("Expected status %d for %s %s, got %d", expectedStatus, method, path, w.Code)
	}
}

func TestGetCachedResponse(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "api-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig()
	cfg.Storage.File.Path = tempDir

	service, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Test non-existent cache
	cacheKey := "test/cache/key"
	response := service.getCachedResponse(cacheKey)
	if response != nil {
		t.Error("Expected nil for non-existent cache")
	}

	// Create a cached response
	testData := []byte(`{"test": "data"}`)
	err = service.store.Save(cacheKey, testData)
	if err != nil {
		t.Fatalf("Failed to save test cache: %v", err)
	}

	// Test existing cache
	response = service.getCachedResponse(cacheKey)
	if response == nil {
		t.Error("Expected non-nil for existing cache")
	}

	if string(response) != string(testData) {
		t.Errorf("Expected cached response %s, got %s", testData, response)
	}
}

func TestCacheResponse(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "api-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := createTestConfig()
	cfg.Storage.File.Path = tempDir

	service, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	cacheKey := "test/cache/response"
	testData := []byte(`{"cached": "response"}`)

	// Cache the response
	service.cacheResponse(cacheKey, testData)

	// Verify it was cached
	if !service.store.FileExists(cacheKey) {
		t.Error("Expected file to exist after caching")
	}

	// Read it back
	cached, err := service.store.ReadFromStorage(cacheKey)
	if err != nil {
		t.Fatalf("Failed to read cached response: %v", err)
	}

	if string(cached) != string(testData) {
		t.Errorf("Expected cached data %s, got %s", testData, cached)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			contains(s[1:], substr) ||
			(len(s) > 0 && s[:len(substr)] == substr))
}
