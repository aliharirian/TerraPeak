# TerraPeak Nginx SSL Setup

This directory contains the Nginx configuration for TerraPeak with SSL support.

## 📁 Files

- `default.conf` - Nginx configuration with SSL proxy
- `docker-compose.nginx.yml` - Docker Compose with Nginx + TerraPeak + MinIO
- `setup-ssl.sh` - SSL certificate setup helper script

## 🚀 Quick Start

### 1. Get SSL Certificates

```bash
# Option 1: Using Certbot (recommended)
sudo certbot certonly --standalone -d tr.tesaco.sbs

# Option 2: Using Docker
docker run --rm -v /etc/letsencrypt:/etc/letsencrypt certbot/certbot certonly --standalone -d tr.tesaco.sbs
```

### 2. Start Services

```bash
# Start all services (Nginx + TerraPeak + MinIO)
docker compose -f docker-compose.nginx.yml up -d

# Check status
docker compose -f docker-compose.nginx.yml ps

# View logs
docker compose -f docker-compose.nginx.yml logs nginx
docker compose -f docker-compose.nginx.yml logs terrapeak
```

### 3. Test SSL

```bash
# Test HTTPS endpoint
curl -I https://tr.tesaco.sbs/healthz

# Test HTTP redirect
curl -I http://tr.tesaco.sbs/healthz
```

## 🔧 Configuration

### Nginx Configuration (`default.conf`)

- **HTTP → HTTPS redirect**: All HTTP traffic redirects to HTTPS
- **SSL termination**: Handles SSL/TLS encryption
- **Reverse proxy**: Proxies requests to TerraPeak container
- **Security headers**: Adds security headers to responses

### Docker Compose (`docker-compose.nginx.yml`)

- **Nginx**: SSL proxy on ports 80/443
- **TerraPeak**: Backend service (internal port 8081)
- **MinIO**: Object storage (ports 9000/9001)
- **Networks**: All services communicate via `terrapeak-network`

## 📋 SSL Certificate Requirements

Your SSL certificates must be located at:
```
/etc/letsencrypt/live/tesaco.sbs/
├── fullchain.pem    # Certificate chain
└── privkey.pem      # Private key
```

## 🔍 Troubleshooting

### Check SSL Certificates
```bash
# Verify certificate exists
ls -la /etc/letsencrypt/live/tesaco.sbs/

# Check certificate validity
openssl x509 -in /etc/letsencrypt/live/tesaco.sbs/fullchain.pem -text -noout
```

### Check Nginx Configuration
```bash
# Test Nginx config
docker compose -f docker-compose.nginx.yml exec nginx nginx -t

# Reload Nginx config
docker compose -f docker-compose.nginx.yml exec nginx nginx -s reload
```

### Check Service Health
```bash
# Check all services
docker compose -f docker-compose.nginx.yml ps

# Check TerraPeak health
curl http://localhost:8081/healthz

# Check Nginx logs
docker compose -f docker-compose.nginx.yml logs nginx
```

## 🌐 Access Points

- **TerraPeak (HTTPS)**: https://tr.tesaco.sbs
- **TerraPeak (HTTP)**: http://tr.tesaco.sbs (redirects to HTTPS)
- **MinIO Console**: http://localhost:9001 (admin: minioadmin/minioadmin)
- **MinIO API**: http://localhost:9000

## 🔄 SSL Certificate Renewal

Let's Encrypt certificates expire every 90 days. Set up automatic renewal:

```bash
# Add to crontab for automatic renewal
sudo crontab -e

# Add this line (runs twice daily)
0 12 * * * /usr/bin/certbot renew --quiet && docker compose -f /path/to/docker-compose.nginx.yml restart nginx
```

## 📝 Notes

- TerraPeak is only accessible internally (no external port)
- All external traffic goes through Nginx
- SSL certificates are mounted read-only from host
- Nginx logs are stored in Docker volume `nginx-logs`

