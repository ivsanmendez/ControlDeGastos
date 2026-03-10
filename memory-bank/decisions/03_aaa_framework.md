# ADR-03: AAA Framework — Authentication, Authorization & Accounting

## Context
The application needs user identity and access control. Expenses must be scoped to individual users, and an audit trail is required for security-sensitive operations. Issue #5.

## Options Considered

### Authentication
| Option | Pros | Cons |
|--------|------|------|
| JWT + Refresh tokens | Stateless access, SPA-friendly, rotation enables revocation | Must store refresh tokens in DB |
| Session cookies | Simple, server-side revocation | Requires session store, less API-friendly |
| OAuth2 (external) | Delegate identity | Over-engineered for personal app |

### Authorization
| Option | Pros | Cons |
|--------|------|------|
| Hardcoded RBAC (role→permissions map) | Simple, testable, no DB table | Adding roles requires code change |
| DB-backed RBAC | Dynamic role management | Overkill for 2 roles |

### Password Hashing
| Option | Pros | Cons |
|--------|------|------|
| bcrypt via outbound port | Domain stays dependency-free, industry standard | Slightly slower than argon2 |
| argon2 | Newer, more configurable | More complex setup |

## Decision

### Authentication
- **Access token**: JWT HS256, 15 min TTL — stateless, fits SPA + API pattern
- **Refresh token**: Random 32 bytes hex, SHA-256 hashed before DB storage, 7 day TTL
- **Refresh strategy**: Single-use rotation with reuse detection — if a revoked token is reused, ALL user sessions are revoked (indicates token theft)
- **Password hashing**: bcrypt via `PasswordHasher` outbound port (domain stays dependency-free)

### Authorization
- **Model**: Permission-based RBAC with 2 roles (`user`, `admin`)
- **Permissions**: `expense:create`, `expense:read:own`, `expense:read:all`, `expense:delete:own`, `expense:delete:all`
- **Role map**: Hardcoded in domain — `user` gets `*:own`, `admin` gets `*:own` + `*:all`
- **Enforcement**: Two levels:
  1. **Middleware** (coarse): Does the role have the required permission?
  2. **Service** (fine): Does the caller own the specific resource? (via `callerID`/`callerRole` params)

### Accounting (Audit)
- **Audit log**: PostgreSQL `audit_logs` table
- **Events**: `register`, `login_success`, `login_failed`, `logout`, `token_refresh`
- **Strategy**: Fire-and-forget — audit failure never blocks auth operations (errors logged to stdout)

### Expense Scoping
- `user_id` FK added to `expenses` table
- Service methods accept `callerID`/`callerRole` for ownership enforcement

## Implementation Phases

| Phase | Scope | Status |
|-------|-------|--------|
| 1 | Add Go dependencies (`x/crypto`, `golang-jwt/jwt/v5`) | Done |
| 2 | User domain (`internal/domain/user/`) — entity, tokens, permissions, audit, ports, service, tests | Done |
| 3 | Update expense domain with `UserID` + caller params | Done |
| 4 | Update port interfaces (`auth.go`, `inbound.go`, `outbound.go`) | Done |
| 5 | Database migrations (users, refresh_tokens, audit_logs, expenses.user_id) | Done |
| 6 | Driven adapters (postgres user/audit repos, bcrypt hasher, JWT issuer, update expense repo) | Done |
| 7 | HTTP adapter (auth handlers, middleware, update routes + expense handler) | Done |
| 8 | Composition root + config (main.go, mise, docker-compose, CI) | Done |

## File Impact

### New Files (16)
- `internal/domain/user/` — `user.go`, `token.go`, `permission.go`, `audit.go`, `ports.go`, `service.go`, `user_test.go`, `service_test.go`
- `internal/port/auth.go`
- `db/migrations/` — `002_create_users.sql`, `003_create_refresh_tokens.sql`, `004_add_user_id_to_expenses.sql`, `005_create_audit_logs.sql`
- `internal/adapter/postgres/user_repo.go`, `internal/adapter/postgres/audit_repo.go`
- `internal/adapter/bcrypt/hasher.go`, `internal/adapter/jwt/issuer.go`
- `internal/adapter/httpapi/auth_handler.go`, `internal/adapter/httpapi/auth_middleware.go`

### Modified Files (10)
- `internal/domain/expense/` — `expense.go`, `service.go`, `expense_test.go`, `service_test.go`
- `internal/port/` — `inbound.go`, `outbound.go`
- `internal/adapter/httpapi/` — `router.go`, `expense_handler.go`
- `internal/adapter/postgres/expense_repo.go`
- `cmd/api/main.go`, `go.mod`, `docker-compose.yml`, `.github/workflows/ci.yml`, `.mise.toml`

## Refresh Token Rotation Flow

```
Client                            Server
  │ POST /auth/login ────────────►│ verify credentials
  │                               │ issue JWT (15min) + refresh (random)
  │                               │ store SHA-256(refresh) in DB
  │                               │ audit: login_success
  │◄── {access_token, refresh} ───┤
  │                               │
  │  ... 15 min later ...         │
  │                               │
  │ POST /auth/refresh ──────────►│ SHA-256(incoming) → DB lookup
  │  {refresh_token: "old"}       │ revoked? → REVOKE ALL → 403
  │                               │ expired? → 401
  │                               │ valid → revoke old, issue new pair
  │                               │ audit: token_refresh
  │◄── {access_token, refresh} ───┤
  │                               │
  │ POST /auth/logout ───────────►│ revoke refresh token
  │                               │ audit: logout
  │◄── 204 ───────────────────────┤
```

## Authorization Flow (per request)

```
Request → RequireAuth (validate JWT, inject Claims)
        → RequirePermission (role has permission?)
        → Handler (extract callerID/role from Claims)
        → Service (ownership check: own resource or admin?)
        → Repository
```

## Consequences
- Go upgraded from 1.23 to 1.24 (required by `golang.org/x/crypto@v0.48.0`)
- Domain-to-domain import: `expense` package imports `user.Role` for ownership checks
- All existing expense API calls become authenticated — breaking change for any existing clients
- Refresh token table grows over time — future cleanup job may be needed
