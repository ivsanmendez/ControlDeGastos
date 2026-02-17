# Cloudflare Tunnel Setup

## Why Cloudflare Tunnel?

Standard HTTPS with Let's Encrypt requires inbound connections on ports 80/443 for ACME challenges. When an ISP blocks these ports (common on residential connections), certificate issuance fails.

Cloudflare Tunnel solves this by creating an **outbound** encrypted connection from the server to Cloudflare's edge — no inbound ports, no port forwarding, no ISP restrictions.

## How It Works

```
Browser → https://cdg.meyis.work
             │
       Cloudflare edge (HTTPS terminated here)
             │  (outbound tunnel, initiated by cloudflared)
       cloudflared container inside pod
             │
       localhost:8080 (Go API + SPA)
```

## Setup (One-Time)

### 1. Create Named Tunnel in Cloudflare Dashboard

1. Go to [one.dash.cloudflare.com](https://one.dash.cloudflare.com) → **Networks → Tunnels**
2. **Create a tunnel** → Cloudflared → name it (e.g., `controldegastos`)
3. On the connector screen, copy the **tunnel token**
4. Add a **Public Hostname**:
   - Subdomain: `cdg`
   - Domain: `meyis.work`
   - Service: `HTTP` → `localhost:8080`
5. Save

### 2. Add Token to Environment

Add to `.env.production`:

```bash
TUNNEL_TOKEN=<your_token_from_cloudflare_dashboard>
```

### 3. Deploy

```bash
./deploy-alt-ports.sh deploy
```

The `controldegastos-cloudflared` container connects automatically using `TUNNEL_TOKEN` from the env file.

## Verifying the Tunnel

```bash
# Check tunnel connection
podman logs controldegastos-cloudflared
# Look for: "Connection established" or "Registered tunnel connection"

# Test public URL
curl https://cdg.meyis.work/health
```

## Current Configuration

| Setting | Value |
|---------|-------|
| Public URL | `https://cdg.meyis.work` |
| Tunnel type | Named tunnel (stable, persistent) |
| HTTPS | Cloudflare edge (automatic, no cert management) |
| Backend | `localhost:8080` inside pod |
| Token | `TUNNEL_TOKEN` in `.env.production` |

## Tunnel vs Quick Tunnel

| | Named Tunnel | Quick Tunnel |
|--|--|--|
| URL | Stable (`cdg.meyis.work`) | Random (`xxx.trycloudflare.com`) |
| Auth | Token in env | None |
| DNS | Configured in dashboard | Auto-assigned |
| Persistence | Survives restarts | Changes on restart |

**Always use named tunnel** in production.

## If the Cloudflared System Service is Installed

`sudo cloudflared service install` installs cloudflared as a systemd service on the host. This conflicts with the container approach (two connectors for same tunnel). Disable it:

```bash
sudo systemctl stop cloudflared
sudo systemctl disable cloudflared
```

The container manages the tunnel lifecycle tied to the pod.
