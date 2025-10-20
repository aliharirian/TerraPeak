package s3

import (
	"testing"

	"github.com/aliharirian/TerraPeak/config"
)

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		skipSSL    bool
		wantHost   string
		wantUseSSL bool
		wantErr    bool
	}{
		{
			name:       "https_url",
			endpoint:   "https://s3.amazonaws.com",
			skipSSL:    false,
			wantHost:   "s3.amazonaws.com",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "http_url",
			endpoint:   "http://localhost:9000",
			skipSSL:    false,
			wantHost:   "localhost:9000",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "https_with_skip_ssl",
			endpoint:   "https://s3.amazonaws.com",
			skipSSL:    true,
			wantHost:   "s3.amazonaws.com",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "http_with_skip_ssl",
			endpoint:   "http://minio.local:9000",
			skipSSL:    true,
			wantHost:   "minio.local:9000",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "plain_hostname",
			endpoint:   "minio.local:9000",
			skipSSL:    false,
			wantHost:   "minio.local:9000",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "plain_hostname_with_skip_ssl",
			endpoint:   "minio.local:9000",
			skipSSL:    true,
			wantHost:   "minio.local:9000",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "localhost",
			endpoint:   "localhost:9000",
			skipSSL:    false,
			wantHost:   "localhost:9000",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "ip_address_without_scheme",
			endpoint:   "192.168.1.100:9000",
			skipSSL:    false,
			wantHost:   "",
			wantUseSSL: false,
			wantErr:    true, // IP with port requires scheme
		},
		{
			name:       "ip_address_with_http",
			endpoint:   "http://192.168.1.100:9000",
			skipSSL:    false,
			wantHost:   "192.168.1.100:9000",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "aws_s3_endpoint",
			endpoint:   "https://s3.us-west-2.amazonaws.com",
			skipSSL:    false,
			wantHost:   "s3.us-west-2.amazonaws.com",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "minio_default",
			endpoint:   "http://127.0.0.1:9000",
			skipSSL:    false,
			wantHost:   "127.0.0.1:9000",
			wantUseSSL: false,
			wantErr:    false,
		},
		{
			name:       "https_with_port",
			endpoint:   "https://s3.example.com:443",
			skipSSL:    false,
			wantHost:   "s3.example.com:443",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "empty_scheme",
			endpoint:   "s3.amazonaws.com",
			skipSSL:    false,
			wantHost:   "s3.amazonaws.com",
			wantUseSSL: true,
			wantErr:    false,
		},
		{
			name:       "hostname_no_port",
			endpoint:   "s3.custom.com",
			skipSSL:    false,
			wantHost:   "s3.custom.com",
			wantUseSSL: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotUseSSL, err := parseEndpoint(tt.endpoint, tt.skipSSL)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotHost != tt.wantHost {
				t.Errorf("parseEndpoint() gotHost = %v, want %v", gotHost, tt.wantHost)
			}

			if gotUseSSL != tt.wantUseSSL {
				t.Errorf("parseEndpoint() gotUseSSL = %v, want %v", gotUseSSL, tt.wantUseSSL)
			}
		})
	}
}

