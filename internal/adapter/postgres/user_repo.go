package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// UserRepo implements user.Repository.
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, u *user.User) error {
	const q = `
		INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, q,
		u.Email,
		u.PasswordHash,
		string(u.Role),
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return user.ErrEmailTaken
		}
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	const q = `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1`

	var u user.User
	var role string

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user %d: %w", id, err)
	}
	u.Role = user.Role(role)
	return &u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	const q = `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1`

	var u user.User
	var role string

	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	u.Role = user.Role(role)
	return &u, nil
}

func (r *UserRepo) SaveRefreshToken(ctx context.Context, t *user.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRowContext(ctx, q,
		t.UserID, t.TokenHash, t.ExpiresAt, t.Revoked, t.CreatedAt,
	).Scan(&t.ID)
}

func (r *UserRepo) FindRefreshTokenByHash(ctx context.Context, hash string) (*user.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at
		FROM refresh_tokens WHERE token_hash = $1`

	var t user.RefreshToken
	err := r.db.QueryRowContext(ctx, q, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.Revoked, &t.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return &t, nil
}

func (r *UserRepo) RevokeRefreshToken(ctx context.Context, id int64) error {
	const q = `UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1`
	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token %d: %w", id, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke refresh token %d: %w", id, err)
	}
	if rows == 0 {
		return user.ErrTokenNotFound
	}
	return nil
}

func (r *UserRepo) RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error {
	const q = `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1 AND revoked = FALSE`
	_, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("revoke all tokens for user %d: %w", userID, err)
	}
	return nil
}
