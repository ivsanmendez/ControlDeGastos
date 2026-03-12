package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	ec "github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense_category"
)

// ExpenseCategoryRepo implements expense_category.Repository.
type ExpenseCategoryRepo struct {
	db *sql.DB
}

func NewExpenseCategoryRepo(db *sql.DB) *ExpenseCategoryRepo {
	return &ExpenseCategoryRepo{db: db}
}

func (r *ExpenseCategoryRepo) Save(ctx context.Context, c *ec.ExpenseCategory) error {
	const q = `
		INSERT INTO expense_categories (name, description, is_active, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, q,
		c.Name,
		c.Description,
		c.IsActive,
		c.UserID,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ec.ErrDuplicate
		}
		return fmt.Errorf("save expense category: %w", err)
	}
	return nil
}

func (r *ExpenseCategoryRepo) FindByID(ctx context.Context, id int64) (*ec.ExpenseCategory, error) {
	const q = `
		SELECT id, name, description, is_active, user_id, created_at, updated_at
		FROM expense_categories
		WHERE id = $1`

	c, err := r.scanOne(ctx, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ec.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find expense category %d: %w", id, err)
	}
	return c, nil
}

func (r *ExpenseCategoryRepo) FindAll(ctx context.Context) ([]ec.ExpenseCategory, error) {
	const q = `
		SELECT id, name, description, is_active, user_id, created_at, updated_at
		FROM expense_categories
		ORDER BY name`

	return r.scanMany(ctx, q)
}

func (r *ExpenseCategoryRepo) FindActive(ctx context.Context) ([]ec.ExpenseCategory, error) {
	const q = `
		SELECT id, name, description, is_active, user_id, created_at, updated_at
		FROM expense_categories
		WHERE is_active = TRUE
		ORDER BY name`

	return r.scanMany(ctx, q)
}

func (r *ExpenseCategoryRepo) Update(ctx context.Context, c *ec.ExpenseCategory) error {
	const q = `
		UPDATE expense_categories
		SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5`

	result, err := r.db.ExecContext(ctx, q, c.Name, c.Description, c.IsActive, c.UpdatedAt, c.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ec.ErrDuplicate
		}
		return fmt.Errorf("update expense category %d: %w", c.ID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update expense category %d: %w", c.ID, err)
	}
	if rows == 0 {
		return ec.ErrNotFound
	}
	return nil
}

func (r *ExpenseCategoryRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM expense_categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return fmt.Errorf("cannot delete: expense category is referenced by expenses")
		}
		return fmt.Errorf("delete expense category %d: %w", id, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete expense category %d: %w", id, err)
	}
	if rows == 0 {
		return ec.ErrNotFound
	}
	return nil
}

// --- Scanners ---

func (r *ExpenseCategoryRepo) scanOne(ctx context.Context, query string, args ...any) (*ec.ExpenseCategory, error) {
	var c ec.ExpenseCategory
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.IsActive,
		&c.UserID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ExpenseCategoryRepo) scanMany(ctx context.Context, query string, args ...any) ([]ec.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list expense categories: %w", err)
	}
	defer rows.Close()

	var categories []ec.ExpenseCategory
	for rows.Next() {
		var c ec.ExpenseCategory
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.IsActive,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan expense category: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list expense categories: %w", err)
	}
	return categories, nil
}
