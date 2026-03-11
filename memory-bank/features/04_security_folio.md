# Feature: Security Folio for Receipts

## Overview
Each receipt signing event generates a unique, persistent security folio that is included in the signed data, stored in the database, and displayed on the printed receipt. An authenticated verification endpoint lets users look up a folio to validate it.

## Folio Format
```
REC-{YYYY}-{NNNNNN}-{XXXXXXXX}
```
- `REC` — document type prefix (Receipt)
- `YYYY` — year of issuance (signing year, not contribution year)
- `NNNNNN` — 6-digit zero-padded sequential number (per-year, atomic)
- `XXXXXXXX` — 8 uppercase hex chars from `crypto/rand` (uniqueness/anti-guessing)

Example: `REC-2026-000001-A3F7B2C1`

Each sign action creates a **new folio** (even for the same contributor+year), providing a full audit trail.

## Database

### Migration 009
Two tables:
1. **`receipt_folio_counters`** — per-year atomic sequence counter (`year PK`, `last_seq`)
2. **`receipt_folios`** — stores every signed receipt (`folio UNIQUE`, contributor_id, receipt_year, signer_name, user_id, canonical_json, signature, certificate, signed_at)

Sequence generation: `INSERT ... ON CONFLICT DO UPDATE SET last_seq = last_seq + 1 RETURNING last_seq` — atomic, race-condition safe.

## Backend

### Domain: `internal/domain/receipt/`
- **`receipt.go`** — `ReceiptFolio` entity, `GenerateFolio()`, `GenerateUUIDSuffix()`, `ErrNotFound`
- **`service.go`** — `Repository` interface (`NextSequence`, `Save`, `FindByFolio`), `Service` (`GenerateNewFolio`, `SaveFolio`, `VerifyFolio`)

### Adapter: `internal/adapter/postgres/receipt_folio_repo.go`
- `NextSequence` — atomic INSERT ON CONFLICT
- `Save` — INSERT RETURNING id
- `FindByFolio` — SELECT, maps `sql.ErrNoRows` → `receipt.ErrNotFound`

### Ports
- `port/inbound.go` — `ReceiptFolioService` interface
- `port/outbound.go` — `ReceiptFolioRepository` type alias

### Permission
- `user.PermReceiptVerify` (`receipt:verify`) — both `RoleUser` and `RoleAdmin`

### Handler Flow (Updated `receipt_handler.go`)
1. Validate request + extract `claims.UserID`
2. Fetch contributor + contributions
3. **Generate folio** via `receiptSvc.GenerateNewFolio(ctx, year)`
4. Build `receiptData` with `Folio` as first field (deterministic JSON)
5. `json.Marshal` → canonical bytes (folio is part of signed data)
6. `signer.Sign(canonical, password)` → signature
7. **Persist** via `receiptSvc.SaveFolio(ctx, &rf)`
8. Return `{folio, data, signature(base64), certificate(base64)}`

### New Endpoint: `GET /receipts/verify/{folio}`
- Authenticated, requires `receipt:verify` permission
- Returns: folio, contributor_id, receipt_year, signer_name, signed_at, canonical_json (base64), signature (base64), certificate (base64)
- 404 if folio not found

### i18n Keys
- `folio_required`, `receipt_folio_not_found`, `failed_to_generate_folio`, `failed_to_save_receipt`

## Frontend

### Type Changes (`use-receipt-signature.ts`)
- `ReceiptData.folio: string`
- `ReceiptSignatureResponse.folio: string`

### Receipt Page (`contribution-receipt-page.tsx`)
- Folio displayed in header after signing (monospace)
- QR code encodes the folio string (short, scannable)
- Folio text displayed below QR code

### Locales
- `receipt.folio`: "Folio de Seguridad" (ES) / "Security Folio" (EN)

### Vite Proxy
- `/receipts` → `http://localhost:8080`

## Acceptance Criteria
- [x] Migration 009 creates counter + folios tables
- [x] Each sign action generates a unique `REC-YYYY-NNNNNN-XXXXXXXX` folio
- [x] Folio is included in canonical JSON before signing
- [x] Signed receipt + folio persisted in `receipt_folios`
- [x] `GET /receipts/verify/{folio}` returns stored receipt data
- [x] Printed receipt shows folio in header and below QR code
- [x] QR code encodes the folio string
- [x] Go compiles (`go vet`), tests pass, frontend builds
