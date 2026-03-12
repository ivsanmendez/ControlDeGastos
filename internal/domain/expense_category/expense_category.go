package expense_category

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("expense category not found")
	ErrDuplicate     = errors.New("expense category name already exists")
	ErrEmptyName     = errors.New("expense category name must not be empty")
	ErrInvalidUserID = errors.New("user ID must be positive")
)

// ExpenseCategory represents a user-defined expense category.
type ExpenseCategory struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates an ExpenseCategory enforcing domain invariants.
func New(userID int64, name, description string) (*ExpenseCategory, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if name == "" {
		return nil, ErrEmptyName
	}

	now := time.Now()
	return &ExpenseCategory{
		Name:        name,
		Description: description,
		IsActive:    true,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
