package port

import (
	"context"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/category"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contributor"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
	ec "github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense_category"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/receipt"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/report"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// ExpenseRepository is the driven port for expense persistence.
type ExpenseRepository = expense.Repository

// EventPublisher is the driven port for publishing domain events.
type EventPublisher = expense.EventPublisher

// EventSubscriber reacts to domain events.
type EventSubscriber interface {
	Handle(ctx context.Context, event expense.Event) error
}

// UserRepository is the driven port for user persistence.
type UserRepository = user.Repository

// PasswordHasher is the driven port for password hashing.
type PasswordHasher = user.PasswordHasher

// TokenIssuer is the driven port for JWT issuance.
type TokenIssuer = user.TokenIssuer

// AuditLogger is the driven port for audit logging.
type AuditLogger = user.AuditLogger

// ContributorRepository is the driven port for contributor persistence.
type ContributorRepository = contributor.Repository

// ContributionRepository is the driven port for contribution persistence.
type ContributionRepository = contribution.Repository

// CategoryRepository is the driven port for contribution category persistence.
type CategoryRepository = category.Repository

// ReceiptFolioRepository is the driven port for receipt folio persistence.
type ReceiptFolioRepository = receipt.Repository

// ReportRepository is the driven port for report aggregation queries.
type ReportRepository = report.Repository

// ExpenseCategoryRepository is the driven port for expense category persistence.
type ExpenseCategoryRepository = ec.Repository

// ReceiptSigner is the driven port for digitally signing receipt data.
// The password is required per-call to decrypt the private key (SAT format).
type ReceiptSigner interface {
	Sign(data []byte, password string) ([]byte, error)
	Certificate() []byte
	Available() bool
}
