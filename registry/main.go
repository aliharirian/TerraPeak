package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aliharirian/TerraPeak/api"
	"github.com/aliharirian/TerraPeak/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "Path to the configuration file")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.Configure(configPath, log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}
	_ = cfg

	logger.Init("TerraPeak", nil, cfg.Log.Level, "15:04:05.0000T2006-01-02")
	logger.Infof("Loaded configuration %s", configPath)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestLogger(&logger.ZerologAdapter{}))

	svc, err := api.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize API service")
	}
	svc.RegisterRoutes(router)

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	log.Info().Str("addr", server.Addr).Msg("Starting Terraform Registry server")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
