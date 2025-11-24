# API Documentation

TerraPeak implements the Terraform Registry Protocol and provides additional proxy endpoints.

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Registry API](#registry-api)
- [Proxy API](#proxy-api)
- [Response Format](#response-format)
- [Error Handling](#error-handling)

## Base URL

All API requests should be made to your TerraPeak instance:

```
https://terrapeak.yourdomain.com
```

## Authentication

Currently, TerraPeak operates without authentication. Future versions will support:
- Basic authentication
- Token-based authentication
- OAuth2

## Rate Limiting

Rate limiting is not currently enforced but may be added in future versions.

## Registry API

### Health Check

Check service health and availability.

**Endpoint:** `GET /healthz`

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Example:**
```bash
curl https://terrapeak.yourdomain.com/healthz
```

### Service Discovery

Terraform service discovery endpoint.

**Endpoint:** `GET /.well-known/terraform.json`

**Response:**
```json
{
  "providers.v1": "/v1/providers/",
  "modules.v1": "/v1/modules/"
}
```

### List Provider Versions

Get all available versions for a provider.

**Endpoint:** `GET /v1/providers/{namespace}/{name}/versions`

**Parameters:**
- `namespace` (path): Provider namespace (e.g., `hashicorp`)
- `name` (path): Provider name (e.g., `aws`)

**Response:**
```json
{
  "versions": [
    {
      "version": "5.31.0",
      "protocols": ["5.0"],
      "platforms": [
        {
          "os": "linux",
          "arch": "amd64"
        },
        {
          "os": "darwin",
          "arch": "arm64"
        }
      ]
    }
  ]
}
```

**Example:**
```bash
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/versions
```

**Cache Headers:**
- `X-Cache-Status: HIT` - Served from cache
- `X-Cache-Status: MISS` - Fetched from upstream

### Get Provider Download URL

Get download URL and metadata for a specific provider version and platform.

**Endpoint:** `GET /v1/providers/{namespace}/{name}/{version}/download/{os}/{arch}`

**Parameters:**
- `namespace` (path): Provider namespace
- `name` (path): Provider name
- `version` (path): Provider version (e.g., `5.31.0`)
- `os` (path): Operating system (`linux`, `darwin`, `windows`)
- `arch` (path): Architecture (`amd64`, `arm64`, `386`)

**Response:**
```json
{
  "protocols": ["5.0"],
  "os": "linux",
  "arch": "amd64",
  "filename": "terraform-provider-aws_5.31.0_linux_amd64.zip",
  "download_url": "https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_linux_amd64.zip",
  "shasums_url": "https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_SHA256SUMS",
  "shasums_signature_url": "https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_SHA256SUMS.sig",
  "shasum": "abc123...",
  "signing_keys": {
    "gpg_public_keys": [
      {
        "key_id": "...",
        "ascii_armor": "..."
      }
    ]
  }
}
```

**Example:**
```bash
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/5.31.0/download/linux/amd64
```

### Download Provider Binary

Download the actual provider binary.

**Endpoint:** `GET /releases.hashicorp.com/terraform-provider-{name}/{version}/{filename}`

**Parameters:**
- `name` (path): Provider name
- `version` (path): Provider version
- `filename` (path): Binary filename

**Response:** Binary file (application/zip)

**Example:**
```bash
curl -O https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_linux_amd64.zip
```

**Cache Headers:**
- `X-Cache-Status: HIT` - Served from cache
- `X-Cache-Status: MISS` - Fetched from upstream and cached

### List Module Versions

Get all available versions for a module.

**Endpoint:** `GET /v1/modules/{namespace}/{name}/{provider}/versions`

**Parameters:**
- `namespace` (path): Module namespace
- `name` (path): Module name
- `provider` (path): Provider name

**Response:**
```json
{
  "modules": [
    {
      "source": "namespace/name/provider",
      "versions": [
        {
          "version": "1.0.0",
          "submodules": []
        }
      ]
    }
  ]
}
```

**Example:**
```bash
curl https://terrapeak.yourdomain.com/v1/modules/terraform-aws-modules/vpc/aws/versions
```

### Download Module

Get module download URL.

**Endpoint:** `GET /v1/modules/{namespace}/{name}/{provider}/{version}/download`

**Response:**
```text
X-Terraform-Get: https://terrapeak.yourdomain.com/github.com/namespace/repo/archive/v1.0.0.tar.gz
```

## Proxy API

### Proxy Information

Get current proxy configuration (sanitized).

**Endpoint:** `GET /proxy/info`

**Response:**
```json
{
  "enabled": true,
  "type": "http",
  "host": "proxy.example.com",
  "port": 8080,
  "authenticated": true
}
```

**Example:**
```bash
curl https://terrapeak.yourdomain.com/proxy/info
```

### HTTP Proxy

Forward HTTP/HTTPS requests through TerraPeak's proxy.

**Endpoint:** `POST /proxy/http/*`

**Headers:**
- `X-Target-URL`: Target URL to proxy
- Standard HTTP headers

**Example:**
```bash
curl -X POST https://terrapeak.yourdomain.com/proxy/http/example \
  -H "X-Target-URL: https://api.example.com/data" \
  -H "Content-Type: application/json" \
  -d '{"key": "value"}'
```

### SOCKS Proxy

SOCKS proxy endpoint for client connections.

**Endpoint:** `POST /proxy/socks`

**Note:** This endpoint is used by SOCKS client libraries, not typically called directly via HTTP.

## Response Format

### Success Response

All successful responses return appropriate status codes and data:

- `200 OK`: Request successful
- `201 Created`: Resource created
- `204 No Content`: Successful, no content to return

### Cache Headers

TerraPeak adds custom headers to indicate cache status:

```
X-Cache-Status: HIT | MISS
X-Cache-Key: /path/to/cached/resource
```

### Content Types

- `application/json`: API responses
- `application/zip`: Provider binaries
- `application/octet-stream`: Generic binary files
- `text/plain`: Checksums and signatures

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "provider_not_found",
    "message": "Provider hashicorp/nonexistent not found",
    "details": "The requested provider does not exist in the registry"
  }
}
```

### HTTP Status Codes

- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `502 Bad Gateway`: Upstream registry unavailable
- `503 Service Unavailable`: Service temporarily unavailable

### Common Error Codes

| Code | Description |
|------|-------------|
| `invalid_request` | Request parameters are invalid |
| `provider_not_found` | Provider does not exist |
| `version_not_found` | Provider version not found |
| `platform_not_supported` | Platform (OS/arch) not available |
| `upstream_error` | Error communicating with upstream registry |
| `storage_error` | Error accessing storage backend |
| `proxy_error` | Error in proxy configuration or connection |

### Example Error Responses

**404 Provider Not Found:**
```json
{
  "errors": ["provider not found"]
}
```

**502 Upstream Unavailable:**
```json
{
  "error": {
    "code": "upstream_error",
    "message": "Unable to reach upstream registry",
    "details": "Connection timeout to registry.terraform.io"
  }
}
```

## Examples

### Complete Provider Download Workflow

```bash
# 1. Check service health
curl https://terrapeak.yourdomain.com/healthz

# 2. List available versions
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/versions | jq .

# 3. Get download URL for specific version
curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/5.31.0/download/linux/amd64 | jq .

# 4. Download provider binary
curl -O https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_linux_amd64.zip

# 5. Verify checksum
curl https://terrapeak.yourdomain.com/releases.hashicorp.com/terraform-provider-aws/5.31.0/terraform-provider-aws_5.31.0_SHA256SUMS
```

### Using with Terraform

Configure Terraform to use TerraPeak:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "terrapeak.yourdomain.com/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
```

Run Terraform:
```bash
terraform init
terraform plan
terraform apply
```

### Testing with curl

```bash
# Test provider versions (cached after first request)
time curl https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/versions

# First request: X-Cache-Status: MISS (~200ms)
# Second request: X-Cache-Status: HIT (~5ms)

# View response headers
curl -I https://terrapeak.yourdomain.com/v1/providers/hashicorp/aws/versions
```

## Rate Limiting (Future)

Future versions will implement rate limiting:

**Headers:**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1609459200
```

**Rate Limit Response:**
```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "retry_after": 60
  }
}
```

## Webhooks (Future)

Future versions will support webhooks for:
- Provider version updates
- Cache invalidation events
- Storage backend changes

For more information, see the [main README](../../README.md) and [Architecture documentation](ARCHITECTURE.md).

