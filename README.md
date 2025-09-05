# TerraPeak

A high-performance caching proxy for Terraform Registry that accelerates provider downloads with intelligent storage backends.

## Quick Start

### Installation

```bash
git clone https://github.com/aliharirian/terrapeak.git
cd terrapeak/registry
go build -o terrapeak
```

### Using Docker

```bash
docker build -t aliharirian/terrapeak .
docker run -p 8081:8081 aliharirian/terrapeak
```
Or useing Docker Compose:

```bash
docker-compose up -d
```
****

## Configuration

Create a `cfg.yml` file:

```yaml
server:
  addr: ":8081"
  domain: "https://your-domain.com"  # Must be HTTPS with valid certificate

terraform:
  registry_url: "https://registry.terraform.io"

storage:
  minio:
    enabled: false  # Set to true for MinIO
    endpoint: "localhost:9000"
    access_key: "minio"
    secret_key: "minio123"
    bucket: "terrapeak-cache"
  file:
    path: "./registry"  # Local storage path
```

> ** Note**
>
> The `server.domain` must use HTTPS with a valid SSL certificate. Terraform requires secure connections for provider downloads and will reject HTTP or self-signed certificates.
> Or over nginx with ssl configured.

### Run

```bash
./terrapeak -c cfg.yml
```

## Usage

### Configure Terraform

Update your Terraform configuration to use TerraPeak:

```hcl
terraform {
  required_providers {
    aws = {
      source = "your-domain.com/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

### API Endpoints

- **Health Check**: `GET /healthz`
- **Provider Versions**: `GET /v1/providers/{namespace}/{name}/versions`
- **Provider Download**: `GET /v1/providers/{namespace}/{name}/{version}/download/{os}/{arch}`

### Example

```bash
# Get AWS provider versions
curl "https://your-domain.com/v1/providers/hashicorp/aws/versions"

# Download AWS provider
curl "https://your-domain.com/v1/providers/hashicorp/aws/5.0.0/download/linux/amd64"
```

## Features

- **Intelligent Caching**: Automatic provider caching
- **Dual Storage**: MinIO object storage or local filesystem
- **High Performance**: Sub-second response times for cached content
- **Drop-in Replacement**: Compatible with Terraform Registry API

## Documentation

For detailed architecture, development guides, and advanced configuration, see the [docs](./docs/Document.md) directory.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new features
4. Run `make test`
5. Submit a pull request

## License

Apache2 License - see [LICENSE](LICENSE) file for details.

## Support

- Create an issue on GitHub
- Check logs for debugging
- Verify configuration file


## RouadMap
- [x] Core Proxy Functionality
- [x] Caching Mechanism
- [x] MinIO Storage Backend
- [x] Local Filesystem Storage Backend
- [ ] CI/CD Integration
- [ ] Suppot for http or socks5 proxies
- [ ] Advanced Caching Policies
- [ ] Web Interface for Management
- [ ] Authentication and Authorization
- [ ] Prometheus Metrics and Monitoring
- [ ] Write Helm Chart for easy deployment on Kubernetes
