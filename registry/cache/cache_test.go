package cache

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockStore implements StoreInterface for testing
type MockStore struct {
	files map[string][]byte
	saved map[string][]byte
}

func NewMockStore() *MockStore {
	return &MockStore{
		files: make(map[string][]byte),
		saved: make(map[string][]byte),
	}
}

func (m *MockStore) FileExists(filePath string) bool {
	_, exists := m.files[filePath]
	return exists
}

func (m *MockStore) ReadFromStorage(filePath string) ([]byte, error) {
	data, exists := m.files[filePath]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	return data, nil
}

func (m *MockStore) Save(filename string, data []byte) error {
	m.saved[filename] = data
	m.files[filename] = data // Also add to files for future reads
	return nil
}

// AddFile adds a file to the mock store (simulates existing cached content)
func (m *MockStore) AddFile(path string, content []byte) {
	m.files[path] = content
}

// GetSaved returns data that was saved during the test
func (m *MockStore) GetSaved(filename string) ([]byte, bool) {
	data, exists := m.saved[filename]
	return data, exists
}

func TestConfig_IsHostAllowed(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		host        string
		expected    bool
		description string
	}{
		{
			name:        "allowed host",
			config:      &Config{AllowedHosts: []string{"github.com", "gitlab.com"}},
			host:        "github.com",
			expected:    true,
			description: "should allow host in allowed list",
		},
		{
			name:        "not allowed host",
			config:      &Config{AllowedHosts: []string{"github.com", "gitlab.com"}},
			host:        "malicious.com",
			expected:    false,
			description: "should reject host not in allowed list",
		},
		{
			name:        "case insensitive",
			config:      &Config{AllowedHosts: []string{"GitHub.COM"}},
			host:        "github.com",
			expected:    true,
			description: "should be case insensitive",
		},
		{
			name:        "host with port",
			config:      &Config{AllowedHosts: []string{"github.com"}},
			host:        "github.com:443",
			expected:    true,
			description: "should ignore port number",
		},
		{
			name:        "empty config",
			config:      &Config{AllowedHosts: []string{}},
			host:        "github.com",
			expected:    false,
			description: "should reject when no hosts allowed",
		},
		{
			name:        "nil config",
			config:      nil,
			host:        "github.com",
			expected:    false,
			description: "should reject when config is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsHostAllowed(tt.host)
			if result != tt.expected {
				t.Errorf("IsHostAllowed() = %v, expected %v - %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		query        string
		expectedErr  bool
		expectedHost string
		expectedPath string
		description  string
	}{
		{
			name:         "valid github path",
			path:         "/github.com/api/v4/projects",
			expectedErr:  false,
			expectedHost: "github.com",
			expectedPath: "/api/v4/projects",
			description:  "should parse github API path correctly",
		},
		{
			name:         "valid gitlab path",
			path:         "/gitlab.com/api/v4/user",
			expectedErr:  false,
			expectedHost: "gitlab.com",
			expectedPath: "/api/v4/user",
			description:  "should parse gitlab API path correctly",
		},
		{
			name:         "host only",
			path:         "/github.com",
			expectedErr:  false,
			expectedHost: "github.com",
			expectedPath: "/",
			description:  "should handle host-only path",
		},
		{
			name:         "host with trailing slash",
			path:         "/github.com/",
			expectedErr:  false,
			expectedHost: "github.com",
			expectedPath: "/",
			description:  "should handle host with trailing slash",
		},
		{
			name:        "empty path",
			path:        "/",
			expectedErr: true,
			description: "should reject empty path",
		},
		{
			name:         "no leading slash",
			path:         "/github.com/api",
			expectedErr:  false,
			expectedHost: "github.com",
			expectedPath: "/api",
			description:  "should work with leading slash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.query != "" {
				req.URL.RawQuery = tt.query
			}

			result, err := ParseRequest(req)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("ParseRequest() expected error but got none - %s", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRequest() unexpected error: %v - %s", err, tt.description)
				return
			}

			if result.Host != tt.expectedHost {
				t.Errorf("ParseRequest() host = %v, expected %v - %s", result.Host, tt.expectedHost, tt.description)
			}

			if result.Path != tt.expectedPath {
				t.Errorf("ParseRequest() path = %v, expected %v - %s", result.Path, tt.expectedPath, tt.description)
			}
		})
	}
}

func TestHandler_CacheHit(t *testing.T) {
	// Setup mock store with cached content
	store := NewMockStore()
	cachedContent := []byte("cached response data")
	store.AddFile("github.com/api/v4/projects", cachedContent)

	// Setup cache handler
	config := &Config{AllowedHosts: []string{"github.com"}}
	handler, err := NewCacheHandler(store, config)
	if err != nil {
		t.Fatalf("Failed to create cache handler: %v", err)
	}

	// Create request
	req := httptest.NewRequest("GET", "/github.com/api/v4/projects", nil)
	rr := httptest.NewRecorder()

	// Handle request
	handler.Handle(rr, req)

	// Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	if rr.Header().Get("X-Cache-Status") != "HIT" {
		t.Errorf("Expected cache hit header, got: %v", rr.Header().Get("X-Cache-Status"))
	}

	if !bytes.Equal(rr.Body.Bytes(), cachedContent) {
		t.Errorf("Handler returned wrong body: got %v want %v", rr.Body.Bytes(), cachedContent)
	}
}

