package port

import (
	"context"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/category"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contributor"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
	ec "github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense_category"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/receipt"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/report"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// ExpenseService is the driving port — the contract that inbound adapters
// (HTTP handlers, AI agents) depend on.
type ExpenseService interface {
	CreateExpense(ctx context.Context, callerID int64, description string, amount float64, categoryID int64, date time.Time) (*expense.Expense, error)
	GetExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64) (*expense.Expense, error)
	ListExpenses(ctx context.Context, callerID int64, callerRole user.Role) ([]expense.ExpenseDetail, error)
	UpdateExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64, description string, amount float64, categoryID int64, date time.Time) (*expense.Expense, error)
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
	CreateContribution(ctx context.Context, callerID int64, contributorID int64, categoryID int64, amount float64, month, year int, paymentDate time.Time, paymentMethod contribution.PaymentMethod) (*contribution.Contribution, error)
	GetContribution(ctx context.Context, id int64) (*contribution.ContributionDetail, error)
	ListContributions(ctx context.Context, contributorID int64, year int) ([]contribution.ContributionDetail, error)
	UpdateContribution(ctx context.Context, id int64, contributorID int64, categoryID int64, amount float64, month, year int, paymentDate time.Time, paymentMethod contribution.PaymentMethod) (*contribution.Contribution, error)
	DeleteContribution(ctx context.Context, id int64) error
}

// CategoryService is the driving port for contribution category use cases.
type CategoryService interface {
	CreateCategory(ctx context.Context, callerID int64, name, description string) (*category.Category, error)
	GetCategory(ctx context.Context, id int64) (*category.Category, error)
	ListCategories(ctx context.Context) ([]category.Category, error)
	ListActiveCategories(ctx context.Context) ([]category.Category, error)
	UpdateCategory(ctx context.Context, id int64, name, description string, isActive bool) (*category.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

// ReceiptFolioService is the driving port for receipt folio use cases.
type ReceiptFolioService interface {
	GenerateNewFolio(ctx context.Context, year int) (folio string, seq int, suffix string, err error)
	SaveFolio(ctx context.Context, rf *receipt.ReceiptFolio) error
	VerifyFolio(ctx context.Context, folio string) (*receipt.ReceiptFolio, error)
}

// ReportService is the driving port for report use cases.
type ReportService interface {
	GetMonthlyBalance(ctx context.Context, year int) (*report.MonthlyBalanceReport, error)
}

// ExpenseCategoryService is the driving port for expense category use cases.
type ExpenseCategoryService interface {
	CreateCategory(ctx context.Context, callerID int64, name, description string) (*ec.ExpenseCategory, error)
	GetCategory(ctx context.Context, id int64) (*ec.ExpenseCategory, error)
	ListCategories(ctx context.Context) ([]ec.ExpenseCategory, error)
	ListActiveCategories(ctx context.Context) ([]ec.ExpenseCategory, error)
	UpdateCategory(ctx context.Context, id int64, name, description string, isActive bool) (*ec.ExpenseCategory, error)
	DeleteCategory(ctx context.Context, id int64) error
}
