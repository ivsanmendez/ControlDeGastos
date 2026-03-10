# Architecture

## Hexagonal Architecture (Ports & Adapters)

The Go backend follows hexagonal architecture to support future AI agent adapters and event-driven orchestration.

```
┌─────────────────────────────────────────────────────────────┐
│                    cmd/api/main.go                          │
│                  (Composition Root)                         │
└────────┬──────────────────────────────────┬─────────────────┘
         │                                  │
┌────────▼────────────┐         ┌───────────▼──────────────┐
│  DRIVING ADAPTERS   │         │   DRIVEN ADAPTERS        │
│  (inbound)          │         │   (outbound)             │
│                     │         │                          │
│  adapter/httpapi/   │         │   adapter/postgres/      │
│  [future] agent/    │         │   adapter/eventbus/      │
│  [future] grpc/     │         │   adapter/certsigner/    │
│                     │         │   adapter/bcrypt/        │
│                     │         │   adapter/jwt/           │
│                     │         │   [future] nats/         │
└────────┬────────────┘         └───────────┬──────────────┘
         │ depends on                       │ implements
         ▼                                  ▼
┌──────────────────────────────────────────────────────────┐
│                     PORTS                                 │
│  port/inbound.go    → ExpenseService, AuthService, etc.  │
│  domain/expense/    → Repository, EventPublisher ifaces  │
│  port/outbound.go   → EventSubscriber, ReceiptSigner     │
└──────────────────────────┬───────────────────────────────┘
                           │
              ┌────────────▼──────────────┐
              │       DOMAIN CORE         │
              │   domain/expense/         │
              │                           │
              │   Entity + factory        │
              │   Domain events           │
              │   Service (use cases)     │
              │   Zero external deps      │
              └───────────────────────────┘
```

## Repository Layout
```
ControlDeContabilidad/
├── cmd/api/main.go              # Composition root (wires all adapters)
├── internal/
│   ├── domain/expense/          # Core hexagon (zero external deps)
│   │   ├── expense.go           # Entity, factory, domain errors
│   │   ├── event.go             # Domain events
│   │   └── service.go           # Service + outbound port interfaces
│   ├── port/
│   │   ├── inbound.go           # Driving port (ExpenseService interface)
│   │   └── outbound.go          # EventSubscriber + type aliases
│   └── adapter/
│       ├── httpapi/             # HTTP driving adapter
│       ├── postgres/            # PostgreSQL driven adapter
│       ├── eventbus/            # In-memory event bus
│       ├── certsigner/          # SAT certificate signer (encrypted PKCS#8)
│       ├── bcrypt/              # Password hashing
│       └── jwt/                 # JWT token issuance
├── web/                         # React SPA (Vite + TypeScript)
├── memory-bank/                 # Project documentation
├── .github/workflows/           # GitHub Actions CI/CD
├── .claude/hooks/               # Claude Code session hooks
├── Dockerfile                   # Multi-stage Go API image
├── docker-compose.yml           # Full stack orchestration
└── CLAUDE.md
```

## Service Architecture
```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│  React SPA  │────▶│   Go API     │────▶│ PostgreSQL │
│  (Vite)     │     │  (net/http)  │     │            │
│  :5173      │     │  :8080       │     │  :5432     │
└─────────────┘     └──────────────┘     └────────────┘
                    ┌──────────────┐
                    │  Event Bus   │ (in-memory, future: NATS)
                    └──────────────┘
```

## Dependency Rule
Arrows always point inward. Domain core imports nothing from adapters or ports. Only `main.go` knows all concrete types.

## Docker Strategy
- **Development**: `docker compose up` runs all three services
- **Production**: Multi-stage Dockerfile builds a minimal Go binary image
- PostgreSQL uses a named volume for data persistence

## Future: Agentic System
- AI agents will be driving adapters in `adapter/agent/`, calling the same `port.ExpenseService` as HTTP handlers
- Event subscribers react to domain events via the event bus
- The in-memory bus can be swapped to NATS/Kafka by adding a new adapter