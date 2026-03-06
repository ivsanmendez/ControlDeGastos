# Production Deployment Guide

## Overview

ControlDeContabilidad is deployed as a containerized application using Podman pods. The Go API serves both the REST API and the React SPA static files. HTTPS is provided by Cloudflare Tunnel — no port forwarding or SSL certificates required.

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│               Podman Pod "controldecontabilidad"                │
│                                                          │
│  ┌──────────────────┐   ┌──────────────────┐           │
│  │  Go API + SPA    │──▶│  PostgreSQL 16   │           │
│  │  localhost:8080  │   │  localhost:5432  │           │
│  └────────┬─────────┘   └──────────────────┘           │
│           │ (exposed as host:8080 for LAN)              │
│  ┌────────▼─────────────────────────────────┐           │
│  │  cloudflared (Cloudflare Tunnel)         │           │
│  │  outbound tunnel → Cloudflare edge       │           │
│  └──────────────────────────────────────────┘           │
└──────────────────────────────────────────────────────────┘
                          │
              Cloudflare edge (HTTPS)
                          │
              https://cdg.meyis.work
```

## Prerequisites

- Podman installed (rootless)
- `TUNNEL_TOKEN` from Cloudflare Zero Trust dashboard
- `.env.production` configured

## Quick Start

### 1. Configure Environment

Create `.env.production` with:

```bash
# PostgreSQL
POSTGRES_USER=controldecontabilidad
POSTGRES_PASSWORD=<strong_password>
POSTGRES_DB=controldecontabilidad

# Database URL (uses localhost — pod networking)
DATABASE_URL=postgres://controldecontabilidad:<same_password>@localhost:5432/controldecontabilidad?sslmode=disable

# API
PORT=8080

# Cloudflare Tunnel
TUNNEL_TOKEN=<token_from_cloudflare_dashboard>
```

Generate a strong password: `openssl rand -hex 24`

**Never commit `.env.production`** — it is gitignored.

### 2. Deploy

```bash
./deploy.sh deploy
```

### 3. Verify

```bash
# Local health check
curl http://localhost:8080/health

# Public URL
curl https://cdg.meyis.work/health

# Check tunnel is connected
podman logs controldecontabilidad-cloudflared
```

## Managing the Application

```bash
./deploy.sh deploy    # Build + start everything
./deploy.sh cleanup   # Stop and remove all containers
./deploy.sh restart   # Restart the pod
./deploy.sh status    # Show container status
./deploy.sh logs      # Follow API logs
```

Direct Podman commands:
```bash
podman pod stop controldecontabilidad     # Stop
podman pod start controldecontabilidad    # Start
podman logs -f controldecontabilidad-api  # API logs
podman logs -f controldecontabilidad-cloudflared  # Tunnel logs
```

## Database Backups

```bash
# Backup
podman exec controldecontabilidad-db \
    pg_dump -U controldecontabilidad controldecontabilidad > backup-$(date +%Y%m%d).sql

# Restore
podman exec -i controldecontabilidad-db \
    psql -U controldecontabilidad controldecontabilidad < backup-20260217.sql
```

## Monitoring

```bash
# Container resource usage
podman stats

# Health check
curl http://localhost:8080/health

# Database ready check
podman exec controldecontabilidad-db pg_isready -U controldecontabilidad
```

## Troubleshooting

### Application won't start

```bash
podman logs controldecontabilidad-api
podman logs controldecontabilidad-db
```

### Tunnel not connecting

```bash
podman logs controldecontabilidad-cloudflared
# Should show: "Connection established" within ~5 seconds
# If not: verify TUNNEL_TOKEN in .env.production
```

### Database connection errors

```bash
grep DATABASE_URL .env.production
# Must use localhost, NOT the service name "db":
# DATABASE_URL=postgres://user:pass@localhost:5432/...
```

### SPA routes not working

The Go API serves `index.html` for all non-API routes. If broken:
1. Verify React build: `podman exec controldecontabilidad-api ls /web/dist/`
2. Verify `STATIC_DIR=/web/dist` is set
3. Rebuild: `./deploy.sh deploy`
