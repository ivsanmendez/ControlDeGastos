# SAT Certificate Signing + Print-Sign Dialog

## Scope
Digital receipt signing using Mexican SAT certificates (`.cer` + `.key`), with a frontend dialog that collects the signer's name and certificate password before printing.

## Problem
- SAT `.key` files are DER-encoded encrypted PKCS#8 (password-protected, not PEM)
- Go stdlib doesn't support encrypted PKCS#8
- Password should not be stored server-side — must be provided per-signing operation
- User needs to enter signer name + password before printing a receipt

## Solution

### Backend
- **`certsigner` adapter** rewired: stores raw encrypted key bytes, decrypts per `Sign()` call
- **`youmark/pkcs8`** library handles DER encrypted PKCS#8 decryption
- **`ReceiptSigner` port**: `Sign(data []byte, password string) ([]byte, error)`
- **Receipt endpoint** changed to `POST /contributions/receipt-signature` with JSON body:
  - `contributor_id`, `year`, `password`, `signer_name`
  - `signer_name` is included in the signed data payload
  - Wrong password returns 401

### Frontend
- **Print button** opens a `ReceiptSignDialog` instead of calling `window.print()` directly
- **Dialog fields**: Signer Name (text) + Certificate Password (password)
- On submit: POST to backend with password + signer name
- On success:
  1. Store signer name + signature response in React state
  2. Close dialog
  3. Render signer name on the "Name: ___" line in signature area
  4. Render `<QRCodeSVG>` with signed payload
  5. `window.print()` after 300ms (let React render first)
- On error: show error message in dialog (e.g. "Invalid certificate password")

### Key Format Support
| Format | Source | Supported |
|--------|--------|-----------|
| DER encrypted PKCS#8 | SAT `.key` | Yes (primary) |
| PEM encrypted PKCS#8 | Generic | Yes (fallback) |
| PEM/DER unencrypted PKCS#8 | Generic | Yes (fallback) |
| PEM/DER PKCS#1 | Generic | Yes (fallback) |

## Files Changed
- `internal/adapter/certsigner/signer.go` — SAT format support, per-call decryption
- `internal/port/outbound.go` — `Sign` takes password param
- `internal/adapter/httpapi/receipt_handler.go` — POST, password + signer_name
- `internal/adapter/httpapi/router.go` — POST route
- `web/src/hooks/use-receipt-signature.ts` — `useMutation` instead of `useQuery`
- `web/src/pages/contribution-receipt-page.tsx` — dialog trigger, QR + signer name state
- `web/src/components/contributions/receipt-sign-dialog.tsx` — new dialog component

## Configuration
- `SIGN_CERT_PATH` — path to `.cer` file
- `SIGN_KEY_PATH` — path to `.key` file
- Both must be set for signing; otherwise endpoint returns 503 and dialog still opens but fails gracefully
