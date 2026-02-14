# ADR-01: Hexagonal Architecture

## Context
The application needs to scale into an agentic system where AI agents act as first-class adapters alongside the HTTP API, and domain events drive agent orchestration.

## Options Considered
1. **Flat layered architecture** (handler → service → repo) — simple but couples adapters to domain
2. **Clean architecture** — similar to hexagonal but with more layers (use cases, entities, gateways)
3. **Hexagonal architecture (ports & adapters)** — domain core defines port interfaces, adapters plug in

## Decision
Hexagonal architecture with:
- **Outbound ports** (Repository, EventPublisher) defined inside domain package (Go idiom: consumer defines the interface)
- **Inbound ports** (ExpenseService) in a dedicated `port/` package shared by all driving adapters
- **In-memory event bus** as initial EventPublisher, replaceable with NATS/Kafka
- **Single bounded context** (expenses) for now

## Consequences
- AI agents and HTTP handlers are architecturally identical — both call `port.ExpenseService`
- Domain logic is fully testable without infrastructure (fake adapters)
- Adding new adapters (gRPC, CLI, message queue) requires no domain changes
- Slightly more packages than a flat structure, but each has a clear, single responsibility
- Event bus swap (in-memory → NATS) only changes `main.go` wiring and one adapter package