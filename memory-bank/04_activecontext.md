# Active Context

## Current Phase
Hexagonal architecture restructuring complete. Domain core, ports, and adapter scaffolding are in place.

## Recent Decisions
- Hexagonal architecture (ports & adapters) for the Go backend
- Single bounded context (expenses) — expandable to multiple later
- Outbound port interfaces defined in domain package (Go idiom: consumer defines interface)
- Inbound port interfaces in dedicated `port/` package (shared by HTTP + future agent adapters)
- In-memory synchronous event bus (replaceable with NATS/Kafka later)
- `github.com/lib/pq` as PostgreSQL driver

## Next Steps
- [ ] Implement PostgreSQL repository (actual SQL queries)
- [ ] Database migration strategy (golang-migrate, goose, etc.)
- [ ] Build React UI for expense management
- [ ] Add domain unit tests with fake adapters

## Open Questions
- Authentication approach (if needed)
- API documentation tooling (OpenAPI/Swagger)
- When to introduce AI agent adapter