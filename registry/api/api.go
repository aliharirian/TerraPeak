package api

import (
	"encoding/json"
	"net/http"

	"github.com/aliharirian/TerraPeak/cache"
	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
	"github.com/aliharirian/TerraPeak/metrics"
	"github.com/aliharirian/TerraPeak/proxy"
	"github.com/aliharirian/TerraPeak/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service struct {
	cfg          *config.Config
	store        *store.Store
	proxyHandler *proxy.Handler
	cacheHandler *cache.Handler
}

func New(cfg *config.Config) (*Service, error) {
	// Initialize store with config
	st, err := store.New(cfg)
	if err != nil {
		logger.Errorf("Failed to initialize store: %v", err)
		return nil, err
	}

	// Initialize proxy handler
	proxyHandler, err := proxy.NewHandler(cfg)
	if err != nil {
		logger.Errorf("Failed to initialize proxy handler: %v", err)
		return nil, err
	}

	// Initialize cache handler with injected proxy HTTP client
	cacheHandler, err := cache.NewCacheHandlerWithClient(st, &cache.Config{
		AllowedHosts:  cfg.Cache.AllowedHosts,
		SkipSSLVerify: cfg.Cache.SkipSSLVerify,
	}, proxyHandler.GetClient().GetClient())
	if err != nil {
		logger.Errorf("Failed to initialize cache handler: %v", err)
		return nil, err
	}

	return &Service{
		cfg:          cfg,
		store:        st,
		proxyHandler: proxyHandler,
		cacheHandler: cacheHandler,
	}, nil
}

func (s *Service) RegisterRoutes(router chi.Router) {
	// Route HEAD requests to GET handlers automatically (no body)
	router.Use(middleware.GetHead)
	// Root endpoint
	router.Get("/", Hello)

	// Health & Metrics endpoint
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { metrics.Health(w) })
	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) { metrics.Metrics() })

	// Terraform registry endpoints
	router.Get("/.well-known/terraform.json", s.WellKnown)
	router.Get("/v1/providers/{namespace}/{name}/versions", s.GetVersionList)
	router.Get("/v1/providers/{namespace}/{name}/{version}/download/{os}/{arch}", s.GetProviderDownloadDetails)
	//router.Get("/v1/modules/{name}", s.GetModule)

	// Proxy endpoints
	router.HandleFunc("/proxy/http/*", s.proxyHandler.HandleHTTPProxy)
	router.HandleFunc("/proxy/socks", s.HandleSOCKSProxy)
	router.Get("/proxy/info", s.GetProxyInfo)

	// Mount cache handler for allowed hosts only
	for _, host := range s.cfg.Cache.AllowedHosts {
		router.HandleFunc("/"+host, s.cacheHandler.Handle)
		router.HandleFunc("/"+host+"/*", s.cacheHandler.Handle)
	}
}

func (s *Service) WellKnown(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, err := responseWriter.Write([]byte(`{"modules.v1": "/v1/modules/", "providers.v1": "/v1/providers/"}`))
	if err != nil {
		return
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "Welcome to the Terraform Registry API"}`))
	if err != nil {
		logger.Errorf("Failed to write response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleSOCKSProxy handles SOCKS proxy connections
func (s *Service) HandleSOCKSProxy(w http.ResponseWriter, r *http.Request) {
	// Hijack the connection for SOCKS protocol
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		logger.Errorf("Failed to hijack connection: %v", err)
		return
	}

	// Handle SOCKS proxy
	s.proxyHandler.HandleSOCKSProxy(conn)
}

// GetProxyInfo returns information about the proxy configuration
func (s *Service) GetProxyInfo(w http.ResponseWriter, r *http.Request) {
	// Get proxy info from the handler's client
	info := s.proxyHandler.GetClient().GetProxyInfo()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(info); err != nil {
		logger.Errorf("Failed to encode proxy info: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetCacheStatus returns the current cache configuration and status
// (Removed cache-specific endpoints that were previously added.)
