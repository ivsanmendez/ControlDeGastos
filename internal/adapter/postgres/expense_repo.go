package postgres

import (
	"context"
	"database/sql"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
)

// ExpenseRepo implements expense.Repository.
type ExpenseRepo struct {
	db *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Save(ctx context.Context, e *expense.Expense) error {
	// TODO: INSERT INTO expenses ... RETURNING id
	return nil
}

func (r *ExpenseRepo) FindByID(ctx context.Context, id int64) (*expense.Expense, error) {
	// TODO: SELECT ... WHERE id = $1
	return nil, expense.ErrNotFound
}

func (r *ExpenseRepo) FindAll(ctx context.Context) ([]expense.Expense, error) {
	// TODO: SELECT ... ORDER BY date DESC
	return nil, nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id int64) error {
	// TODO: DELETE FROM expenses WHERE id = $1
	return nil
}