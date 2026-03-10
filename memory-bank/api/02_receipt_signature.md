# Receipt Digital Signature

## Endpoint

### `POST /contributions/receipt-signature`

**Auth**: Required (Bearer JWT) + `contribution:read` permission

Signs a contribution receipt for a given contributor and year using a SAT certificate.

### Request Body
```json
{
  "contributor_id": 42,
  "year": 2026,
  "password": "sat-certificate-password",
  "signer_name": "Juan Perez"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `contributor_id` | int64 | Yes | Contributor to generate receipt for |
| `year` | int | Yes | Fiscal year |
| `password` | string | Yes | Password to decrypt SAT `.key` file |
| `signer_name` | string | Yes | Name of the person signing the receipt |

### Success Response (200)
```json
{
  "data": {
    "contributor_id": 42,
    "house_number": "12A",
    "contributor_name": "Maria Garcia",
    "year": 2026,
    "payments": [
      { "month": 1, "amount": 500.00 },
      { "month": 2, "amount": 500.00 }
    ],
    "total": 1000.00,
    "signer_name": "Juan Perez",
    "generated_at": "2026-03-09T12:00:00Z"
  },
  "signature": "<base64-encoded RSA-SHA256 signature>",
  "certificate": "<base64-encoded DER X.509 certificate>"
}
```

### Error Responses

| Status | Condition |
|--------|-----------|
| 400 | Missing or invalid fields |
| 401 | Invalid certificate password |
| 404 | Contributor not found |
| 503 | Signing not configured (no `SIGN_CERT_PATH` / `SIGN_KEY_PATH`) |

## Configuration

| Env Variable | Description |
|-------------|-------------|
| `SIGN_CERT_PATH` | Path to `.cer` file (DER or PEM X.509 certificate) |
| `SIGN_KEY_PATH` | Path to `.key` file (DER-encoded encrypted PKCS#8, SAT format) |

Both must be set for signing to be available. If either is missing, endpoint returns 503.

## Signing Algorithm

1. Build canonical JSON of `receiptData` (includes `signer_name`)
2. SHA-256 hash
3. RSA PKCS#1 v1.5 signature
4. Response includes base64-encoded signature + certificate for verification

## SAT Key Format

Mexican SAT certificates use DER-encoded encrypted PKCS#8 for private keys (`.key` files). The `youmark/pkcs8` library handles decryption. The key is decrypted per-request — the password is never stored server-side.
