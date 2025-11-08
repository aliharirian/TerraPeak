# TerraPeak

[![CI](https://github.com/aliharirian/TerraPeak/actions/workflows/ci.yml/badge.svg)](https://github.com/aliharirian/TerraPeak/actions/workflows/ci.yml)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://hub.docker.com/r/aliharirian/terrapeak)

**A high-performance caching proxy for Terraform Registry that accelerates provider downloads with intelligent storage backends and comprehensive proxy support.**

TerraPeak acts as a transparent caching layer between your Terraform workflows and the official Terraform Registry, dramatically reducing download times and bandwidth usage for frequently accessed providers. It also provides both outbound proxy client functionality and inbound proxy server capabilities for corporate environments.

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

The easiest way to get started with TerraPeak:

```bash
# Clone the repository
git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak

# Start TerraPeak with S3/MinIO storage backend
docker-compose up -d

# Check if services are running
docker-compose ps
```

This will start:
- **TerraPeak** on port `8081` with S3/MinIO caching
- **MinIO** object storage on ports `9000` (API) and `9001` (Console)
- **Nginx** reverse proxy with SSL termination (if configured on `.nginx/docker-compose.yml` path)

### Option 2: Docker Run

```bash
# Pull the latest image
docker pull aliharirian/terrapeak:latest

# Run with default configuration
docker run -d \
  --name terrapea:latest \
  -p 8081:8081 \
  -v $(pwd)/cfg.yml:/app/cfg.yml:ro \
  aliharirian/terrapeak:latest
```

### Option 3: Build from Source

```bash
# Clone and build
git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak/registry
go build -o terrapeak

# Run with configuration
./terrapeak -c ../cfg.yml
```
Or useing builded bainary file on Github Packages.

> **💡 Pro Tip**: Use Docker Compose for the complete setup with S3/MinIO storage backend and nginx reverse proxy.

## ⚙️ Configuration

TerraPeak uses a YAML configuration file (`cfg.yml`) to customize behavior. Here's a complete example:

```yaml
server:
  addr: ":8081"                     # Server listen address
  domain: "https://tp.example.com"  # Public domain (HTTPS required)

log:
  level: "info"                     # Log level: debug, info, warn, error

terraform:
  registry_url: "https://registry.terraform.io"  # Upstream registry

storage:
  # If you want to use S3/MinIO Object Storage
  s3:
    enabled: true                  # Enable S3/MinIO object storage
    endpoint: "http://minio:9000"  # S3/MinIO server endpoint
    region: "us-east-1"            # AWS region for S3/MinIO
    access_key: "minioadmin"       # S3/MinIO access key
    secret_key: "minioadmin"       # S3/MinIO secret key
    bucket: "proxy-cache"          # Storage bucket name
    skip_ssl_verify: true          # Skip SSL verification (dev only)

  # If you want to use File Storage disable S3 Storage
  file:
    path: "/data/registry"         # Local filesystem path

# Proxy configuration (optional)
proxy:
  enabled: false                   # Enable proxy functionality
  type: "http"                     # Proxy type: http, socks5, socks4
  host: "127.0.0.1"               # Proxy server hostname or IP
  port: 8080                      # Proxy server port
  username: ""                    # Authentication username (optional)
  password: ""                    # Authentication password (optional)
```

### 🔐 SSL Requirements

> **⚠️ Important**: The `server.domain` must use HTTPS with a valid SSL certificate. Terraform requires secure connections for provider downloads and will reject HTTP or self-signed certificates.

**Options for SSL:**
- Use a reverse proxy (nginx) with Let's Encrypt certificates
- Configure your own SSL certificates
- Use a cloud load balancer with SSL termination

### 🌐 Proxy Configuration

TerraPeak supports comprehensive proxy functionality for corporate environments:

#### Outbound Proxy (Client Mode)
Configure TerraPeak to route all HTTP requests through a corporate proxy:

```yaml
proxy:
  enabled: true
  type: "http"           # http, socks5, socks4
  host: "proxy.company.com"
  port: 8080
  username: "user"       # Optional
  password: "pass"       # Optional
```

#### Supported Proxy Types
- **HTTP/HTTPS**: Standard HTTP proxy with CONNECT method support
- **SOCKS5**: Full SOCKS5 protocol with authentication
- **SOCKS4**: SOCKS4 protocol (username only)

#### Corporate Environment Benefits
- **Network Compliance**: Route all external requests through approved proxies
- **Security**: Centralized traffic monitoring and filtering
- **Bandwidth Control**: Optimize corporate bandwidth usage
- **Audit Trail**: Complete request logging through proxy infrastructure

### 🏃‍♂️ Running TerraPeak

```bash
# With configuration file
./terrapeak -c cfg.yml

# Or using Docker
docker run -v $(pwd)/cfg.yml:/app/cfg.yml:ro aliharirian/terrapeak:latest
```

## 📖 Usage

### 🔧 Configure Terraform

Update your Terraform configuration to use TerraPeak as your provider registry:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "tp.example.com/hashicorp/aws"  # Your TerraPeak domain
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "tp.example.com/hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}
```

### 🌐 API Endpoints

TerraPeak implements the Terraform Registry API specification and provides additional proxy endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Health check endpoint |
| `/v1/providers/{namespace}/{name}/versions` | GET | List provider versions |
| `/v1/providers/{namespace}/{name}/{version}/download/{os}/{arch}` | GET | Download provider binary |
| `/proxy/info` | GET | Get proxy configuration information |
| `/proxy/http/*` | POST | HTTP proxy endpoint |
| `/proxy/socks` | POST | SOCKS proxy endpoint |

### 🧪 Testing the API

```bash
# Health check
curl "https://tp.example.com/healthz"

# Get AWS provider versions
curl "https://tp.example.com/v1/providers/hashicorp/aws/versions"

# Download AWS provider for Linux AMD64
curl "https://tp.example.com/v1/providers/hashicorp/aws/5.0.0/download/linux/amd64"

# Get Kubernetes provider versions
curl "https://tp.example.com/v1/providers/hashicorp/kubernetes/versions"

# Get proxy configuration info
curl "https://tp.example.com/proxy/info"
```

### 🚀 Performance Benefits

- **First download**: Provider fetched from upstream registry and cached
- **Subsequent downloads**: Served from cache with sub-second response times
- **Bandwidth savings**: Reduce external registry traffic by up to 90%
- **Offline capability**: Cached providers available even when upstream is down

## ✨ Features

### 🚀 Performance & Reliability
- **Intelligent Caching**: Automatic provider caching with configurable storage backends
- **High Performance**: Sub-second response times for cached content
- **Flexible Storage**: S3/MinIO object storage or local filesystem support
- **Interface-Based Architecture**: Clean separation of storage backends with Go interfaces
- **Drop-in Replacement**: Fully compatible with Terraform Registry API
- **Proxy Support**: Outbound proxy client for corporate environments

### 🛠️ Developer Experience
- **Easy Setup**: Docker Compose configuration for quick deployment
- **Flexible Configuration**: YAML-based configuration with comprehensive options
- **Health Monitoring**: Built-in health checks and logging
- **SSL Ready**: HTTPS support with reverse proxy configuration

### 🔧 Storage Options
- **S3/MinIO Integration**: Scalable object storage for production environments
- **Local Filesystem**: Simple file-based caching for development
- **Interface-Based Design**: Clean Go interfaces for easy backend switching
- **Automatic Selection**: Smart backend selection based on configuration

### 🌐 Proxy Capabilities
- **Outbound Proxy Client**: Route HTTP requests through HTTP, SOCKS4, or SOCKS5 proxies
- **Inbound Proxy Server**: Act as a proxy server for HTTP and SOCKS protocols
- **Authentication Support**: Username/password authentication for all proxy types
- **Corporate Environment Ready**: Perfect for restricted network environments

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────┐    ┌─────────────────┐
│   Terraform     │───▶│   TerraPeak   │───▶│ Terraform       │
│   CLI/CI/CD     │    │   (Proxy)    │    │ Registry        │
└─────────────────┘    └──────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────┐
                       │   Storage    │
                       │ (S3/FS/...)  │
                       └──────────────┘

┌─────────────────┐    ┌──────────────┐    ┌─────────────────┐
│   Corporate     │───▶│   TerraPeak   │───▶│   Corporate     │
│   Proxy         │    │   (Client)   │    │   Network       │
└─────────────────┘    └──────────────┘    └─────────────────┘
```

**How it works:**
1. Terraform requests a provider from TerraPeak
2. TerraPeak checks local cache first
3. If not cached, fetches from upstream registry (optionally through corporate proxy)
4. Caches the provider for future requests
5. Returns provider to Terraform

**Storage Architecture:**
- **Interface-Based Design**: Clean Go interfaces (`Storage`) for all storage backends
- **Automatic Selection**: Smart backend selection based on configuration
- **S3/MinIO Backend**: Scalable object storage for production environments
- **Filesystem Backend**: Simple file-based caching for development
- **Extensible**: Easy to add new storage backends by implementing the `Storage` interface

**Proxy functionality:**
- **Outbound**: TerraPeak can route all HTTP requests through configured corporate proxies
- **Inbound**: TerraPeak can act as a proxy server for other applications
- **Authentication**: Supports username/password authentication for proxy connections

## 📚 Documentation

For detailed architecture, development guides, and advanced configuration, see the [docs](./docs/Document.md) directory.

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feat/amazing-feature`
3. **Add tests** for new features
4. **Run tests**: `make test`
5. **Commit changes**: `git commit -m 'feat: Add amazing feature'`
6. **Push to branch**: `git push origin feat/amazing-feature`
7. **Submit a pull request**

### Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/TerraPeak.git
cd TerraPeak

# Setup development environment
make dev-setup

# Build and test
make build
make test

# Run with development config
make run-dev
```

### 🔧 Development Workflow

TerraPeak includes comprehensive development tools and pre-commit checks:

#### **Pre-Commit Commands**
```bash
# Standard pre-commit checks (recommended)
make pre-commit

# Quick checks for development
make pre-commit-quick

# Full comprehensive checks
make pre-commit-full
```

#### **Available Make Commands**
```bash
# Development
make dev-setup          # Setup development environment
make build              # Build the application
make test               # Run all tests
make test-unit          # Run unit tests only
make test-integration   # Run integration tests
make test-coverage      # Run tests with coverage
make test-api           # Test API endpoints
make fmt                # Format code
make vet                # Run go vet
make lint               # Run linter
make clean              # Clean build artifacts

# Docker
make docker-build       # Build Docker image
make docker-run         # Run in Docker
make docker-compose-up  # Start with docker-compose

# Git workflow
make git-commit         # Run checks and prepare commit
make git-push           # Run full checks and prepare push
```

#### **Code Quality**
- **Automatic Formatting**: Code is automatically formatted with `go fmt`
- **Linting**: Comprehensive linting with golangci-lint
- **Testing**: Unit, integration, and API tests with coverage
- **Interface Design**: Clean Go interfaces for storage backends
- **Pre-commit Hooks**: Automated quality checks before commits

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Need help? Here are your options:

- **🐛 Bug Reports**: Create an issue on GitHub with detailed logs
- **💡 Feature Requests**: Open a discussion or issue
- **📖 Documentation**: Check the [docs](./docs/Document.md) directory
- **🔍 Troubleshooting**:
  - Verify your `cfg.yml` configuration
  - Check container logs: `docker logs terrapeak`
  - Ensure SSL certificates are valid
  - Test connectivity: `curl https://tp.example.com/healthz`

## 🗺️ Roadmap

### ✅ Completed
- [x] Core Proxy Functionality
- [x] Caching Mechanism
- [x] S3/MinIO Storage Backend
- [x] Local Filesystem Storage Backend
- [x] Go Interface Architecture for Storage
- [x] Docker Compose Setup
- [x] Nginx Reverse Proxy Configuration
- [x] CI/CD Integration
- [x] HTTP/HTTPS/SOCKS5 Proxy Support
- [x] Outbound Proxy Client
- [x] Inbound Proxy Server
- [x] Proxy Authentication
- [x] Pre-commit Workflow
- [x] Development Environment Setup
- [x] Comprehensive Testing Suite

### 📋 Planned
- [ ] Web Interface for Management
- [ ] Advanced Caching Policies
- [ ] Authentication and Authorization
- [ ] Prometheus Metrics and Monitoring
- [ ] Helm Chart for Kubernetes deployment
- [ ] Multi-region support
- [ ] Provider version pinning
- [ ] Cache warming strategies

---

<div align="center">

**⭐ Star this repository if you find it useful!**

Made with ❤️ by [Ali Haririan](https://github.com/aliharirian)

</div>
