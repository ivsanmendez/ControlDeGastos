package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
)

// ContributionRepo implements contribution.Repository.
type ContributionRepo struct {
	db *sql.DB
}

func NewContributionRepo(db *sql.DB) *ContributionRepo {
	return &ContributionRepo{db: db}
}

func (r *ContributionRepo) Save(ctx context.Context, c *contribution.Contribution) error {
	const q = `
		INSERT INTO contributions (contributor_id, category_id, amount, month, year, payment_date, payment_method, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, q,
		c.ContributorID,
		c.CategoryID,
		c.Amount,
		c.Month,
		c.Year,
		c.PaymentDate,
		string(c.PaymentMethod),
		c.UserID,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return contribution.ErrDuplicate
		}
		return fmt.Errorf("save contribution: %w", err)
	}
	return nil
}

func (r *ContributionRepo) Update(ctx context.Context, c *contribution.Contribution) error {
	const q = `
		UPDATE contributions
		SET contributor_id = $1, category_id = $2, amount = $3, month = $4, year = $5,
		    payment_date = $6, payment_method = $7, updated_at = $8
		WHERE id = $9`

	result, err := r.db.ExecContext(ctx, q,
		c.ContributorID,
		c.CategoryID,
		c.Amount,
		c.Month,
		c.Year,
		c.PaymentDate,
		string(c.PaymentMethod),
		c.UpdatedAt,
		c.ID,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return contribution.ErrDuplicate
		}
		return fmt.Errorf("update contribution %d: %w", c.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update contribution %d: %w", c.ID, err)
	}
	if rows == 0 {
		return contribution.ErrNotFound
	}
	return nil
}

func (r *ContributionRepo) FindByID(ctx context.Context, id int64) (*contribution.Contribution, error) {
	const q = `
		SELECT id, contributor_id, category_id, amount, month, year, payment_date, payment_method, user_id, created_at, updated_at
		FROM contributions
		WHERE id = $1`

	c, err := r.scanOne(ctx, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, contribution.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find contribution %d: %w", id, err)
	}
	return c, nil
}

func (r *ContributionRepo) FindAll(ctx context.Context) ([]contribution.Contribution, error) {
	const q = `
		SELECT id, contributor_id, category_id, amount, month, year, payment_date, payment_method, user_id, created_at, updated_at
		FROM contributions
		ORDER BY year DESC, month DESC`

	return r.scanMany(ctx, q)
}

func (r *ContributionRepo) FindByContributorAndYear(ctx context.Context, contributorID int64, year int) ([]contribution.Contribution, error) {
	const q = `
		SELECT id, contributor_id, category_id, amount, month, year, payment_date, payment_method, user_id, created_at, updated_at
		FROM contributions
		WHERE contributor_id = $1 AND year = $2
		ORDER BY month`

	return r.scanMany(ctx, q, contributorID, year)
}

func (r *ContributionRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM contributions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete contribution %d: %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete contribution %d: %w", id, err)
	}
	if rows == 0 {
		return contribution.ErrNotFound
	}
	return nil
}

// --- Detailed (JOIN) queries ---

const detailSelect = `
	SELECT c.id, c.contributor_id, c.category_id, c.amount, c.month, c.year, c.payment_date, c.payment_method, c.user_id, c.created_at, c.updated_at,
	       ct.house_number, ct.name, ct.phone,
	       cc.name
	FROM contributions c
	JOIN contributors ct ON ct.id = c.contributor_id
	JOIN contribution_categories cc ON cc.id = c.category_id`

func (r *ContributionRepo) FindDetailedByID(ctx context.Context, id int64) (*contribution.ContributionDetail, error) {
	q := detailSelect + ` WHERE c.id = $1`

	d, err := r.scanDetail(ctx, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, contribution.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find contribution detail %d: %w", id, err)
	}
	return d, nil
}

func (r *ContributionRepo) FindAllDetailed(ctx context.Context) ([]contribution.ContributionDetail, error) {
	q := detailSelect + ` ORDER BY c.year DESC, c.month DESC, ct.house_number`
	return r.scanDetails(ctx, q)
}

func (r *ContributionRepo) FindDetailedByContributorAndYear(ctx context.Context, contributorID int64, year int) ([]contribution.ContributionDetail, error) {
	q := detailSelect + ` WHERE c.contributor_id = $1 AND c.year = $2 ORDER BY c.month`
	return r.scanDetails(ctx, q, contributorID, year)
}

// --- Scanners ---

func (r *ContributionRepo) scanOne(ctx context.Context, query string, args ...any) (*contribution.Contribution, error) {
	var c contribution.Contribution
	var method string

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.ID,
		&c.ContributorID,
		&c.CategoryID,
		&c.Amount,
		&c.Month,
		&c.Year,
		&c.PaymentDate,
		&method,
		&c.UserID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	c.PaymentMethod = contribution.PaymentMethod(method)
	return &c, nil
}

func (r *ContributionRepo) scanMany(ctx context.Context, query string, args ...any) ([]contribution.Contribution, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list contributions: %w", err)
	}
	defer rows.Close()

	var contributions []contribution.Contribution
	for rows.Next() {
		var c contribution.Contribution
		var method string

		if err := rows.Scan(
			&c.ID,
			&c.ContributorID,
			&c.CategoryID,
			&c.Amount,
			&c.Month,
			&c.Year,
			&c.PaymentDate,
			&method,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan contribution: %w", err)
		}
		c.PaymentMethod = contribution.PaymentMethod(method)
		contributions = append(contributions, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list contributions: %w", err)
	}
	return contributions, nil
}

func (r *ContributionRepo) scanDetail(ctx context.Context, query string, args ...any) (*contribution.ContributionDetail, error) {
	var d contribution.ContributionDetail
	var method string

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&d.ID,
		&d.ContributorID,
		&d.CategoryID,
		&d.Amount,
		&d.Month,
		&d.Year,
		&d.PaymentDate,
		&method,
		&d.UserID,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.HouseNumber,
		&d.ContributorName,
		&d.Phone,
		&d.CategoryName,
	)
	if err != nil {
		return nil, err
	}
	d.PaymentMethod = contribution.PaymentMethod(method)
	return &d, nil
}

func (r *ContributionRepo) scanDetails(ctx context.Context, query string, args ...any) ([]contribution.ContributionDetail, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list contribution details: %w", err)
	}
	defer rows.Close()

	var details []contribution.ContributionDetail
	for rows.Next() {
		var d contribution.ContributionDetail
		var method string

		if err := rows.Scan(
			&d.ID,
			&d.ContributorID,
			&d.CategoryID,
			&d.Amount,
			&d.Month,
			&d.Year,
			&d.PaymentDate,
			&method,
			&d.UserID,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.HouseNumber,
			&d.ContributorName,
			&d.Phone,
			&d.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("scan contribution detail: %w", err)
		}
		d.PaymentMethod = contribution.PaymentMethod(method)
		details = append(details, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list contribution details: %w", err)
	}
	return details, nil
}
