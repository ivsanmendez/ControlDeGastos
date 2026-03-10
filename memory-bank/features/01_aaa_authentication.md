# Feature: AAA Framework (Authentication, Authorization & Accounting)

Issue: #5

## Scope
Backend-only implementation. No React UI (separate issue).

### Authentication
- Full lifecycle: Registration → Login → Session Management (refresh token rotation) → Logout
- JWT HS256 access tokens (15 min TTL)
- Random refresh tokens with SHA-256 hashing, stored in DB (7 day TTL)
- Single-use rotation with reuse detection (revoked reuse → revoke ALL user sessions)
- bcrypt password hashing via outbound port

### Authorization
- Permission-based RBAC with 2 roles: `user`, `admin`
- 5 permissions: `expense:create`, `expense:read:own`, `expense:read:all`, `expense:delete:own`, `expense:delete:all`
- Two-level enforcement: middleware (permission check) + service (ownership check)
- Expense scoping via `user_id` FK

### Accounting
- Audit log table in PostgreSQL
- Events: `register`, `login_success`, `login_failed`, `logout`, `token_refresh`
- Fire-and-forget: audit failure never blocks auth operations

## Acceptance Criteria
- [ ] User can register with email/password
- [ ] User can login and receive JWT + refresh token
- [ ] JWT validates on protected routes, returns 401 if missing/invalid
- [ ] Refresh token rotation works (old token revoked, new pair issued)
- [ ] Reuse detection: replaying revoked refresh token revokes ALL user sessions
- [ ] User can logout (refresh token revoked)
- [ ] Expenses scoped to user (user sees only own, admin sees all)
- [ ] Ownership enforced on GET/DELETE (403 for non-owner non-admin)
- [ ] Audit log entries created for all auth events
- [ ] All domain tests pass with fakes
- [ ] `go build`, `go test`, `go vet` all pass
- [ ] Docker image builds

## Architecture
See ADR-03 (`memory-bank/decisions/03_aaa_framework.md`) for full design and rationale.
