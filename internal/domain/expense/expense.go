package expense

import (
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("expense not found")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrEmptyDescription  = errors.New("description cannot be empty")
	ErrInvalidUserID     = errors.New("user ID must be positive")
	ErrInvalidCategoryID = errors.New("category ID must be positive")
	ErrForbidden         = errors.New("access denied")
)

type Expense struct {
	ID          int64
	UserID      int64
	Description string
	Amount      float64
	CategoryID  int64
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ExpenseDetail includes denormalized category name for list views.
type ExpenseDetail struct {
	ID           int64
	UserID       int64
	Description  string
	Amount       float64
	CategoryID   int64
	CategoryName string
	Date         time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// New creates an Expense enforcing domain invariants.
func New(userID int64, description string, amount float64, categoryID int64, date time.Time) (*Expense, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if description == "" {
		return nil, ErrEmptyDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}
	now := time.Now()
	return &Expense{
		UserID:      userID,
		Description: description,
		Amount:      amount,
		CategoryID:  categoryID,
		Date:        date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}