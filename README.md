# TerraPeak

[![CI](https://github.com/aliharirian/TerraPeak/actions/workflows/ci.yml/badge.svg)](https://github.com/aliharirian/TerraPeak/actions/workflows/ci.yml)

---

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://hub.docker.com/r/aliharirian/terrapeak)

**A high-performance caching proxy for Terraform Registry that accelerates provider downloads with intelligent storage backends.**

TerraPeak acts as a transparent caching layer between your Terraform workflows and the official Terraform Registry, dramatically reducing download times and bandwidth usage for frequently accessed providers.

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

The easiest way to get started with TerraPeak:

```bash
# Clone the repository
git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak

# Start TerraPeak with MinIO storage backend
docker-compose up -d

# Check if services are running
docker-compose ps
```

This will start:
- **TerraPeak** on port `8081` with MinIO caching
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

> **💡 Pro Tip**: Use Docker Compose for the complete setup with MinIO storage backend and nginx reverse proxy.

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
  # If you want to use Minio Object Storage
  minio:
    enabled: true                  # Enable MinIO object storage
    endpoint: "http://minio:9000"  # MinIO server endpoint
    region: "us-east-1"            # AWS region for MinIO
    access_key: "minioadmin"       # MinIO access key
    secret_key: "minioadmin"       # MinIO secret key
    bucket: "proxy-cache"          # Storage bucket name
    skip_ssl_verify: true          # Skip SSL verification (dev only)

  # If you want to use File Storage disable Minio Storage
  file:
    path: "/data/registry"         # Local filesystem path
```

### 🔐 SSL Requirements

> **⚠️ Important**: The `server.domain` must use HTTPS with a valid SSL certificate. Terraform requires secure connections for provider downloads and will reject HTTP or self-signed certificates.

**Options for SSL:**
- Use a reverse proxy (nginx) with Let's Encrypt certificates
- Configure your own SSL certificates
- Use a cloud load balancer with SSL termination

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

TerraPeak implements the Terraform Registry API specification:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Health check endpoint |
| `/v1/providers/{namespace}/{name}/versions` | GET | List provider versions |
| `/v1/providers/{namespace}/{name}/{version}/download/{os}/{arch}` | GET | Download provider binary |

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
- **Dual Storage**: MinIO object storage or local filesystem support
- **Drop-in Replacement**: Fully compatible with Terraform Registry API

### 🛠️ Developer Experience
- **Easy Setup**: Docker Compose configuration for quick deployment
- **Flexible Configuration**: YAML-based configuration with comprehensive options
- **Health Monitoring**: Built-in health checks and logging
- **SSL Ready**: HTTPS support with reverse proxy configuration

### 🔧 Storage Options
- **MinIO Integration**: Scalable object storage for production environments
- **Local Filesystem**: Simple file-based caching for development
- **Configurable Backends**: Easy switching between storage types

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
                       │  (MinIO/FS)  │
                       └──────────────┘
```

**How it works:**
1. Terraform requests a provider from TerraPeak
2. TerraPeak checks local cache first
3. If not cached, fetches from upstream registry
4. Caches the provider for future requests
5. Returns provider to Terraform

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

# Build and test
cd registry
go build -o terrapeak
go test ./...

# Run with development config
./terrapeak -c ../cfg.yml
```

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
- [x] MinIO Storage Backend
- [x] Local Filesystem Storage Backend
- [x] Docker Compose Setup
- [x] Nginx Reverse Proxy Configuration
- [x] CI/CD Integration

### 🚧 In Progress
- [ ] Implement Go interface for store package

### 📋 Planned
- [ ] Support for HTTP/HTTPS/SOCKS5 Proxy
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
