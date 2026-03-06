#!/bin/bash
# ControlDeContabilidad - Cloudflare Named Tunnel Deployment
# Public URL: https://cdg.meyis.work
# Requires TUNNEL_TOKEN in .env.production

set -e

PUBLIC_URL="https://cdg.meyis.work"

echo "🚀 ControlDeContabilidad Deployment (Cloudflare Tunnel)"
echo "   Public: $PUBLIC_URL"
echo

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo "❌ Error: .env.production not found!"
    exit 1
fi

# Check TUNNEL_TOKEN is set in the env file
if ! grep -q "^TUNNEL_TOKEN=" .env.production; then
    echo "❌ Error: TUNNEL_TOKEN not found in .env.production"
    echo "   Add: TUNNEL_TOKEN=<your_cloudflare_tunnel_token>"
    exit 1
fi

# Function to stop and clean up
cleanup() {
    echo "🧹 Cleaning up existing deployment..."
    podman stop controldecontabilidad-api controldecontabilidad-db controldecontabilidad-cloudflared 2>/dev/null || true
    podman rm controldecontabilidad-api controldecontabilidad-db controldecontabilidad-cloudflared 2>/dev/null || true
    podman pod stop controldecontabilidad 2>/dev/null || true
    podman pod rm controldecontabilidad 2>/dev/null || true
}

# Function to deploy
deploy() {
    echo "📦 Building application image..."
    podman build -t controldecontabilidad:latest -f Dockerfile .

    echo "🔧 Creating pod..."
    podman pod create --name controldecontabilidad \
        -p 8080:8080

    echo "🐘 Starting PostgreSQL..."
    podman run -d --pod controldecontabilidad --name controldecontabilidad-db \
        --env-file .env.production \
        -v controldecontabilidad_pgdata:/var/lib/postgresql/data \
        docker.io/library/postgres:16-alpine

    echo "⏳ Waiting for database to be ready..."
    sleep 10

    echo "🚀 Starting API..."
    podman run -d --pod controldecontabilidad --name controldecontabilidad-api \
        --env-file .env.production \
        -e STATIC_DIR=/web/dist \
        controldecontabilidad:latest

    echo "🌐 Starting Cloudflare Tunnel..."
    podman run -d --pod controldecontabilidad --name controldecontabilidad-cloudflared \
        --env-file .env.production \
        docker.io/cloudflare/cloudflared:latest \
        tunnel run

    echo
    echo "✅ Deployment complete!"
    echo
    echo "📊 Status:"
    podman ps --pod --filter pod=controldecontabilidad
    echo
    echo "🌐 Application: $PUBLIC_URL"
    echo "   Local:       http://localhost:8080"
    echo
    echo "📝 Useful commands:"
    echo "  View API logs:    podman logs -f controldecontabilidad-api"
    echo "  View tunnel logs: podman logs -f controldecontabilidad-cloudflared"
    echo "  Stop:             podman pod stop controldecontabilidad"
    echo "  Start:            podman pod start controldecontabilidad"
    echo "  Remove:           ./deploy.sh cleanup"
}

# Parse command
case "${1:-deploy}" in
    deploy)
        cleanup
        deploy
        ;;
    cleanup)
        cleanup
        echo "✅ Cleanup complete!"
        ;;
    restart)
        echo "🔄 Restarting..."
        podman pod restart controldecontabilidad
        echo "✅ Restart complete!"
        echo "   Application: $PUBLIC_URL"
        ;;
    status)
        echo "📊 Status:"
        podman ps --pod --filter pod=controldecontabilidad
        echo
        echo "🌐 Application: $PUBLIC_URL"
        ;;
    logs)
        podman logs -f controldecontabilidad-api
        ;;
    *)
        echo "Usage: $0 {deploy|cleanup|restart|status|logs}"
        exit 1
        ;;
esac
