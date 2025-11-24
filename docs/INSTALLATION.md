# Installation Guide

This guide provides detailed installation methods for TerraPeak.

## Table of Contents

- [Docker Compose (Recommended)](#docker-compose-recommended)
- [Docker](#docker)
- [Binary Installation](#binary-installation)
- [Build from Source](#build-from-source)
- [Kubernetes Deployment](#kubernetes-deployment)

## Docker Compose (Recommended)

The easiest way to get started with TerraPeak including MinIO storage backend.

### Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

### Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/aliharirian/TerraPeak.git
   cd TerraPeak
   ```

2. **Configure the application**
   ```bash
   # Copy and edit the configuration file
   vim registry/.cfg.default.yml
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Verify services are running**
   ```bash
   docker-compose ps
   docker-compose logs -f terrapeak
   ```

This will start:
- **TerraPeak Registry** on port `8081`
- **MinIO** object storage on ports `9000` (API) and `9001` (Console)

### Access Points

- TerraPeak Registry: `http://localhost:8081`
- MinIO Console: `http://localhost:9001` (credentials: minioadmin/minioadmin)
- Health Check: `http://localhost:8081/healthz`

## Docker

Run TerraPeak as a standalone container.

### Pull from Docker Hub

```bash
docker pull aliharirian/terrapeak-registry:latest
```

### Run with default configuration

```bash
docker run -d \
  --name terrapeak-registry \
  -p 8081:8081 \
  -v $(pwd)/cfg.yml:/app/cfg.yml:ro \
  -v terrapeak-data:/data/registry \
  aliharirian/terrapeak-registry:latest
```

### Run with environment variables

```bash
docker run -d \
  --name terrapeak-registry \
  -p 8081:8081 \
  -e SERVER_ADDR=":8081" \
  -e LOG_LEVEL="info" \
  -e TERRAFORM_REGISTRY_URL="https://registry.terraform.io" \
  -v terrapeak-data:/data/registry \
  aliharirian/terrapeak-registry:latest
```

### Run with S3/MinIO backend

```bash
docker run -d \
  --name terrapeak-registry \
  -p 8081:8081 \
  -e STORAGE_S3_ENABLED="true" \
  -e STORAGE_S3_ENDPOINT="http://minio:9000" \
  -e STORAGE_S3_ACCESS_KEY="minioadmin" \
  -e STORAGE_S3_SECRET_KEY="minioadmin" \
  -e STORAGE_S3_BUCKET="proxy-cache" \
  -v $(pwd)/cfg.yml:/app/cfg.yml:ro \
  aliharirian/terrapeak-registry:latest
```

## Binary Installation

Download pre-built binaries from GitHub Releases.

### Download and Install

```bash
# Set your desired version
VERSION="1.0.0"
OS="linux"      # linux, darwin
ARCH="amd64"    # amd64, arm64

# Download the binary
wget https://github.com/aliharirian/TerraPeak/releases/download/v${VERSION}/terrapeak-registry_${OS}_${ARCH}.tar.gz

# Extract
tar -xzf terrapeak-registry_${OS}_${ARCH}.tar.gz

# Make executable
chmod +x terrapeak-registry_${OS}_${ARCH}

# Move to system path (optional)
sudo mv terrapeak-registry_${OS}_${ARCH} /usr/local/bin/terrapeak-registry
```

### Run the binary

```bash
# Create configuration file
cat > cfg.yml <<EOF
server:
  addr: ":8081"
  domain: "localhost"

log:
  level: "info"

terraform:
  registry_url: "https://registry.terraform.io"

storage:
  file:
    path: "./data/registry"
EOF

# Run TerraPeak
terrapeak-registry -c cfg.yml
```

### Create systemd service (Linux)

```bash
sudo tee /etc/systemd/system/terrapeak.service <<EOF
[Unit]
Description=TerraPeak Registry
After=network.target

[Service]
Type=simple
User=terrapeak
Group=terrapeak
ExecStart=/usr/local/bin/terrapeak-registry -c /etc/terrapeak/cfg.yml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# Create user and directories
sudo useradd -r -s /bin/false terrapeak
sudo mkdir -p /etc/terrapeak /var/lib/terrapeak
sudo cp cfg.yml /etc/terrapeak/
sudo chown -R terrapeak:terrapeak /etc/terrapeak /var/lib/terrapeak

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable terrapeak
sudo systemctl start terrapeak
sudo systemctl status terrapeak
```

## Build from Source

Build TerraPeak from the source code.

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak/registry

# Download dependencies
go mod download

# Build the binary
go build -o terrapeak .

# Or use Make
make build
```

### Build with custom flags

```bash
VERSION="dev"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse HEAD)

go build \
  -trimpath \
  -ldflags "-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
  -o terrapeak \
  .
```

### Run

```bash
# Copy default configuration
cp .cfg.default.yml ../cfg.yml

# Edit configuration
nano ../cfg.yml

# Run
./terrapeak -c ../cfg.yml
```

## Post-Installation

### Verify Installation

```bash
# Check health
curl http://localhost:8081/healthz

# Test API
curl http://localhost:8081/v1/providers/hashicorp/aws/versions
```

### Configure Terraform

Update your Terraform configuration to use TerraPeak:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "terrapeak.yourdomain.com/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

### Setup SSL/TLS (Production)

For production deployments, configure SSL/TLS using:

- **Reverse Proxy**: Nginx, Traefik, or HAProxy with Let's Encrypt
- **Cloud Load Balancer**: AWS ALB, GCP Load Balancer, Azure Application Gateway
- **Kubernetes Ingress**: cert-manager with Let's Encrypt

Example nginx configuration:

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

## Troubleshooting

### Common Issues

**Connection refused**
```bash
# Check if service is running
docker ps
# or
systemctl status terrapeak

# Check logs
docker logs terrapeak-registry
# or
journalctl -u terrapeak -f
```

**Permission denied on data directory**
```bash
# Fix permissions
sudo chown -R 1000:1000 /data/registry
```

**MinIO connection failed**
```bash
# Verify MinIO is accessible
curl http://localhost:9000/minio/health/live

# Check network connectivity
docker network inspect terrapeak_default
```

## Upgrade

### Docker

```bash
docker-compose pull
docker-compose up -d
```

### Binary

```bash
# Download new version
wget https://github.com/aliharirian/TerraPeak/releases/download/v${NEW_VERSION}/terrapeak-registry_${OS}_${ARCH}.tar.gz

# Stop service
sudo systemctl stop terrapeak

# Replace binary
sudo tar -xzf terrapeak-registry_${OS}_${ARCH}.tar.gz -C /usr/local/bin/

# Start service
sudo systemctl start terrapeak
```

## Uninstall

### Docker Compose

```bash
docker-compose down -v
```

### Binary

```bash
sudo systemctl stop terrapeak
sudo systemctl disable terrapeak
sudo rm /etc/systemd/system/terrapeak.service
sudo rm /usr/local/bin/terrapeak-registry
sudo rm -rf /etc/terrapeak /var/lib/terrapeak
```

For more information, see the [main README](../README.md).

