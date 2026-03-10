package port

import (
	"context"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contributor"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// ExpenseService is the driving port — the contract that inbound adapters
// (HTTP handlers, AI agents) depend on.
type ExpenseService interface {
	CreateExpense(ctx context.Context, callerID int64, description string, amount float64, category expense.Category, date time.Time) (*expense.Expense, error)
	GetExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64) (*expense.Expense, error)
	ListExpenses(ctx context.Context, callerID int64, callerRole user.Role) ([]expense.Expense, error)
	DeleteExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64) error
}

// ContributorService is the driving port for contributor use cases.
type ContributorService interface {
	CreateContributor(ctx context.Context, callerID int64, houseNumber, name, phone string) (*contributor.Contributor, error)
	GetContributor(ctx context.Context, id int64) (*contributor.Contributor, error)
	ListContributors(ctx context.Context) ([]contributor.Contributor, error)
	UpdateContributor(ctx context.Context, id int64, name, phone string) (*contributor.Contributor, error)
	DeleteContributor(ctx context.Context, id int64) error
}

// ContributionService is the driving port for contribution use cases.
type ContributionService interface {
	CreateContribution(ctx context.Context, callerID int64, contributorID int64, amount float64, month, year int, paymentDate time.Time, paymentMethod contribution.PaymentMethod) (*contribution.Contribution, error)
	GetContribution(ctx context.Context, id int64) (*contribution.ContributionDetail, error)
	ListContributions(ctx context.Context, contributorID int64, year int) ([]contribution.ContributionDetail, error)
	DeleteContribution(ctx context.Context, id int64) error
}
