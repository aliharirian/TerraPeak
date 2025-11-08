package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
)

// StoreInterface defines the interface for the store that the cache handler will use
// This matches the methods available in the existing store package
type StoreInterface interface {
	FileExists(filePath string) bool
	ReadFromStorage(filePath string) ([]byte, error)
	Save(filename string, data []byte) error
}

// Handler handles HTTP requests with transparent caching and proxying
type Handler struct {
	store      StoreInterface
	config     *Config
	httpClient *http.Client
}

// NewCacheHandler creates a new cache handler with the given store and configuration
func NewCacheHandler(store StoreInterface, config *Config) (*Handler, error) {
	if store == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cache config: %w", err)
	}

	return NewCacheHandlerWithClient(store, config, nil)
}

// NewCacheHandlerWithClient creates a new cache handler with injected HTTP client
func NewCacheHandlerWithClient(store StoreInterface, config *Config, httpClient *http.Client) (*Handler, error) {
	if store == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cache config: %w", err)
	}

	return &Handler{
		store:      store,
		config:     config,
		httpClient: httpClient,
	}, nil
}

// Handle is the main HTTP handler that implements caching and proxying logic
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request to extract host and path
	proxyReq, err := ParseRequest(r)
	if err != nil {
		logger.Debugf("Invalid cache request path %s: %v", r.URL.Path, err)
		http.NotFound(w, r)
		return
	}

	// Check if the host is allowed
	if !h.config.IsHostAllowed(proxyReq.Host) {
		logger.Warnf("Host %s is not in allowed hosts list", proxyReq.Host)
		http.Error(w, "Forbidden: Host not allowed", http.StatusForbidden)
		return
	}

	logger.Infof("Processing request for %s%s", proxyReq.Host, proxyReq.Path)

	// Generate cache key for this request
	cacheKey := GenerateCacheKey(proxyReq)

	// Check if content exists in cache
	if h.store.FileExists(cacheKey) {
		logger.Infof("Cache HIT: Serving cached content for %s", cacheKey)
		h.serveCachedContent(w, cacheKey)
		return
	}

	// Cache miss - need to proxy to upstream and cache the result
	logger.Infof("Cache MISS: Proxying request to upstream %s", proxyReq.Host)
	h.proxyAndCache(w, proxyReq, cacheKey)
}

// serveCachedContent serves content from the cache
func (h *Handler) serveCachedContent(w http.ResponseWriter, cacheKey string) {
	data, err := h.store.ReadFromStorage(cacheKey)
	if err != nil {
		logger.Errorf("Failed to read cached content for %s: %v", cacheKey, err)
		http.Error(w, "Internal server error reading cache", http.StatusInternalServerError)
		return
	}

	// Set cache headers
	w.Header().Set("X-Cache-Status", "HIT")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))

	// Write response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		logger.Errorf("Failed to write cached response: %v", err)
	}

	logger.Infof("Successfully served cached content for %s (%d bytes)", cacheKey, len(data))
}

// proxyAndCache proxies the request to upstream server and caches the successful response
func (h *Handler) proxyAndCache(w http.ResponseWriter, proxyReq *ProxyRequest, cacheKey string) {
	// Make upstream request with SSL verification config
	resp, err := MakeUpstreamRequestWithConfig(proxyReq, h.httpClient, h.config.SkipSSLVerify)
	if err != nil {
		logger.Errorf("Upstream request failed for %s: %v", proxyReq.Host, err)
		http.Error(w, "Upstream server error", http.StatusBadGateway)
		return
	}

	// Copy response headers to client (excluding hop-by-hop headers)
	copyResponseHeaders(w.Header(), resp.Headers)
	w.Header().Set("X-Cache-Status", "MISS")

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Write response body to client
	if _, err := w.Write(resp.Body); err != nil {
		logger.Errorf("Failed to write response to client: %v", err)
		return
	}

	// Cache the response if it was successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := h.store.Save(cacheKey, resp.Body); err != nil {
			logger.Warnf("Failed to cache response for %s: %v", cacheKey, err)
			// Don't return error to client as the response was already sent
		} else {
			logger.Infof("Successfully cached response for %s (%d bytes)", cacheKey, len(resp.Body))
		}
	} else {
		logger.Infof("Not caching response for %s (status: %d)", cacheKey, resp.StatusCode)
	}

	logger.Infof("Successfully proxied request to %s (%d bytes, status: %d)",
		proxyReq.Host, len(resp.Body), resp.StatusCode)
}

// copyResponseHeaders copies headers from upstream response to client response
func copyResponseHeaders(dst, src http.Header) {
	// Headers that should not be forwarded
	skipHeaders := map[string]bool{
		"Connection":        true,
		"Keep-Alive":        true,
		"Transfer-Encoding": true,
		"Upgrade":           true,
		"Proxy-Connection":  true,
		"Trailer":           true,
	}

	for key, values := range src {
		if !skipHeaders[key] {
			dst[key] = values
		}
	}
}

// Config holds cache configuration settings
type Config struct {
	// AllowedHosts is a list of upstream hosts that are allowed to be proxied
	AllowedHosts []string `yaml:"allowed_hosts"`

	// SkipSSLVerify disables SSL certificate verification for upstream requests
	// WARNING: This should only be used in development or with trusted hosts
	SkipSSLVerify bool `yaml:"skip_ssl_verify"`
}

