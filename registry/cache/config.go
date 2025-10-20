package cache

import (
	"fmt"
	"strings"
)

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
