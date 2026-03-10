package user

import (
	"context"
	"log"
	"time"
)

// Service orchestrates authentication use cases.
type Service struct {
	repo   Repository
	hasher PasswordHasher
	tokens TokenIssuer
	audit  AuditLogger
}

func NewService(repo Repository, hasher PasswordHasher, tokens TokenIssuer, audit AuditLogger) *Service {
	return &Service{repo: repo, hasher: hasher, tokens: tokens, audit: audit}
}

func (s *Service) Register(ctx context.Context, email, password string, info AuditInfo) (*User, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	u, err := New(email, hash, RoleUser)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}

	s.logAudit(ctx, &u.ID, AuditRegister, info, map[string]string{"email": u.Email})
	return u, nil
}

func (s *Service) Login(ctx context.Context, email, password string, info AuditInfo) (*User, TokenPair, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logAudit(ctx, nil, AuditLoginFailed, info, map[string]string{"email": email})
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	if err := s.hasher.Compare(u.PasswordHash, password); err != nil {
		s.logAudit(ctx, &u.ID, AuditLoginFailed, info, map[string]string{"email": email})
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	pair, err := s.tokens.Issue(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, TokenPair{}, err
	}

	rt := &RefreshToken{
		UserID:    u.ID,
		TokenHash: HashToken(pair.RefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, TokenPair{}, err
	}

	s.logAudit(ctx, &u.ID, AuditLoginSuccess, info, nil)
	return u, pair, nil
}

func (s *Service) RefreshToken(ctx context.Context, rawRefresh string, info AuditInfo) (TokenPair, error) {
	hash := HashToken(rawRefresh)
	stored, err := s.repo.FindRefreshTokenByHash(ctx, hash)
	if err != nil {
		return TokenPair{}, ErrTokenNotFound
	}

	// Reuse detection: if this token was already revoked, someone is replaying it.
	if stored.Revoked {
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, stored.UserID)
		s.logAudit(ctx, &stored.UserID, AuditTokenRefresh, info, map[string]string{"reuse_detected": "true"})
		return TokenPair{}, ErrTokenRevoked
	}

	if stored.IsExpired() {
		return TokenPair{}, ErrTokenExpired
	}

	// Revoke the old token.
	if err := s.repo.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return TokenPair{}, err
	}

	// Look up user to get current email/role for the new JWT.
	u, err := s.repo.FindByID(ctx, stored.UserID)
	if err != nil {
		return TokenPair{}, err
	}

	pair, err := s.tokens.Issue(u.ID, u.Email, u.Role)
	if err != nil {
		return TokenPair{}, err
	}

	rt := &RefreshToken{
		UserID:    u.ID,
		TokenHash: HashToken(pair.RefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return TokenPair{}, err
	}

	s.logAudit(ctx, &u.ID, AuditTokenRefresh, info, nil)
	return pair, nil
}

func (s *Service) Logout(ctx context.Context, rawRefresh string, info AuditInfo) error {
	hash := HashToken(rawRefresh)
	stored, err := s.repo.FindRefreshTokenByHash(ctx, hash)
	if err != nil {
		return ErrTokenNotFound
	}
	if err := s.repo.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return err
	}
	s.logAudit(ctx, &stored.UserID, AuditLogout, info, nil)
	return nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// logAudit fires-and-forgets an audit entry. Errors are logged but never returned.
func (s *Service) logAudit(ctx context.Context, userID *int64, action AuditAction, info AuditInfo, metadata map[string]string) {
	entry := NewAuditEntry(userID, action, info, metadata)
	if err := s.audit.Log(ctx, entry); err != nil {
		log.Printf("audit log error: %v", err)
	}
}
