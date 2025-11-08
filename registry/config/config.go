package config

import (
	"errors"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr         string `yaml:"addr"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
		IdleTimeout  int    `yaml:"idle_timeout"`
		Domain       string `yaml:"domain"`
	} `yaml:"server"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Terraform struct {
		RegistryUrl string `yaml:"registry_url"`
	} `yaml:"terraform"`

	Storage struct {
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
	}

	ServeIf bool `yaml:"serve_if"`

	Proxy struct {
		Enabled  bool   `yaml:"enabled"`
		Type     string `yaml:"type"` // "http", "socks5", "socks4"
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"proxy"`

	Cache struct {
		AllowedHosts  []string `yaml:"allowed_hosts"`
		SkipSSLVerify bool     `yaml:"skip_ssl_verify"`
		Rewrites      []struct {
			Prefix string `yaml:"prefix"`
			Host   string `yaml:"host"`
		} `yaml:"rewrites"`
	} `yaml:"cache"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate(logger zerolog.Logger) error {
	if c.Terraform.RegistryUrl == "" {
		logger.Error().Msg("terraform.registry_url is required but not set")
		return errors.New("terraform.registry_url is required")
	}

	if c.Server.Addr == "" {
		logger.Error().Msg("server.addr is required but not set")
		return errors.New("server.addr is required")
	}

	// Validate cache config if any allowed hosts are set
	if len(c.Cache.AllowedHosts) == 0 {
		logger.Warn().Msg("cache.allowed_hosts is empty - cache functionality will be limited")
	}

	logger.Debug().
		Str("registry_url", c.Terraform.RegistryUrl).
		Str("server_addr", c.Server.Addr).
		Int("allowed_hosts", len(c.Cache.AllowedHosts)).
		Msg("Configuration validated successfully")

	return nil
}

var (
	once    sync.Once
	global  *Config
	loadErr error
)

// readYAMLInto reads YAML from path and merges into cfg (must be a pointer).
func readYAMLInto(cfg *Config, path string, logger zerolog.Logger) error {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error().Str("path", path).Err(err).Msg("read config failed")
		return err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		logger.Error().Str("path", path).Err(err).Msg("parse config failed")
		return err
	}
	return nil
}

// Configure loads config by first reading .cfg.default.yml, then merging user config (once).
func Configure(userPath string, logger zerolog.Logger) (*Config, error) {
	once.Do(func() {
		cfg := &Config{}

		const defaultPath = ".cfg.default.yml"
		if _, err := os.Stat(defaultPath); err == nil {
			if err := readYAMLInto(cfg, defaultPath, logger); err != nil {
				loadErr = err
				return
			}
		} else {
			logger.Warn().Str("path", defaultPath).Err(err).Msg("default config not found; continuing without it")
		}

		if userPath != "" {
			if _, err := os.Stat(userPath); err == nil {
				if err := readYAMLInto(cfg, userPath, logger); err != nil {
					loadErr = err
					return
				}
			} else {
				logger.Warn().Str("path", userPath).Err(err).Msg("user config file not found; using defaults only")
			}
		}

		// Validate the configuration
		if err := cfg.Validate(logger); err != nil {
			loadErr = err
			return
		}

		global = cfg
	})

	return global, loadErr
}

// Get returns the loaded config (or nil if Configure hasn't been called).
func Get() *Config { return global }
