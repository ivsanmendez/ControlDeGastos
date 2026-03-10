package user

import "context"

// Repository is the outbound port for user and refresh-token persistence.
type Repository interface {
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)

	SaveRefreshToken(ctx context.Context, t *RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id int64) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error
}

// PasswordHasher is the outbound port for password hashing.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// TokenIssuer is the outbound port for JWT issuance.
type TokenIssuer interface {
	Issue(userID int64, email string, role Role) (TokenPair, error)
}

// AuditLogger is the outbound port for persisting audit entries.
type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}
