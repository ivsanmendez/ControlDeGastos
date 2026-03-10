package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// AuditRepo implements user.AuditLogger.
type AuditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Log(ctx context.Context, entry user.AuditEntry) error {
	const q = `
		INSERT INTO audit_logs (user_id, action, ip_address, user_agent, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	var metaJSON []byte
	if len(entry.Metadata) > 0 {
		var err error
		metaJSON, err = json.Marshal(entry.Metadata)
		if err != nil {
			return fmt.Errorf("marshal audit metadata: %w", err)
		}
	}

	_, err := r.db.ExecContext(ctx, q,
		entry.UserID,
		string(entry.Action),
		entry.IP,
		entry.UserAgent,
		metaJSON,
		entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}
