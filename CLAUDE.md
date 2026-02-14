# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

ControlDeGastos — full-stack expense tracking app. Go REST API + React SPA + PostgreSQL, all orchestrated via Docker Compose.

## Memory Bank

Project documentation lives in `memory-bank/`. Each folder has a `README.md` as its index.

A **SessionStart hook** (`.claude/hooks/load-memory-bank.sh`) automatically loads all memory bank files into context at the beginning of every session and after context compaction. No manual reading needed.

Files use `NN_name.md` numbering — lower number = higher priority:
1. `01_projectbrief.md` — scope and goals
2. `02_techstack.md` — tools, versions, ports
3. `03_architecture.md` — repo layout and service diagram
4. `04_activecontext.md` — current phase, decisions, next steps
5. `05_progress.md` — what works, what's left, known issues
6. `06_kanban.md` — Kanban board state and GitHub Project link

Subdirectories for growing documentation:
- `decisions/` — Architecture Decision Records
- `features/` — feature design docs
- `api/` — endpoint specs and schemas
- `frontend/` — component and state docs

Update these files when making significant changes. Use the `NN_` naming convention when adding new files.

## Tool Management

Uses **mise** (`.mise.toml`) to manage Go 1.23 and Node 22. All commands should be prefixed with `mise exec --` or use `mise run <task>`.

### Environment Profiles
```
.mise.toml              # Base: tools + shared env vars (PORT, APP_NAME) + tasks
.mise.development.toml  # Dev: DATABASE_URL → localhost, LOG_LEVEL=debug
.mise.test.toml         # Test: DATABASE_URL → test DB, LOG_LEVEL=warn
.mise.local.toml        # Local secrets (gitignored)
```
Activate a profile: `export MISE_ENV=development` (or `test`). Mise merges the profile file on top of the base config. Env vars auto-activate when entering the project directory.

## Commands

### Docker (primary workflow)
```bash
docker compose up              # Start all services (API + DB + frontend)
docker compose up -d           # Start detached
docker compose build           # Rebuild images
docker compose down            # Stop all services
docker compose down -v         # Stop and remove volumes (resets DB)
```

### mise tasks
```bash
mise run build                 # Build Go API binary to bin/api
mise run dev:api               # Run Go API locally (port 8080)
mise run dev:web               # Run Vite dev server (port 5173)
mise run test                  # Run all Go tests
mise run test:web              # Run frontend tests
mise run lint                  # go vet ./...
mise run lint:web              # ESLint frontend
```

### Running a single Go test
```bash
mise exec -- go test -run TestName ./internal/domain/expense/...
```

### Running frontend commands
```bash
mise exec -- npm --prefix web run dev
mise exec -- npm --prefix web run build
```

## Architecture (Hexagonal)

- **`cmd/api/main.go`** — Composition root (wires all adapters, only place that knows concrete types)
- **`internal/domain/expense/`** — Core hexagon: entity, factory, domain events, Service, outbound port interfaces (zero external deps)
- **`internal/port/`** — `inbound.go` (ExpenseService driving port), `outbound.go` (EventSubscriber + type aliases)
- **`internal/adapter/httpapi/`** — HTTP driving adapter (handlers call `port.ExpenseService`)
- **`internal/adapter/postgres/`** — PostgreSQL driven adapter (implements `expense.Repository`)
- **`internal/adapter/eventbus/`** — In-memory event bus (implements `expense.EventPublisher`)
- **`web/`** — React SPA (Vite + TypeScript)
- Dependency rule: domain imports nothing from adapters. Adapters depend on ports. Only `main.go` imports everything.
- Go 1.22+ enhanced `net/http` mux for routing (no external router)
- API serves JSON on `:8080`, React dev server on `:5173`

## Docker

Container-first development. PostgreSQL runs in Docker always. The Go API has a multi-stage Dockerfile (builder + minimal alpine runtime).

Environment variable `DATABASE_URL` connects the API to PostgreSQL.

## Kanban Workflow

Project follows Kanban methodology. Board: https://github.com/users/ivsanmendez/projects/2

**Columns**: Backlog → Todo → In Progress → Review → Done

Work items are GitHub issues. When starting work:
1. Move the issue to "In Progress" on the board
2. Create a feature branch from `main`
3. When done, move to "Review" and open a PR
4. After merge, issue moves to "Done"

Update `memory-bank/06_kanban.md` after board state changes.

**Labels**: `backend`, `frontend`, `infrastructure`, `domain`, `agentic`

## Commit Convention

All commits **must** use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): short description
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `build`

**Scopes**: `api`, `domain`, `web`, `infra`, `agent`, `memory-bank`

**Examples**:
```
feat(api): add expense CRUD handlers
fix(domain): validate negative amounts
docs(memory-bank): update architecture diagram
chore(infra): add database migration task to mise
ci: add lint step for conventional commits
test(domain): add service unit tests with fakes
```

Breaking changes use `!` after scope: `feat(api)!: change expense response format`

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`) runs on push/PR to main:
- **api** job: `go vet`, `go test -race`, `go build` (with PostgreSQL service)
- **web** job: `npm ci`, `npm run lint`, `npm run build`
- **docker** job: verifies `docker build` succeeds