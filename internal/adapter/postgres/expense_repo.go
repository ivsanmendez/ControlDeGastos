package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	const q = `
		INSERT INTO expenses (description, amount, category, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRowContext(ctx, q,
		e.Description,
		e.Amount,
		string(e.Category),
		e.Date,
		e.CreatedAt,
		e.UpdatedAt,
	).Scan(&e.ID)
}

func (r *ExpenseRepo) FindByID(ctx context.Context, id int64) (*expense.Expense, error) {
	const q = `
		SELECT id, description, amount, category, date, created_at, updated_at
		FROM expenses
		WHERE id = $1`

	var e expense.Expense
	var category string

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&e.ID,
		&e.Description,
		&e.Amount,
		&category,
		&e.Date,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, expense.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find expense %d: %w", id, err)
	}

	e.Category = expense.Category(category)
	return &e, nil
}

func (r *ExpenseRepo) FindAll(ctx context.Context) ([]expense.Expense, error) {
	const q = `
		SELECT id, description, amount, category, date, created_at, updated_at
		FROM expenses
		ORDER BY date DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []expense.Expense
	for rows.Next() {
		var e expense.Expense
		var category string

		if err := rows.Scan(
			&e.ID,
			&e.Description,
			&e.Amount,
			&category,
			&e.Date,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}

		e.Category = expense.Category(category)
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}

	return expenses, nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM expenses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete expense %d: %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete expense %d: %w", id, err)
	}
	if rows == 0 {
		return expense.ErrNotFound
	}

	return nil
}
