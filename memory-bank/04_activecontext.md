# Active Context

## Current Phase
SAT certificate signing + print-sign dialog flow implemented. Contributions receipt system now supports Mexican SAT `.cer`/`.key` format with per-request password decryption and a browser dialog for signer name + password before printing.

## Recent Decisions
- Hexagonal architecture (ports & adapters) for the Go backend
- Single bounded context (expenses) — expandable to multiple later
- Outbound port interfaces defined in domain package (Go idiom: consumer defines interface)
- Inbound port interfaces in dedicated `port/` package (shared by HTTP + future agent adapters)
- In-memory synchronous event bus (replaceable with NATS/Kafka later)
- `github.com/lib/pq` as PostgreSQL driver
- goose for database migrations (`db/migrations/`)
- Convention-aligned PR template (ADR-02)
- Production deployment: Podman pod + Cloudflare Named Tunnel
- AAA framework design (ADR-03):
  - JWT HS256 access tokens (15 min) + refresh token rotation (7 day, SHA-256 hashed)
  - Permission-based RBAC (2 roles: user/admin, hardcoded role→permission map)
  - Audit log to PostgreSQL (fire-and-forget)
  - bcrypt via outbound port (domain stays dependency-free)
  - Go upgraded from 1.23 → 1.24 (required by `golang.org/x/crypto@v0.48.0`)
- SAT certificate signing:
  - `youmark/pkcs8` for DER-encoded encrypted PKCS#8 (SAT `.key` format)
  - Private key decrypted per-request (password not stored server-side)
  - `ReceiptSigner.Sign(data, password)` port interface
  - Receipt endpoint changed from GET to POST (accepts password + signer_name)
  - Frontend: print-sign dialog flow (dialog → sign → render QR + name → print)

## Next Steps
- [x] Implement PostgreSQL repository (actual SQL queries) — #1
- [x] Database migration strategy — #2 (goose)
- [x] Add domain unit tests with fake adapters — #3
- [x] AAA framework — #5
- [x] SAT certificate signing + print-sign dialog
- [ ] Build React UI for expense management — #4

## Open Questions
- API documentation tooling (OpenAPI/Swagger)
- When to introduce AI agent adapter
