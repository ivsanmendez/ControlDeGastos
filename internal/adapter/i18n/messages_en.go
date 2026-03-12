package i18n

// MessagesEN contains all English translations for API error messages.
var MessagesEN = map[string]string{
	// Generic
	"invalid_request_body": "invalid request body",
	"invalid_id":           "invalid id",

	// Auth
	"registration_failed":              "registration failed",
	"login_failed":                     "login failed",
	"invalid_or_expired_refresh_token": "invalid or expired refresh token",
	"token_refresh_failed":             "token refresh failed",
	"invalid_refresh_token":            "invalid refresh token",
	"logout_failed":                    "logout failed",
	"no_claims_in_context":             "no claims in context",
	"user_not_found":                   "user not found",

	// Auth middleware
	"missing_or_invalid_auth_header": "missing or invalid authorization header",
	"invalid_or_expired_token":       "invalid or expired token",
	"insufficient_permissions":       "insufficient permissions",

	// Expenses
	"expense_not_found": "expense not found",

	// Contributors
	"contributor_not_found": "contributor not found",

	// Contributions
	"invalid_contributor_id":       "invalid contributor_id",
	"invalid_year":                 "invalid year",
	"invalid_payment_date_format":  "invalid payment_date format, expected YYYY-MM-DD",
	"contribution_not_found":       "contribution not found",

	// Receipt
	"receipt_signing_not_configured":    "receipt signing is not configured",
	"contributor_id_and_year_required":  "contributor_id and year are required",
	"password_required":                 "password is required",
	"signer_name_required":             "signer_name is required",
	"failed_to_load_contributions":     "failed to load contributions",
	"failed_to_serialize_receipt_data":  "failed to serialize receipt data",
	"invalid_certificate_password":     "invalid certificate password",
	"failed_to_sign_receipt":           "failed to sign receipt",
	"folio_required":                   "folio is required",
	"receipt_folio_not_found":          "receipt folio not found",
	"failed_to_generate_folio":         "failed to generate folio",
	"failed_to_save_receipt":           "failed to save receipt",

	// Expense categories
	"expense_category_not_found": "expense category not found",

	// Reports
	"report_query_failed": "report query failed",
}
