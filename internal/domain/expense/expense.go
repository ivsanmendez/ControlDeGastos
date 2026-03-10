package expense

import (
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("expense not found")
	ErrInvalidAmount    = errors.New("amount must be positive")
	ErrEmptyDescription = errors.New("description cannot be empty")
	ErrInvalidUserID    = errors.New("user ID must be positive")
	ErrForbidden        = errors.New("access denied")
)

type Category string

const (
	CategoryFood      Category = "food"
	CategoryTransport Category = "transport"
	CategoryHousing   Category = "housing"
	CategoryOther     Category = "other"
)

type Expense struct {
	ID          int64
	UserID      int64
	Description string
	Amount      float64
	Category    Category
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates an Expense enforcing domain invariants.
func New(userID int64, description string, amount float64, category Category, date time.Time) (*Expense, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if description == "" {
		return nil, ErrEmptyDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	now := time.Now()
	return &Expense{
		UserID:      userID,
		Description: description,
		Amount:      amount,
		Category:    category,
		Date:        date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}