# Auth API

## Endpoints

| Method | Path | Auth | Request Body | Response |
|--------|------|------|-------------|----------|
| POST | `/auth/register` | Public | `{email, password}` | 201 `{id, email, role}` |
| POST | `/auth/login` | Public | `{email, password}` | 200 `{access_token, refresh_token}` |
| POST | `/auth/refresh` | Public | `{refresh_token}` | 200 `{access_token, refresh_token}` |
| POST | `/auth/logout` | Bearer | `{refresh_token}` | 204 |
| GET | `/auth/me` | Bearer | — | 200 `{id, email, role}` |

## Authentication

- **Access token**: JWT HS256, 15 min TTL
  - Header: `Authorization: Bearer <jwt>`
  - Claims: `uid` (int64), `email` (string), `role` (string), `exp`, `iat`
- **Refresh token**: Random 32 bytes hex, 7 day TTL, single-use with rotation

## Request/Response Examples

### Register
```json
// POST /auth/register
// Request
{"email": "user@example.com", "password": "password123"}

// Response 201
{"id": 1, "email": "user@example.com", "role": "user"}
```

### Login
```json
// POST /auth/login
// Request
{"email": "user@example.com", "password": "password123"}

// Response 200
{"access_token": "eyJhbGci...", "refresh_token": "a1b2c3d4..."}
```

### Refresh
```json
// POST /auth/refresh
// Request
{"refresh_token": "a1b2c3d4..."}

// Response 200
{"access_token": "eyJhbGci...", "refresh_token": "e5f6g7h8..."}
```

### Logout
```json
// POST /auth/logout (requires Bearer token)
// Request
{"refresh_token": "a1b2c3d4..."}

// Response 204 (no body)
```

### Me
```json
// GET /auth/me (requires Bearer token)
// Response 200
{"id": 1, "email": "user@example.com", "role": "user"}
```

## Error Responses

All errors return `{"error": "message"}`.

| Status | Meaning |
|--------|---------|
| 400 | Invalid request body |
| 401 | Invalid/expired credentials or token |
| 403 | Token revoked (reuse detection) |
| 409 | Email already taken (register) |
| 422 | Validation error (weak password, invalid email) |

## Expense Endpoints (Updated)

All expense endpoints now require authentication and enforce ownership:

| Method | Path | Auth | Permission |
|--------|------|------|------------|
| POST | `/expenses` | Bearer | `expense:create` |
| GET | `/expenses` | Bearer | `expense:read:own` (user: own only, admin: all) |
| GET | `/expenses/{id}` | Bearer | `expense:read:own` (403 if not owner and not admin) |
| DELETE | `/expenses/{id}` | Bearer | `expense:delete:own` (403 if not owner and not admin) |
