package user

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound           = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// New creates a User enforcing domain invariants.
func New(email, passwordHash string, role Role) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}
	now := time.Now()
	return &User{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ValidatePassword checks minimum password requirements (called before hashing).
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	at := strings.Index(email, "@")
	if at < 1 {
		return false
	}
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")
	return dot > 0 && dot < len(domain)-1
}
