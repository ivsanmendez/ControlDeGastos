package i18n

// MessagesES contains all Spanish translations for API error messages.
var MessagesES = map[string]string{
	// Generic
	"invalid_request_body": "cuerpo de solicitud inválido",
	"invalid_id":           "id inválido",

	// Auth
	"registration_failed":              "error en el registro",
	"login_failed":                     "error al iniciar sesión",
	"invalid_or_expired_refresh_token": "token de actualización inválido o expirado",
	"token_refresh_failed":             "error al actualizar el token",
	"invalid_refresh_token":            "token de actualización inválido",
	"logout_failed":                    "error al cerrar sesión",
	"no_claims_in_context":             "sin claims en el contexto",
	"user_not_found":                   "usuario no encontrado",

	// Auth middleware
	"missing_or_invalid_auth_header": "encabezado de autorización faltante o inválido",
	"invalid_or_expired_token":       "token inválido o expirado",
	"insufficient_permissions":       "permisos insuficientes",

	// Expenses
	"expense_not_found": "gasto no encontrado",

	// Contributors
	"contributor_not_found": "contribuyente no encontrado",

	// Contributions
	"invalid_contributor_id":       "contributor_id inválido",
	"invalid_year":                 "año inválido",
	"invalid_payment_date_format":  "formato de payment_date inválido, se esperaba YYYY-MM-DD",
	"contribution_not_found":       "contribución no encontrada",

	// Receipt
	"receipt_signing_not_configured":    "la firma de recibos no está configurada",
	"contributor_id_and_year_required":  "contributor_id y year son requeridos",
	"password_required":                 "la contraseña es requerida",
	"signer_name_required":             "el nombre del firmante es requerido",
	"failed_to_load_contributions":     "no se pudieron cargar las contribuciones",
	"failed_to_serialize_receipt_data":  "no se pudieron serializar los datos del recibo",
	"invalid_certificate_password":     "contraseña de certificado inválida",
	"failed_to_sign_receipt":           "no se pudo firmar el recibo",
	"folio_required":                   "el folio es requerido",
	"receipt_folio_not_found":          "folio de recibo no encontrado",
	"failed_to_generate_folio":         "no se pudo generar el folio",
	"failed_to_save_receipt":           "no se pudo guardar el recibo",

	// Expense categories
	"expense_category_not_found": "categoría de gasto no encontrada",

	// Reports
	"report_query_failed": "error al generar el reporte",
}
