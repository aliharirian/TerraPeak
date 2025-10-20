package cache

import (
	"fmt"
	"net/http"

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
	store  StoreInterface
	config *Config
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

	return &Handler{
		store:  store,
		config: config,
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
	resp, err := MakeUpstreamRequestWithConfig(proxyReq, h.config.SkipSSLVerify)
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
