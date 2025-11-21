package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "terrapeak",
		Name:      "http_request_duration_seconds",
		Help:      "Histogram of HTTP request latencies (seconds)",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path"})

	cacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "cache_hits_total",
		Help:      "Total number of cache hits",
	})

	cacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "cache_misses_total",
		Help:      "Total number of cache misses",
	})

	upstreamRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "upstream_requests_total",
		Help:      "Total number of requests made to upstream servers",
	}, []string{"status"})

	storageErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "storage_errors_total",
		Help:      "Total number of storage-related errors",
	})

	proxyRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "terrapeak",
		Name:      "proxy_requests_total",
		Help:      "Total number of proxy requests handled",
	}, []string{"type"})
)

// Metrics serves Prometheus metrics on the provided ResponseWriter.
func Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// MetricsLogAndServe logs the incoming request then serves Prometheus metrics.
// Use this when you want the metrics package to also perform logging similar to
// how `Health` is used for health checks.
func MetricsLogAndServe(w http.ResponseWriter, r *http.Request) {
	// Log remote address and basic request info
	logger.Infof("HTTP GET /metrics - from %s", r.RemoteAddr)
	Metrics(w, r)
}

// Health is kept in a separate file (health.go) and remains unchanged.

// HTTPMiddleware instruments requests and records metrics. Note: it uses the
// request path as the "path" label which may be high-cardinality. Consider
// using route names or reducing label cardinality in production.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK, written: 0}
		next.ServeHTTP(srw, r)

		path := r.URL.Path
		statusStr := strconv.Itoa(srw.status)
		httpRequests.WithLabelValues(r.Method, path, statusStr).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status  int
	written int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusResponseWriter) Write(b []byte) (int, error) {
	n, err := s.ResponseWriter.Write(b)
	s.written += n
	return n, err
}

// Helper functions used by other packages to increment specific metrics.
func IncCacheHit() {
	cacheHits.Inc()
}

func IncCacheMiss() {
	cacheMisses.Inc()
}

func IncUpstreamRequest(status int) {
	upstreamRequests.WithLabelValues(strconv.Itoa(status)).Inc()
}

func IncStorageError() {
	storageErrors.Inc()
}

func IncProxyRequest(proxyType string) {
	proxyRequests.WithLabelValues(proxyType).Inc()
}
