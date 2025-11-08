package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aliharirian/TerraPeak/api"
	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createIntegrationTestConfig(tempDir string) *config.Config {
	return &config.Config{
		Server: struct {
			Addr         string `yaml:"addr"`
			ReadTimeout  int    `yaml:"read_timeout"`
			WriteTimeout int    `yaml:"write_timeout"`
			IdleTimeout  int    `yaml:"idle_timeout"`
			Domain       string `yaml:"domain"`
		}{
			Addr:         ":0", // Let OS choose port
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
			Domain:       "https://test.example.com",
		},
		Log: struct {
			Level string `yaml:"level"`
		}{
			Level: "info",
		},
		Terraform: struct {
			RegistryUrl string `yaml:"registry_url"`
		}{
			RegistryUrl: "https://registry.terraform.io",
		},
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
		Cache: struct {
			AllowedHosts  []string `yaml:"allowed_hosts"`
			SkipSSLVerify bool     `yaml:"skip_ssl_verify"`
			Rewrites      []struct {
				Prefix string `yaml:"prefix"`
				Host   string `yaml:"host"`
			} `yaml:"rewrites"`
		}{
			AllowedHosts:  []string{"github.com", "registry.terraform.io", "gitlab.com"},
			SkipSSLVerify: true,
		},
		ServeIf: true,
	}
}

func TestFullIntegration(t *testing.T) {
	// Create temporary directory for test storage
	tempDir, err := os.MkdirTemp("", "terrapeak-integration-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize logger
	logger.Init("TerraPeak-Test", nil, "info", "15:04:05.0000T2006-01-02")

	// Create test configuration
	cfg := createIntegrationTestConfig(tempDir)

	// Create API service
	svc, err := api.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create API service: %v", err)
	}

	// Create router and register routes
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestLogger(&logger.ZerologAdapter{}))

	svc.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test all endpoints
	t.Run("Root endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("Failed to call root endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Health endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/healthz")
		if err != nil {
			t.Fatalf("Failed to call health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Metrics endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/metrics")
		if err != nil {
			t.Fatalf("Failed to call metrics endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Well-known endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/.well-known/terraform.json")
		if err != nil {
			t.Fatalf("Failed to call well-known endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}

		if result["providers.v1"] != "/v1/providers/" {
			t.Errorf("Expected providers.v1 to be '/v1/providers/', got %v", result["providers.v1"])
		}
	})
}

func TestCachingBehavior(t *testing.T) {
	// Create temporary directory for test storage
	tempDir, err := os.MkdirTemp("", "terrapeak-cache-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock upstream server
	requestCount := 0
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"versions":[{"version":"5.0.0","protocols":["5.0"]}],"request_count":%d}`, requestCount)
	}))
	defer mockUpstream.Close()

	// Initialize logger
	logger.Init("TerraPeak-Test", nil, "debug", "15:04:05.0000T2006-01-02")

	// Create test configuration
	cfg := createIntegrationTestConfig(tempDir)
	cfg.Terraform.RegistryUrl = mockUpstream.URL

	// Create API service
	svc, err := api.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create API service: %v", err)
	}

	// Create router and register routes
	router := chi.NewRouter()
	svc.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Make first request (should be cache miss)
	resp1, err := http.Get(server.URL + "/v1/providers/hashicorp/aws/versions")
	if err != nil {
		t.Fatalf("Failed to make first request: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 on first request, got %d", resp1.StatusCode)
	}

	cacheStatus1 := resp1.Header.Get("X-Cache-Status")
	if cacheStatus1 != "MISS" {
		t.Errorf("Expected X-Cache-Status: MISS on first request, got %s", cacheStatus1)
	}

	// Wait a brief moment to ensure cache is written
	time.Sleep(100 * time.Millisecond)

	// Make second request (should be cache hit)
	resp2, err := http.Get(server.URL + "/v1/providers/hashicorp/aws/versions")
	if err != nil {
		t.Fatalf("Failed to make second request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 on second request, got %d", resp2.StatusCode)
	}

	cacheStatus2 := resp2.Header.Get("X-Cache-Status")
	if cacheStatus2 != "HIT" {
		t.Errorf("Expected X-Cache-Status: HIT on second request, got %s", cacheStatus2)
	}

	// Verify that upstream was only called once
	if requestCount != 1 {
		t.Errorf("Expected upstream to be called once, but was called %d times", requestCount)
	}
}

func TestErrorHandling(t *testing.T) {
	// Create temporary directory for test storage
	tempDir, err := os.MkdirTemp("", "terrapeak-error-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize logger
	logger.Init("TerraPeak-Test", nil, "error", "15:04:05.0000T2006-01-02")

	// Create test configuration with invalid upstream
	cfg := createIntegrationTestConfig(tempDir)
	cfg.Terraform.RegistryUrl = "http://nonexistent-upstream.example.com"

	// Create API service
	svc, err := api.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create API service: %v", err)
	}

	// Create router and register routes
	router := chi.NewRouter()
	svc.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test provider versions endpoint with unreachable upstream
	resp, err := http.Get(server.URL + "/v1/providers/hashicorp/aws/versions")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status 502 (Bad Gateway) for unreachable upstream, got %d", resp.StatusCode)
	}
}

func TestConcurrentRequests(t *testing.T) {
	// Create temporary directory for test storage
	tempDir, err := os.MkdirTemp("", "terrapeak-concurrent-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock upstream server with delay
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"versions":[{"version":"5.0.0","protocols":["5.0"]}]}`))
	}))
	defer mockUpstream.Close()

	// Initialize logger
	logger.Init("TerraPeak-Test", nil, "info", "15:04:05.0000T2006-01-02")

	// Create test configuration
	cfg := createIntegrationTestConfig(tempDir)
	cfg.Terraform.RegistryUrl = mockUpstream.URL

	// Create API service
	svc, err := api.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create API service: %v", err)
	}

	// Create router and register routes
	router := chi.NewRouter()
	svc.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Make multiple concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get(server.URL + "/v1/providers/hashicorp/aws/versions")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests to complete")
		}
	}
}
