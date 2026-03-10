package contributor

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("contributor not found")
	ErrDuplicate        = errors.New("contributor with this house number already exists")
	ErrEmptyHouseNumber = errors.New("house number cannot be empty")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrInvalidUserID    = errors.New("user ID must be positive")
)

type Contributor struct {
	ID          int64
	HouseNumber string
	Name        string
	Phone       string
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository is the outbound port for contributor persistence.
type Repository interface {
	Save(ctx context.Context, c *Contributor) error
	FindByID(ctx context.Context, id int64) (*Contributor, error)
	FindAll(ctx context.Context) ([]Contributor, error)
	Update(ctx context.Context, c *Contributor) error
	Delete(ctx context.Context, id int64) error
}

// New creates a Contributor enforcing domain invariants.
func New(userID int64, houseNumber, name, phone string) (*Contributor, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if houseNumber == "" {
		return nil, ErrEmptyHouseNumber
	}
	if name == "" {
		return nil, ErrEmptyName
	}

	now := time.Now()
	return &Contributor{
		HouseNumber: houseNumber,
		Name:        name,
		Phone:       phone,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
