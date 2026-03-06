package port

import (
	"context"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
)

// ExpenseService is the driving port — the contract that inbound adapters
// (HTTP handlers, AI agents) depend on.
type ExpenseService interface {
	CreateExpense(ctx context.Context, description string, amount float64, category expense.Category, date time.Time) (*expense.Expense, error)
	GetExpense(ctx context.Context, id int64) (*expense.Expense, error)
	ListExpenses(ctx context.Context) ([]expense.Expense, error)
	DeleteExpense(ctx context.Context, id int64) error
}