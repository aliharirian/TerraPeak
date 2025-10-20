package cache

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
	return MakeUpstreamRequestWithConfig(proxyReq, false)
}

// MakeUpstreamRequestWithConfig performs the actual HTTP request to the upstream server with custom configuration
func MakeUpstreamRequestWithConfig(proxyReq *ProxyRequest, skipSSLVerify bool) (*ProxyResponse, error) {
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

	// Create HTTP client with optional SSL verification skip
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipSSLVerify,
			},
		},
	}

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
