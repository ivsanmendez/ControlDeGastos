package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contributor"
)

// ContributorRepo implements contributor.Repository.
type ContributorRepo struct {
	db *sql.DB
}

func NewContributorRepo(db *sql.DB) *ContributorRepo {
	return &ContributorRepo{db: db}
}

func (r *ContributorRepo) Save(ctx context.Context, c *contributor.Contributor) error {
	const q = `
		INSERT INTO contributors (house_number, name, phone, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, q,
		c.HouseNumber,
		c.Name,
		c.Phone,
		c.UserID,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&c.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return contributor.ErrDuplicate
		}
		return fmt.Errorf("save contributor: %w", err)
	}
	return nil
}

func (r *ContributorRepo) FindByID(ctx context.Context, id int64) (*contributor.Contributor, error) {
	const q = `
		SELECT id, house_number, name, phone, user_id, created_at, updated_at
		FROM contributors
		WHERE id = $1`

	c, err := r.scanOne(ctx, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, contributor.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find contributor %d: %w", id, err)
	}
	return c, nil
}

func (r *ContributorRepo) FindAll(ctx context.Context) ([]contributor.Contributor, error) {
	const q = `
		SELECT id, house_number, name, phone, user_id, created_at, updated_at
		FROM contributors
		ORDER BY house_number`

	return r.scanMany(ctx, q)
}

func (r *ContributorRepo) Update(ctx context.Context, c *contributor.Contributor) error {
	const q = `
		UPDATE contributors SET name = $1, phone = $2, updated_at = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, q, c.Name, c.Phone, c.UpdatedAt, c.ID)
	if err != nil {
		return fmt.Errorf("update contributor %d: %w", c.ID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update contributor %d: %w", c.ID, err)
	}
	if rows == 0 {
		return contributor.ErrNotFound
	}
	return nil
}

func (r *ContributorRepo) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM contributors WHERE id = $1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return fmt.Errorf("cannot delete contributor: still referenced by contributions")
		}
		return fmt.Errorf("delete contributor %d: %w", id, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete contributor %d: %w", id, err)
	}
	if rows == 0 {
		return contributor.ErrNotFound
	}
	return nil
}

func (r *ContributorRepo) scanOne(ctx context.Context, query string, args ...any) (*contributor.Contributor, error) {
	var c contributor.Contributor
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&c.ID,
		&c.HouseNumber,
		&c.Name,
		&c.Phone,
		&c.UserID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContributorRepo) scanMany(ctx context.Context, query string, args ...any) ([]contributor.Contributor, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list contributors: %w", err)
	}
	defer rows.Close()

	var contributors []contributor.Contributor
	for rows.Next() {
		var c contributor.Contributor
		if err := rows.Scan(
			&c.ID,
			&c.HouseNumber,
			&c.Name,
			&c.Phone,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan contributor: %w", err)
		}
		contributors = append(contributors, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list contributors: %w", err)
	}
	return contributors, nil
}