// IsHostAllowed checks if the given host is in the allowed hosts list
func (c *Config) IsHostAllowed(host string) bool {
	if c == nil || len(c.AllowedHosts) == 0 {
		return false
	}

	// Normalize host by converting to lowercase and removing any port
	normalizedHost := strings.ToLower(host)
	if colonIndex := strings.Index(normalizedHost, ":"); colonIndex != -1 {
		normalizedHost = normalizedHost[:colonIndex]
	}

	for _, allowedHost := range c.AllowedHosts {
		if strings.ToLower(allowedHost) == normalizedHost {
			return true
		}
	}
	return false
}

// Validate ensures the cache configuration is valid
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("cache config cannot be nil")
	}

	if len(c.AllowedHosts) == 0 {
		return fmt.Errorf("cache config must specify at least one allowed host")
	}

	for i, host := range c.AllowedHosts {
		if strings.TrimSpace(host) == "" {
			return fmt.Errorf("allowed host at index %d cannot be empty", i)
		}
	}

	return nil
}

// ProxyRequest contains all the information needed to proxy a request
type ProxyRequest struct {
	Host        string
	Path        string
	Method      string
	Headers     http.Header
	Body        io.ReadCloser
	QueryString string
}

// ProxyResponse contains the response from the upstream server
type ProxyResponse struct {
	StatusCode    int
	Headers       http.Header
	Body          []byte
	ContentLength int64
}

// ParseRequest extracts the target host and path from an incoming HTTP request
// Expected format: /{host}/{path...}
// Example: /github.com/api/v4/projects -> host=github.com, path=/api/v4/projects
func ParseRequest(r *http.Request) (*ProxyRequest, error) {
	if r == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		return nil, fmt.Errorf("invalid path: must start with host")
	}

	// Split path into host and remaining path
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid path: must contain host")
	}

	host := parts[0]
	if host == "" {
		return nil, fmt.Errorf("invalid path: host cannot be empty")
	}

	// Reconstruct the path for the upstream server
	var upstreamPath string
	if len(parts) == 2 {
		upstreamPath = "/" + parts[1]
	} else {
		upstreamPath = "/"
	}

	return &ProxyRequest{
		Host:        host,
		Path:        upstreamPath,
		Method:      r.Method,
		Headers:     r.Header.Clone(),
		Body:        r.Body,
		QueryString: r.URL.RawQuery,
	}, nil
}

// BuildUpstreamURL constructs the full upstream URL from the proxy request
func (pr *ProxyRequest) BuildUpstreamURL() string {
	upstreamURL := fmt.Sprintf("https://%s%s", pr.Host, pr.Path)
	if pr.QueryString != "" {
		upstreamURL += "?" + pr.QueryString
	}
	return upstreamURL
}

// MakeUpstreamRequest performs the actual HTTP request to the upstream server
func MakeUpstreamRequest(proxyReq *ProxyRequest) (*ProxyResponse, error) {
	return MakeUpstreamRequestWithConfig(proxyReq, nil, false)
}

// MakeUpstreamRequestWithConfig performs the actual HTTP request to the upstream server with custom configuration
func MakeUpstreamRequestWithConfig(proxyReq *ProxyRequest, httpClient *http.Client, skipSSLVerify bool) (*ProxyResponse, error) {
	if proxyReq == nil {
		return nil, fmt.Errorf("proxy request cannot be nil")
	}

	upstreamURL := proxyReq.BuildUpstreamURL()

	// Create HTTP request
	req, err := http.NewRequest(proxyReq.Method, upstreamURL, proxyReq.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream request: %w", err)
	}

	// Copy headers (excluding hop-by-hop headers)
	copyHeaders(req.Header, proxyReq.Headers)

	// Prefer injected client; otherwise create a minimal client honoring env proxy and TLS skip
	client := httpClient
	if client == nil {
		transport := &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSLVerify},
		}
		client = &http.Client{Transport: transport}
	} else {
		// Clone and disable global timeout; control timeout with context per request
		cloned := *client
		cloned.Timeout = 0
		client = &cloned
	}

	// Per-request timeout (allow large artifact downloads)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read upstream response: %w", err)
	}

	return &ProxyResponse{
		StatusCode:    resp.StatusCode,
		Headers:       resp.Header.Clone(),
		Body:          body,
		ContentLength: resp.ContentLength,
	}, nil
}

// copyHeaders copies headers from source to destination, excluding hop-by-hop headers
func copyHeaders(dst, src http.Header) {
	// Headers that should not be forwarded (hop-by-hop headers)
	hopByHopHeaders := map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}

	for key, values := range src {
		if !hopByHopHeaders[key] {
			dst[key] = values
		}
	}
}

// GenerateCacheKey generates a cache key from the proxy request
// This is used as the file path in storage
func GenerateCacheKey(proxyReq *ProxyRequest) string {
	key := fmt.Sprintf("%s%s", proxyReq.Host, proxyReq.Path)
	if proxyReq.QueryString != "" {
		// URL encode the query string to make it filesystem-safe
		encoded := url.QueryEscape(proxyReq.QueryString)
		key += "?" + encoded
	}
	// Remove leading slash if present to make it a valid file path
	return strings.TrimPrefix(key, "/")
}