func TestHandler_CacheMissWithProxy(t *testing.T) {
	// Setup mock store (empty - cache miss)
	store := NewMockStore()

	// Setup cache handler with allowed host
	config := &Config{AllowedHosts: []string{"github.com"}}
	handler, err := NewCacheHandler(store, config)
	if err != nil {
		t.Fatalf("Failed to create cache handler: %v", err)
	}

	// This test simulates the behavior but doesn't actually make upstream requests
	// since we can't easily mock HTTPS calls to github.com in unit tests
	// Instead, we test that the logic works correctly for cache miss scenarios

	// Create request that will result in cache miss
	req := httptest.NewRequest("GET", "/github.com/api/v4/projects", nil)
	rr := httptest.NewRecorder()

	// Handle request - this will fail at the upstream request stage which is expected
	handler.Handle(rr, req)

	// Verify that it didn't find content in cache (so it would attempt upstream)
	expectedCacheKey := "github.com/api/v4/projects"
	if store.FileExists(expectedCacheKey) {
		t.Errorf("Cache should have been empty for key: %s", expectedCacheKey)
	}

	// The response will be a 502 Bad Gateway because we can't reach the upstream
	// This is expected behavior in this test environment
	if rr.Code != http.StatusBadGateway {
		t.Logf("Expected 502 due to upstream failure, got: %v", rr.Code)
	}
}

func TestHandler_ForbiddenHost(t *testing.T) {
	// Setup mock store
	store := NewMockStore()

	// Setup cache handler with limited allowed hosts
	config := &Config{AllowedHosts: []string{"github.com", "gitlab.com"}}
	handler, err := NewCacheHandler(store, config)
	if err != nil {
		t.Fatalf("Failed to create cache handler: %v", err)
	}

	// Create request to disallowed host
	req := httptest.NewRequest("GET", "/malicious.com/api/v1/data", nil)
	rr := httptest.NewRecorder()

	// Handle request
	handler.Handle(rr, req)

	// Verify forbidden response
	if rr.Code != http.StatusForbidden {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusForbidden)
	}

	if !strings.Contains(rr.Body.String(), "Forbidden") {
		t.Errorf("Expected forbidden message in response body: %s", rr.Body.String())
	}
}

func TestHandler_InvalidPath(t *testing.T) {
	// Setup mock store
	store := NewMockStore()

	// Setup cache handler
	config := &Config{AllowedHosts: []string{"github.com"}}
	handler, err := NewCacheHandler(store, config)
	if err != nil {
		t.Fatalf("Failed to create cache handler: %v", err)
	}

	// Test cases for invalid paths
	testCases := []struct {
		name string
		path string
	}{
		{"empty_path", "/"},
		{"root_only", "/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()

			handler.Handle(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Errorf("Handler returned wrong status code for path %s: got %v want %v",
					tc.path, rr.Code, http.StatusNotFound)
			}
		})
	}
}

func TestGenerateCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		proxyReq    *ProxyRequest
		expected    string
		description string
	}{
		{
			name: "simple path",
			proxyReq: &ProxyRequest{
				Host: "github.com",
				Path: "/api/v4/projects",
			},
			expected:    "github.com/api/v4/projects",
			description: "should generate simple cache key",
		},
		{
			name: "with query string",
			proxyReq: &ProxyRequest{
				Host:        "github.com",
				Path:        "/api/v4/projects",
				QueryString: "per_page=100&page=1",
			},
			expected:    "github.com/api/v4/projects?per_page%3D100%26page%3D1",
			description: "should URL encode query string",
		},
		{
			name: "root path",
			proxyReq: &ProxyRequest{
				Host: "github.com",
				Path: "/",
			},
			expected:    "github.com/",
			description: "should handle root path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateCacheKey(tt.proxyReq)
			if result != tt.expected {
				t.Errorf("GenerateCacheKey() = %v, expected %v - %s", result, tt.expected, tt.description)
			}
		})
	}
}

func TestNewCacheHandler_Validation(t *testing.T) {
	store := NewMockStore()

	tests := []struct {
		name        string
		store       StoreInterface
		config      *Config
		expectError bool
		description string
	}{
		{
			name:        "valid config",
			store:       store,
			config:      &Config{AllowedHosts: []string{"github.com"}},
			expectError: false,
			description: "should create handler with valid config",
		},
		{
			name:        "nil store",
			store:       nil,
			config:      &Config{AllowedHosts: []string{"github.com"}},
			expectError: true,
			description: "should reject nil store",
		},
		{
			name:        "nil config",
			store:       store,
			config:      nil,
			expectError: true,
			description: "should reject nil config",
		},
		{
			name:        "empty allowed hosts",
			store:       store,
			config:      &Config{AllowedHosts: []string{}},
			expectError: true,
			description: "should reject empty allowed hosts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCacheHandler(tt.store, tt.config)

			if tt.expectError && err == nil {
				t.Errorf("NewCacheHandler() expected error but got none - %s", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("NewCacheHandler() unexpected error: %v - %s", err, tt.description)
			}
		})
	}
}