func TestParseEndpoint_EdgeCases(t *testing.T) {
	t.Run("scheme_priority_over_skipSSL", func(t *testing.T) {
		// When scheme is https, skipSSL should force it to false
		host, useSSL, err := parseEndpoint("https://example.com", true)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if useSSL {
			t.Error("Expected useSSL=false when skipSSL=true, got true")
		}
		if host != "example.com" {
			t.Errorf("Expected host=example.com, got %s", host)
		}
	})

	t.Run("http_ignores_skipSSL_false", func(t *testing.T) {
		// When scheme is http, it should be http regardless of skipSSL
		host, useSSL, err := parseEndpoint("http://example.com", false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if useSSL {
			t.Error("Expected useSSL=false for http scheme, got true")
		}
		if host != "example.com" {
			t.Errorf("Expected host=example.com, got %s", host)
		}
	})

	t.Run("no_scheme_defaults_to_ssl_unless_skipped", func(t *testing.T) {
		// No scheme should default to SSL unless skipSSL is true
		host1, useSSL1, err := parseEndpoint("example.com", false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !useSSL1 {
			t.Error("Expected useSSL=true when no scheme and skipSSL=false")
		}

		host2, useSSL2, err := parseEndpoint("example.com", true)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if useSSL2 {
			t.Error("Expected useSSL=false when skipSSL=true")
		}

		if host1 != host2 {
			t.Error("Hosts should be the same")
		}
	})
}

func TestParseEndpoint_URLParsing(t *testing.T) {
	t.Run("url_with_path", func(t *testing.T) {
		// URL with path should still extract host correctly
		host, useSSL, err := parseEndpoint("https://s3.amazonaws.com/bucket", false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if host != "s3.amazonaws.com" {
			t.Errorf("Expected host=s3.amazonaws.com, got %s", host)
		}
		if !useSSL {
			t.Error("Expected useSSL=true for https")
		}
	})

	t.Run("url_with_query", func(t *testing.T) {
		// URL with query params should extract host correctly
		host, useSSL, err := parseEndpoint("https://s3.amazonaws.com?region=us-east-1", false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if host != "s3.amazonaws.com" {
			t.Errorf("Expected host=s3.amazonaws.com, got %s", host)
		}
		if !useSSL {
			t.Error("Expected useSSL=true for https")
		}
	})
}

func TestParseEndpoint_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		skipSSL  bool
		wantHost string
		wantSSL  bool
	}{
		{
			name:     "aws_s3_standard",
			endpoint: "https://s3.amazonaws.com",
			skipSSL:  false,
			wantHost: "s3.amazonaws.com",
			wantSSL:  true,
		},
		{
			name:     "minio_local_dev",
			endpoint: "http://localhost:9000",
			skipSSL:  false,
			wantHost: "localhost:9000",
			wantSSL:  false,
		},
		{
			name:     "minio_docker",
			endpoint: "http://minio:9000",
			skipSSL:  false,
			wantHost: "minio:9000",
			wantSSL:  false,
		},
		{
			name:     "digital_ocean_spaces",
			endpoint: "https://nyc3.digitaloceanspaces.com",
			skipSSL:  false,
			wantHost: "nyc3.digitaloceanspaces.com",
			wantSSL:  true,
		},
		{
			name:     "self_signed_cert",
			endpoint: "https://internal-s3.company.local:9000",
			skipSSL:  true,
			wantHost: "internal-s3.company.local:9000",
			wantSSL:  false,
		},
		{
			name:     "backblaze_b2",
			endpoint: "https://s3.us-west-000.backblazeb2.com",
			skipSSL:  false,
			wantHost: "s3.us-west-000.backblazeb2.com",
			wantSSL:  true,
		},
		{
			name:     "wasabi",
			endpoint: "https://s3.wasabisys.com",
			skipSSL:  false,
			wantHost: "s3.wasabisys.com",
			wantSSL:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, useSSL, err := parseEndpoint(tt.endpoint, tt.skipSSL)
			if err != nil {
				t.Fatalf("parseEndpoint() error = %v", err)
			}
			if host != tt.wantHost {
				t.Errorf("host = %v, want %v", host, tt.wantHost)
			}
			if useSSL != tt.wantSSL {
				t.Errorf("useSSL = %v, want %v", useSSL, tt.wantSSL)
			}
		})
	}
}

// TestNew_ConfigValidation tests the New function with various configurations
// Note: This requires a running S3/MinIO instance for integration testing
// For unit tests, we would need to mock the minio client
func TestNew_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     func() *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing_endpoint",
			cfg: func() *config.Config {
				cfg := &config.Config{}
				cfg.Storage.S3.Endpoint = ""
				cfg.Storage.S3.AccessKey = "access"
				cfg.Storage.S3.SecretKey = "secret"
				cfg.Storage.S3.Bucket = "bucket"
				return cfg
			},
			wantErr: true,
			errMsg:  "S3 configuration is incomplete: endpoint, access key, secret key, and bucket must be set",
		},
		{
			name: "missing_access_key",
			cfg: func() *config.Config {
				cfg := &config.Config{}
				cfg.Storage.S3.Endpoint = "http://localhost:9000"
				cfg.Storage.S3.AccessKey = ""
				cfg.Storage.S3.SecretKey = "secret"
				cfg.Storage.S3.Bucket = "bucket"
				return cfg
			},
			wantErr: true,
			errMsg:  "S3 configuration is incomplete: endpoint, access key, secret key, and bucket must be set",
		},
		{
			name: "missing_secret_key",
			cfg: func() *config.Config {
				cfg := &config.Config{}
				cfg.Storage.S3.Endpoint = "http://localhost:9000"
				cfg.Storage.S3.AccessKey = "access"
				cfg.Storage.S3.SecretKey = ""
				cfg.Storage.S3.Bucket = "bucket"
				return cfg
			},
			wantErr: true,
			errMsg:  "S3 configuration is incomplete: endpoint, access key, secret key, and bucket must be set",
		},
		{
			name: "missing_bucket",
			cfg: func() *config.Config {
				cfg := &config.Config{}
				cfg.Storage.S3.Endpoint = "http://localhost:9000"
				cfg.Storage.S3.AccessKey = "access"
				cfg.Storage.S3.SecretKey = "secret"
				cfg.Storage.S3.Bucket = ""
				return cfg
			},
			wantErr: true,
			errMsg:  "S3 configuration is incomplete: endpoint, access key, secret key, and bucket must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg()
			_, err := New(cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("New() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
