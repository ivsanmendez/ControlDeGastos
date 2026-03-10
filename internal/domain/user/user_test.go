package user_test

import (
	"testing"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

func TestNew_Valid(t *testing.T) {
	u, err := user.New("test@example.com", "hashed", user.RoleUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Email != "test@example.com" {
		t.Errorf("email = %q, want %q", u.Email, "test@example.com")
	}
	if u.Role != user.RoleUser {
		t.Errorf("role = %q, want %q", u.Role, user.RoleUser)
	}
	if u.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestNew_NormalizesEmail(t *testing.T) {
	u, err := user.New("  TEST@Example.COM  ", "hashed", user.RoleUser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "test@example.com" {
		t.Errorf("email = %q, want lowercased/trimmed", u.Email)
	}
}

func TestNew_InvalidEmail(t *testing.T) {
	cases := []string{"", "noat", "@nodomain", "user@", "user@.com", "user@domain."}
	for _, email := range cases {
		_, err := user.New(email, "hashed", user.RoleUser)
		if err != user.ErrInvalidEmail {
			t.Errorf("New(%q): expected ErrInvalidEmail, got %v", email, err)
		}
	}
}

func TestNew_EmptyHash(t *testing.T) {
	_, err := user.New("test@example.com", "", user.RoleUser)
	if err == nil {
		t.Error("expected error for empty hash")
	}
}

func TestValidatePassword(t *testing.T) {
	if err := user.ValidatePassword("short"); err != user.ErrWeakPassword {
		t.Errorf("expected ErrWeakPassword, got %v", err)
	}
	if err := user.ValidatePassword("longenough"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHashToken(t *testing.T) {
	h1 := user.HashToken("token-abc")
	h2 := user.HashToken("token-abc")
	if h1 != h2 {
		t.Error("same input should produce same hash")
	}
	h3 := user.HashToken("token-xyz")
	if h1 == h3 {
		t.Error("different input should produce different hash")
	}
	if len(h1) != 64 {
		t.Errorf("expected 64-char hex, got %d chars", len(h1))
	}
}

func TestRoleHasPermission(t *testing.T) {
	if !user.RoleHasPermission(user.RoleUser, user.PermExpenseCreate) {
		t.Error("user should have expense:create")
	}
	if user.RoleHasPermission(user.RoleUser, user.PermExpenseReadAll) {
		t.Error("user should NOT have expense:read:all")
	}
	if !user.RoleHasPermission(user.RoleAdmin, user.PermExpenseReadAll) {
		t.Error("admin should have expense:read:all")
	}
	if !user.RoleHasPermission(user.RoleAdmin, user.PermExpenseReadOwn) {
		t.Error("admin should also have expense:read:own")
	}
	if user.RoleHasPermission("unknown", user.PermExpenseCreate) {
		t.Error("unknown role should have no permissions")
	}
}

func TestRefreshToken_IsUsable(t *testing.T) {
	rt := &user.RefreshToken{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if !rt.IsUsable() {
		t.Error("expected usable")
	}

	rt.Revoked = true
	if rt.IsUsable() {
		t.Error("revoked token should not be usable")
	}
}

func TestRefreshToken_IsExpired(t *testing.T) {
	rt := &user.RefreshToken{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	if !rt.IsExpired() {
		t.Error("expected expired")
	}
	if rt.IsUsable() {
		t.Error("expired token should not be usable")
	}
}
