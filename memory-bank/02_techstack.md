# Tech Stack & Tooling

## Runtime & Languages
| Tool       | Version | Managed by |
|------------|---------|------------|
| Go         | 1.23    | mise       |
| Node.js    | 22      | mise       |
| PostgreSQL | 16      | Docker     |

## Mise Environment Profiles
| File | Purpose | Committed |
|------|---------|-----------|
| `.mise.toml` | Base: tools, shared env vars, tasks | Yes |
| `.mise.development.toml` | Dev: local DB URL, debug logging | Yes |
| `.mise.test.toml` | Test: test DB URL, warn logging | Yes |
| `.mise.local.toml` | Personal secrets/overrides | No (gitignored) |

Activate with `export MISE_ENV=development` (or `test`).

## Key Decisions
- **mise** manages language runtimes (Go, Node) and environment profiles via `.mise.toml`
- **Docker Compose** is the primary way to run the full stack (API + DB + frontend)
- **GitHub Actions** for CI/CD (build, test, lint)
- Go standard library `net/http` for routing (Go 1.22+ enhanced mux)
- Vite for frontend build tooling

## Ports (default)
| Service    | Port |
|------------|------|
| Go API     | 8080 |
| React dev  | 5173 |
| PostgreSQL | 5432 |