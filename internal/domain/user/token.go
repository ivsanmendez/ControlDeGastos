package user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

var (
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrTokenNotFound = errors.New("refresh token not found")
)

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

// IsExpired reports whether the token has passed its expiration time.
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsable reports whether the token is valid for use.
func (t *RefreshToken) IsUsable() bool {
	return !t.Revoked && !t.IsExpired()
}

// HashToken returns the SHA-256 hex digest of a raw token string.
func HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// TokenPair holds an access token and a raw refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
