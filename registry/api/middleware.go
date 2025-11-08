package api

import (
	"net/http"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/go-chi/chi/v5/middleware"
)

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Version, X-Request-ID")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// APIVersionMiddleware adds API version header to responses
func APIVersionMiddleware(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-API-Version", version)
			next.ServeHTTP(w, r)
		})
	}
}

// RequestLoggingMiddleware logs API requests with timing
func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		requestID := middleware.GetReqID(r.Context())

		logger.Infof("API request: %s %s (requestId: %s, remoteAddr: %s)", 
			r.Method, r.URL.Path, requestID, r.RemoteAddr)

		next.ServeHTTP(ww, r)

		logger.Infof("Request completed: %s (requestId: %s, duration: %v, status: %d, bytes: %d)",
			r.URL.Path, requestID, time.Since(start), ww.Status(), ww.BytesWritten())
	})
}
