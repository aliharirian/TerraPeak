<div align="center">

# TerraPeak

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://hub.docker.com/r/aliharirian/terrapeak)

</div>

A high-performance caching proxy for Terraform Registry that accelerates provider downloads with intelligent storage backends.

TerraPeak acts as a transparent caching layer between your Terraform workflows and the official Terraform Registry, dramatically reducing download times and bandwidth usage for frequently accessed providers.

### Performance Benefits

- **First download**: Provider fetched from upstream and cached (X-Cache-Status: MISS)
- **Subsequent downloads**: Served from cache with sub-second response (~5ms, X-Cache-Status: HIT)
- **Bandwidth savings**: Reduce external traffic by up to 95%
- **Offline capability**: Cached providers available even when upstream is down

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone and start
git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak
docker-compose up -d

# Verify
curl http://localhost:8081/healthz
```

This starts TerraPeak on port `8081`.

### Using Docker

```bash
docker pull aliharirian/terrapeak-registry:latest

docker run -d \
  --name terrapeak-registry \
  -p 8081:8081 \
  -v $(pwd)/cfg.yml:/app/cfg.yml:ro \
  aliharirian/terrapeak-registry:latest
```

### Using Binary

Download from [GitHub Releases](https://github.com/aliharirian/TerraPeak/releases):

```bash
# Download and extract
wget https://github.com/aliharirian/TerraPeak/releases/download/v1.0.0/terrapeak-registry_linux_amd64.tar.gz
tar -xzf terrapeak-registry_linux_amd64.tar.gz

# Run
./terrapeak-registry -c cfg.yml
```

For detailed installation methods, see [INSTALLATION.md](docs/INSTALLATION.md).

## Configuration

Create a `cfg.yml` file to configure TerraPeak:

```yaml
server:
  addr: ":8081"
  domain: "https://terrapeak.yourdomain.com"  # HTTPS required for production

log:
  level: "info"  # debug, info, warn, error

# Terraform registry configuration
registry:
  # Upstream Terraform registry URL
  url: "https://registry.terraform.io"

# Storage backend (choose one)
storage:
  # Option 1: S3/MinIO (recommended for production)
  s3:
    enabled: true
    endpoint: "http://minio:9000"
    region: "us-east-1"
    access_key: "minioadmin"
    secret_key: "minioadmin"
    bucket: "proxy-cache"
    skip_ssl_verify: false

  # Option 2: Local filesystem
  file:
    path: "/data/registry"

# Cache configuration
cache:
  # Allowed hosts for caching (providers, modules, etc.)
  allowed_hosts:
    - registry.terraform.io
    - releases.hashicorp.com
    - github.com
    - gitlab.com
  skip_ssl_verify: false  # Only for development

# Optional: HTTP proxy for outbound requests
proxy:
  enabled: false
  type: "http"  # http, socks5, socks4
  host: "proxy.company.com"
  port: 8080
  username: ""  # Optional
  password: ""  # Optional
```

> ### Configuration Notes
> - **server.domain**: Must use HTTPS with valid SSL certificate for Terraform compatibility
> - **registry.url**: Upstream Terraform registry (usually `https://registry.terraform.io`)
> - **storage**: Enable either S3/MinIO or filesystem, not both
> - **proxy**: Configure if you need to route requests through a corporate HTTP proxy

## Usage

### Configure Terraform

Update your Terraform configuration to use TerraPeak:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "terrapeak.yourdomain.com/hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "terrapeak.yourdomain.com/hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}
```

Then run Terraform as usual:

```bash
terraform init
terraform plan
terraform apply
```

### Test the API

```bash
# Health check
curl https://terrapeak.yourdomain.com/healthz

# List AWS provider versions
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/versions

# Get download metadata
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/5.0.0/download/linux/amd64
```

For complete API documentation, see [API.md](docs/registry/API.md).

## SSL/TLS Setup

Terraform requires HTTPS with valid SSL certificates. Setup options:

1. **Nginx Reverse Proxy** with Let's Encrypt (recommended)
2. **Cloud Load Balancer** with SSL termination
3. **Kubernetes Ingress** with cert-manager

Example Nginx configuration:

```nginx
server {
    listen 443 ssl http2;
    server_name terrapeak.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/terrapeak.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/terrapeak.yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Documentation

- **[Installation Guide](docs/INSTALLATION.md)** - Detailed installation methods
- **[API Documentation](docs/registry/API.md)** - Complete API reference
- **[Architecture](docs/registry/ARCHITECTURE.md)** - System design and architecture
- **[Contributing](CONTRIBUTING.md)** - How to contribute

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code of conduct
- Development workflow
- Coding standards
- Testing requirements
- Submitting pull requests

## Roadmap

### Planned Features

- [ ] Web UI for cache management and monitoring
- [ ] Support Terraform Module Caching
- [ ] Helm chart for Kubernetes deployment
- [ ] Advanced caching policies (TTL, LRU, size limits)
- [ ] Authentication and authorization system if possible

## Support

Need help?

- **Documentation**: Start with [INSTALLATION.md](docs/INSTALLATION.md) and [API.md](docs/registry/API.md)
- **Bug Reports**: [Create an issue](https://github.com/aliharirian/TerraPeak/issues/new)
- **Questions**: Use [GitHub Discussions](https://github.com/aliharirian/TerraPeak/discussions)

### Troubleshooting

**Connection Issues:**
```bash
# Check service status
docker-compose ps
docker-compose logs -f terrapeak

# Verify health
curl http://localhost:8081/healthz
```

**Configuration Issues:**
```bash
# Validate configuration
./terrapeak-registry -c cfg.yml --validate

# Enable debug logging
./terrapeak-registry -c cfg.yml --log-level=debug
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

⭐ **Star this repository if you find it useful!**

Made with ❤️ by [Ali Haririan](https://github.com/aliharirian)

</div>
