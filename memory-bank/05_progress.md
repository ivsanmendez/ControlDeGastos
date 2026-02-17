# Progress Log

## What Works
- Hexagonal architecture implemented (domain/port/adapter layers)
- Domain core: Expense entity with factory + validation, domain events, Service with CRUD use cases
- Port interfaces: `ExpenseService` (inbound), `Repository`/`EventPublisher` (outbound), `EventSubscriber`
- HTTP adapter: health endpoint + expense CRUD handlers wired via `RegisterRoutes`
- PostgreSQL adapter: connection helper + repository stubs (not yet implemented)
- In-memory event bus adapter
- Composition root in `cmd/api/main.go` wires all layers
- React SPA scaffolded with Vite + TypeScript (`web/`)
- Docker Compose (API + PostgreSQL 16 + React dev server)
- GitHub Actions CI pipeline
- Memory bank with SessionStart hook auto-loading
- CLAUDE.md
- Production deployment via Podman pod + Cloudflare Named Tunnel
  - Public URL: https://cdg.meyis.work
  - Script: `./deploy-alt-ports.sh`
  - Pod: API + PostgreSQL + cloudflared containers

## What's Left to Build
- PostgreSQL repository implementation (actual SQL)
- Database migrations
- React expense management UI
- Domain unit tests
- Authentication (TBD)
- AI agent adapter (future)

## Known Issues
- PostgreSQL adapter methods are stubs (return nil/ErrNotFound)
- No database migrations yet