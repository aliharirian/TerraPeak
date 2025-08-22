package api

import (
	"net/http"

	"github.com/aliharirian/TerraPeak/config"
	"github.com/aliharirian/TerraPeak/logger"
	"github.com/aliharirian/TerraPeak/metrics"
	"github.com/aliharirian/TerraPeak/store"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	cfg   *config.Config
	store *store.Store
}

func New(cfg *config.Config) (*Service, error) {
	// Initialize store with config
	st, err := store.New(cfg)
	if err != nil {
		logger.Errorf("Failed to initialize store: %v", err)
		return nil, err
	}

	return &Service{
		cfg:   cfg,
		store: st,
	}, nil
}

func (s *Service) RegisterRoutes(router chi.Router) {
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

	// Proxy and cache endpoints
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) { s.store.HandleRequest(w, r) })

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
