#!/bin/bash
# ControlDeGastos - Cloudflare Named Tunnel Deployment
# Public URL: https://cdg.meyis.work
# Requires TUNNEL_TOKEN in .env.production

set -e

PUBLIC_URL="https://cdg.meyis.work"

echo "🚀 ControlDeGastos Deployment (Cloudflare Tunnel)"
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
    podman stop controldegastos-api controldegastos-db controldegastos-cloudflared 2>/dev/null || true
    podman rm controldegastos-api controldegastos-db controldegastos-cloudflared 2>/dev/null || true
    podman pod stop controldegastos 2>/dev/null || true
    podman pod rm controldegastos 2>/dev/null || true
}

# Function to deploy
deploy() {
    echo "📦 Building application image..."
    podman build -t controldegastos:latest -f Dockerfile .

    echo "🔧 Creating pod..."
    podman pod create --name controldegastos \
        -p 8080:8080

    echo "🐘 Starting PostgreSQL..."
    podman run -d --pod controldegastos --name controldegastos-db \
        --env-file .env.production \
        -v controldegastos_pgdata:/var/lib/postgresql/data \
        docker.io/library/postgres:16-alpine

    echo "⏳ Waiting for database to be ready..."
    sleep 10

    echo "🚀 Starting API..."
    podman run -d --pod controldegastos --name controldegastos-api \
        --env-file .env.production \
        -e STATIC_DIR=/web/dist \
        controldegastos:latest

    echo "🌐 Starting Cloudflare Tunnel..."
    podman run -d --pod controldegastos --name controldegastos-cloudflared \
        --env-file .env.production \
        docker.io/cloudflare/cloudflared:latest \
        tunnel run

    echo
    echo "✅ Deployment complete!"
    echo
    echo "📊 Status:"
    podman ps --pod --filter pod=controldegastos
    echo
    echo "🌐 Application: $PUBLIC_URL"
    echo "   Local:       http://localhost:8080"
    echo
    echo "📝 Useful commands:"
    echo "  View API logs:    podman logs -f controldegastos-api"
    echo "  View tunnel logs: podman logs -f controldegastos-cloudflared"
    echo "  Stop:             podman pod stop controldegastos"
    echo "  Start:            podman pod start controldegastos"
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
        podman pod restart controldegastos
        echo "✅ Restart complete!"
        echo "   Application: $PUBLIC_URL"
        ;;
    status)
        echo "📊 Status:"
        podman ps --pod --filter pod=controldegastos
        echo
        echo "🌐 Application: $PUBLIC_URL"
        ;;
    logs)
        podman logs -f controldegastos-api
        ;;
    *)
        echo "Usage: $0 {deploy|cleanup|restart|status|logs}"
        exit 1
        ;;
esac
