package user

import "time"

type AuditAction string

const (
	AuditLoginSuccess AuditAction = "login_success"
	AuditLoginFailed  AuditAction = "login_failed"
	AuditLogout       AuditAction = "logout"
	AuditTokenRefresh AuditAction = "token_refresh"
	AuditRegister     AuditAction = "register"
)

type AuditEntry struct {
	ID        int64
	UserID    *int64
	Action    AuditAction
	IP        string
	UserAgent string
	Metadata  map[string]string
	CreatedAt time.Time
}

// AuditInfo carries request metadata for audit logging.
type AuditInfo struct {
	IP        string
	UserAgent string
}

// NewAuditEntry creates an AuditEntry with the current timestamp.
func NewAuditEntry(userID *int64, action AuditAction, info AuditInfo, metadata map[string]string) AuditEntry {
	return AuditEntry{
		UserID:    userID,
		Action:    action,
		IP:        info.IP,
		UserAgent: info.UserAgent,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}
}
