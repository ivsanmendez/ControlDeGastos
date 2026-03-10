package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// --- Fakes ---

type fakeRepo struct {
	users         map[int64]*user.User
	usersByEmail  map[string]*user.User
	tokens        map[int64]*user.RefreshToken
	tokensByHash  map[string]*user.RefreshToken
	nextUserID    int64
	nextTokenID   int64
	saveUserErr   error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		users:        make(map[int64]*user.User),
		usersByEmail: make(map[string]*user.User),
		tokens:       make(map[int64]*user.RefreshToken),
		tokensByHash: make(map[string]*user.RefreshToken),
		nextUserID:   1,
		nextTokenID:  1,
	}
}

func (r *fakeRepo) Save(_ context.Context, u *user.User) error {
	if r.saveUserErr != nil {
		return r.saveUserErr
	}
	if _, exists := r.usersByEmail[u.Email]; exists {
		return user.ErrEmailTaken
	}
	u.ID = r.nextUserID
	r.nextUserID++
	cp := *u
	r.users[u.ID] = &cp
	r.usersByEmail[u.Email] = &cp
	return nil
}

func (r *fakeRepo) FindByID(_ context.Context, id int64) (*user.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *fakeRepo) FindByEmail(_ context.Context, email string) (*user.User, error) {
	u, ok := r.usersByEmail[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *fakeRepo) SaveRefreshToken(_ context.Context, t *user.RefreshToken) error {
	t.ID = r.nextTokenID
	r.nextTokenID++
	cp := *t
	r.tokens[t.ID] = &cp
	r.tokensByHash[t.TokenHash] = &cp
	return nil
}

func (r *fakeRepo) FindRefreshTokenByHash(_ context.Context, hash string) (*user.RefreshToken, error) {
	t, ok := r.tokensByHash[hash]
	if !ok {
		return nil, user.ErrTokenNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *fakeRepo) RevokeRefreshToken(_ context.Context, id int64) error {
	t, ok := r.tokens[id]
	if !ok {
		return user.ErrTokenNotFound
	}
	t.Revoked = true
	// Also update the hash-indexed map.
	if byHash, exists := r.tokensByHash[t.TokenHash]; exists {
		byHash.Revoked = true
	}
	return nil
}

func (r *fakeRepo) RevokeAllUserRefreshTokens(_ context.Context, userID int64) error {
	for _, t := range r.tokens {
		if t.UserID == userID {
			t.Revoked = true
			if byHash, exists := r.tokensByHash[t.TokenHash]; exists {
				byHash.Revoked = true
			}
		}
	}
	return nil
}

// fakeHasher uses a simple prefix for deterministic testing.
type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (fakeHasher) Compare(hash, password string) error {
	if hash == "hashed:"+password {
		return nil
	}
	return errors.New("mismatch")
}

// fakeIssuer returns predictable tokens.
type fakeIssuer struct {
	callCount int
}

func (f *fakeIssuer) Issue(userID int64, email string, role user.Role) (user.TokenPair, error) {
	f.callCount++
	return user.TokenPair{
		AccessToken:  fmt.Sprintf("jwt-%d-%d", userID, f.callCount),
		RefreshToken: fmt.Sprintf("refresh-%d-%d", userID, f.callCount),
	}, nil
}

// fakeAudit captures audit entries.
type fakeAudit struct {
	entries []user.AuditEntry
}

func (a *fakeAudit) Log(_ context.Context, entry user.AuditEntry) error {
	a.entries = append(a.entries, entry)
	return nil
}

// --- Helpers ---

var testCtx = context.Background()

func newTestService() (*user.Service, *fakeRepo, *fakeIssuer, *fakeAudit) {
	repo := newFakeRepo()
	issuer := &fakeIssuer{}
	audit := &fakeAudit{}
	svc := user.NewService(repo, fakeHasher{}, issuer, audit)
	return svc, repo, issuer, audit
}

var testAuditInfo = user.AuditInfo{IP: "127.0.0.1", UserAgent: "test-agent"}

// --- Tests ---

func TestRegister_HappyPath(t *testing.T) {
	svc, repo, _, audit := newTestService()

	u, err := svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == 0 {
		t.Error("expected ID to be set")
	}
	if u.Email != "user@example.com" {
		t.Errorf("email = %q, want %q", u.Email, "user@example.com")
	}
	if u.Role != user.RoleUser {
		t.Errorf("role = %q, want %q", u.Role, user.RoleUser)
	}
	if _, ok := repo.users[u.ID]; !ok {
		t.Error("user not found in repo")
	}
	if len(audit.entries) != 1 || audit.entries[0].Action != user.AuditRegister {
		t.Errorf("expected register audit entry, got %v", audit.entries)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, err := svc.Register(testCtx, "user@example.com", "short", testAuditInfo)
	if !errors.Is(err, user.ErrWeakPassword) {
		t.Errorf("expected ErrWeakPassword, got %v", err)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, err := svc.Register(testCtx, "bademail", "password123", testAuditInfo)
	if !errors.Is(err, user.ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, err := svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = svc.Register(testCtx, "user@example.com", "password456", testAuditInfo)
	if !errors.Is(err, user.ErrEmailTaken) {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
}

func TestLogin_HappyPath(t *testing.T) {
	svc, _, _, audit := newTestService()

	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	audit.entries = nil // reset after register

	u, pair, err := svc.Login(testCtx, "user@example.com", "password123", testAuditInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "user@example.com" {
		t.Errorf("email = %q", u.Email)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Error("expected non-empty tokens")
	}
	if len(audit.entries) != 1 || audit.entries[0].Action != user.AuditLoginSuccess {
		t.Errorf("expected login_success audit, got %v", audit.entries)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _, audit := newTestService()
	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	audit.entries = nil

	_, _, err := svc.Login(testCtx, "user@example.com", "wrong", testAuditInfo)
	if !errors.Is(err, user.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
	if len(audit.entries) != 1 || audit.entries[0].Action != user.AuditLoginFailed {
		t.Errorf("expected login_failed audit, got %v", audit.entries)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, _, err := svc.Login(testCtx, "nobody@example.com", "password123", testAuditInfo)
	if !errors.Is(err, user.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRefreshToken_HappyPath(t *testing.T) {
	svc, _, _, audit := newTestService()
	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	_, pair, _ := svc.Login(testCtx, "user@example.com", "password123", testAuditInfo)
	audit.entries = nil

	newPair, err := svc.RefreshToken(testCtx, pair.RefreshToken, testAuditInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newPair.AccessToken == "" || newPair.RefreshToken == "" {
		t.Error("expected non-empty tokens")
	}
	if newPair.AccessToken == pair.AccessToken {
		t.Error("new access token should differ from old")
	}
	if len(audit.entries) != 1 || audit.entries[0].Action != user.AuditTokenRefresh {
		t.Errorf("expected token_refresh audit, got %v", audit.entries)
	}
}

func TestRefreshToken_ReuseDetection(t *testing.T) {
	svc, _, _, audit := newTestService()
	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	_, pair, _ := svc.Login(testCtx, "user@example.com", "password123", testAuditInfo)

	// First refresh succeeds.
	_, err := svc.RefreshToken(testCtx, pair.RefreshToken, testAuditInfo)
	if err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}
	audit.entries = nil

	// Second use of same token → reuse detection.
	_, err = svc.RefreshToken(testCtx, pair.RefreshToken, testAuditInfo)
	if !errors.Is(err, user.ErrTokenRevoked) {
		t.Errorf("expected ErrTokenRevoked, got %v", err)
	}
	// Should have audit entry with reuse_detected.
	if len(audit.entries) != 1 || audit.entries[0].Metadata["reuse_detected"] != "true" {
		t.Errorf("expected reuse_detected audit, got %v", audit.entries)
	}
}

func TestRefreshToken_Expired(t *testing.T) {
	svc, repo, _, _ := newTestService()
	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	_, pair, _ := svc.Login(testCtx, "user@example.com", "password123", testAuditInfo)

	// Manually expire the token.
	hash := user.HashToken(pair.RefreshToken)
	if stored, ok := repo.tokensByHash[hash]; ok {
		stored.ExpiresAt = time.Now().Add(-time.Hour)
		repo.tokens[stored.ID].ExpiresAt = stored.ExpiresAt
	}

	_, err := svc.RefreshToken(testCtx, pair.RefreshToken, testAuditInfo)
	if !errors.Is(err, user.ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestRefreshToken_NotFound(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, err := svc.RefreshToken(testCtx, "nonexistent", testAuditInfo)
	if !errors.Is(err, user.ErrTokenNotFound) {
		t.Errorf("expected ErrTokenNotFound, got %v", err)
	}
}

func TestLogout_HappyPath(t *testing.T) {
	svc, _, _, audit := newTestService()
	svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)
	_, pair, _ := svc.Login(testCtx, "user@example.com", "password123", testAuditInfo)
	audit.entries = nil

	err := svc.Logout(testCtx, pair.RefreshToken, testAuditInfo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(audit.entries) != 1 || audit.entries[0].Action != user.AuditLogout {
		t.Errorf("expected logout audit, got %v", audit.entries)
	}

	// Using same token after logout should fail on refresh.
	_, err = svc.RefreshToken(testCtx, pair.RefreshToken, testAuditInfo)
	if !errors.Is(err, user.ErrTokenRevoked) {
		t.Errorf("expected ErrTokenRevoked after logout, got %v", err)
	}
}

func TestLogout_InvalidToken(t *testing.T) {
	svc, _, _, _ := newTestService()

	err := svc.Logout(testCtx, "nonexistent", testAuditInfo)
	if !errors.Is(err, user.ErrTokenNotFound) {
		t.Errorf("expected ErrTokenNotFound, got %v", err)
	}
}

func TestGetUser(t *testing.T) {
	svc, _, _, _ := newTestService()
	registered, _ := svc.Register(testCtx, "user@example.com", "password123", testAuditInfo)

	u, err := svc.GetUser(testCtx, registered.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "user@example.com" {
		t.Errorf("email = %q", u.Email)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc, _, _, _ := newTestService()

	_, err := svc.GetUser(testCtx, 999)
	if !errors.Is(err, user.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
