# Progress Log

## What Works
- Hexagonal architecture implemented (domain/port/adapter layers)
- Domain core: Expense entity with factory + validation, domain events, Service with CRUD use cases
- Port interfaces: `ExpenseService` (inbound), `Repository`/`EventPublisher` (outbound), `EventSubscriber`
- HTTP adapter: health endpoint + expense CRUD handlers wired via `RegisterRoutes`
- PostgreSQL adapter: full CRUD implementation (Save, FindByID, FindAll, Delete)
- In-memory event bus adapter
- Composition root in `cmd/api/main.go` wires all layers
- React SPA scaffolded with Vite + TypeScript (`web/`)
- Docker Compose (API + PostgreSQL 16 + React dev server)
- GitHub Actions CI pipeline
- Memory bank with SessionStart hook auto-loading
- CLAUDE.md
- Database migrations with goose (`db/migrations/`)
- Domain unit tests (entity + service with fake adapters)
- Convention-aligned PR template (`.github/pull_request_template.md`)
- Production deployment via Podman pod + Cloudflare Named Tunnel
  - Public URL: https://cdg.meyis.work
  - Script: `./deploy.sh`
  - Pod: API + PostgreSQL + cloudflared containers
- Contributor + contribution management (CRUD endpoints + React UI)
- Receipt digital signature system (certsigner adapter, SAT certificate support)

## Recently Completed — SAT Certificate Signing + Print-Sign Dialog
- **certsigner adapter** rewritten for SAT format:
  - Supports DER-encoded encrypted PKCS#8 (`.key`), PEM encrypted PKCS#8, and unencrypted PKCS#8/PKCS#1 fallback
  - Private key stored as raw encrypted bytes, decrypted per `Sign()` call with password
  - Added `github.com/youmark/pkcs8` dependency
- **ReceiptSigner port** updated: `Sign(data []byte, password string) ([]byte, error)`
- **Receipt endpoint** changed from `GET` to `POST /contributions/receipt-signature`:
  - Request body: `{ contributor_id, year, password, signer_name }`
  - `signer_name` included in signed data
  - Returns 401 on wrong password, 503 if signer not configured
- **Frontend print-sign dialog** (`receipt-sign-dialog.tsx`):
  - Dialog with Signer Name + Certificate Password fields
  - Uses `useMutation` (was `useQuery`) — triggered on demand
  - On success: renders signer name on signature line + QR code with signed payload, then `window.print()`

## Previously Completed — AAA Framework (#5)
- User domain: entity, tokens, permissions, audit, ports, service + 24 tests
- Expense domain updated: `UserID` field, caller params, ownership checks + 18 tests
- Port interfaces: `AuthService` driving port, updated `ExpenseService`
- Database migrations: users, refresh_tokens, audit_logs, expenses.user_id
- Driven adapters: postgres user/audit repos, bcrypt hasher, JWT issuer, updated expense repo
- HTTP adapter: auth handlers, middleware (RequireAuth, RequirePermission), updated routes
- Composition root wires all new services
- Config: JWT_SECRET in mise, docker-compose, CI
- Go upgraded 1.23 → 1.24, Dockerfile + CI + mise updated

## What's Left to Build
- React expense management UI (#4)
- AI agent driving adapter (#6)
- Persistent event bus (#7)

## Known Issues
_(none)_
