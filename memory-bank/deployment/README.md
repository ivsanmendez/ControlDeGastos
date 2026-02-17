# Deployment

Production deployment guides and operational procedures.

## Documents

| File | Summary |
|------|---------|
| `01_production_deployment.md` | Complete deployment guide — quick start, env setup, management commands |
| `02_cloudflare_tunnel.md` | Cloudflare Tunnel setup — how it works, configuration, troubleshooting |

## Quick Start

```bash
# 1. Configure
nano .env.production   # set DB password + TUNNEL_TOKEN

# 2. Deploy
./deploy-alt-ports.sh deploy

# 3. Verify
curl https://cdg.meyis.work/health
```

## Convention

Use `NN_topic.md` format for deployment-related documentation.
