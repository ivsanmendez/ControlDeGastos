# Reverse Proxy Configuration

## Overview

For production deployments, use a reverse proxy to:
- Handle HTTPS/SSL termination
- Provide domain-based routing
- Add security headers
- Enable rate limiting (optional)
- Serve multiple applications on one server (optional)

## Option 1: Caddy (Recommended - Easiest)

### Why Caddy?
- Automatic HTTPS with Let's Encrypt
- Zero-configuration SSL renewal
- Simple configuration syntax
- Perfect for single-service deployments

### Installation

```bash
# Ubuntu/Debian
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

### Configuration

Create `/etc/caddy/Caddyfile`:

```caddy
# Replace with your domain
expenses.yourdomain.com {
    # Reverse proxy to Go API (serves both API and SPA)
    reverse_proxy localhost:8080

    # Optional: Enable compression
    encode gzip

    # Optional: Security headers
    header {
        # Enable HSTS
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        # Prevent clickjacking
        X-Frame-Options "SAMEORIGIN"
        # Prevent MIME sniffing
        X-Content-Type-Options "nosniff"
        # XSS protection
        X-XSS-Protection "1; mode=block"
    }

    # Optional: Request logging
    log {
        output file /var/log/caddy/access.log
    }
}
```

### Start Caddy

```bash
# Test configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Reload configuration
sudo systemctl reload caddy

# View logs
sudo journalctl -u caddy -f
```

**That's it!** Caddy automatically:
- Obtains SSL certificate from Let's Encrypt
- Renews certificates before expiration
- Redirects HTTP to HTTPS

### Quick One-Liner (No Config File)

For testing:
```bash
caddy reverse-proxy --from expenses.yourdomain.com --to localhost:8080
```

## Option 2: Nginx

### Why Nginx?
- Most popular web server
- Battle-tested and highly performant
- Extensive documentation
- More manual SSL setup required

### Installation

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx
```

### Configuration

Create `/etc/nginx/sites-available/controldegastos`:

```nginx
# HTTP - redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name expenses.yourdomain.com;

    # Redirect all HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

# HTTPS
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name expenses.yourdomain.com;

    # SSL certificates (certbot will add these)
    # ssl_certificate /etc/letsencrypt/live/expenses.yourdomain.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/expenses.yourdomain.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/controldegastos-access.log;
    error_log /var/log/nginx/controldegastos-error.log;

    # Reverse proxy to Go API
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;

        # WebSocket support (if needed in future)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';

        # Proxy headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_cache_bypass $http_upgrade;
    }

    # Optional: Increase client body size for file uploads
    client_max_body_size 10M;
}
```

### Enable and Get SSL Certificate

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/controldegastos /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Obtain SSL certificate
sudo certbot --nginx -d expenses.yourdomain.com

# Reload nginx
sudo systemctl reload nginx
```

### Auto-renewal

Certbot sets up auto-renewal automatically. Test it:

```bash
sudo certbot renew --dry-run
```

## Option 3: Traefik (Container-Native)

### Why Traefik?
- Designed for containers/microservices
- Automatic service discovery
- Automatic SSL via Let's Encrypt
- Great for multi-service deployments

### Configuration

Add Traefik to `docker-compose.prod.yml`:

```yaml
services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.insecure=false"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.email=your@email.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "letsencrypt:/letsencrypt"
    networks:
      - app-network
    restart: unless-stopped

  api:
    # ... existing api service config ...
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`expenses.yourdomain.com`)"
      - "traefik.http.routers.api.entrypoints=websecure"
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
      - "traefik.http.services.api.loadbalancer.server.port=8080"
      # HTTP to HTTPS redirect
      - "traefik.http.routers.api-http.rule=Host(`expenses.yourdomain.com`)"
      - "traefik.http.routers.api-http.entrypoints=web"
      - "traefik.http.routers.api-http.middlewares=redirect-to-https"
      - "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme=https"
    # Remove port mapping (Traefik handles it)
    # ports:
    #   - "8080:8080"

volumes:
  letsencrypt:
    driver: local
```

Start with:
```bash
podman-compose -f docker-compose.prod.yml up -d
```

## DNS Configuration

Before using any reverse proxy with a domain:

1. **Point your domain to this server**:
   - Add an A record: `expenses.yourdomain.com` → `<server-public-ip>`
   - Or CNAME if using subdomain: `expenses` → `your-server.com`

2. **Verify DNS propagation**:
   ```bash
   nslookup expenses.yourdomain.com
   dig expenses.yourdomain.com
   ```

3. **Test HTTP first** (before SSL):
   ```bash
   curl http://expenses.yourdomain.com
   ```

## Firewall Configuration

Allow HTTP and HTTPS traffic:

```bash
# UFW (Ubuntu)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Firewalld (RHEL/Fedora)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

## Recommendation

**For this deployment, I recommend Caddy**:
- ✅ Automatic HTTPS (zero config SSL)
- ✅ Simple one-file configuration
- ✅ Automatic certificate renewal
- ✅ Perfect for single-service deployments
- ✅ Just works™

If you already use Nginx or need more complex routing, use Nginx.
If you plan to run multiple containerized services, use Traefik.

## Testing

Once configured:

```bash
# Test HTTPS
curl https://expenses.yourdomain.com/health

# Check certificate
openssl s_client -connect expenses.yourdomain.com:443 -servername expenses.yourdomain.com < /dev/null
```

## Next Steps

- Set up monitoring (see `03_monitoring.md`)
- Configure automated backups (see `04_backups.md`)
- Review security best practices (see `05_security.md`)
