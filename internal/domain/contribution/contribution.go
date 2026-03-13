package contribution

import (
	"errors"
	"time"
)

var (
	ErrNotFound             = errors.New("contribution not found")
	ErrDuplicate            = errors.New("contribution already exists for this contributor/category/month/year")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrInvalidContributorID = errors.New("contributor ID must be positive")
	ErrInvalidCategoryID    = errors.New("category ID must be positive")
	ErrInvalidMonth         = errors.New("month must be between 1 and 12")
	ErrInvalidYear          = errors.New("year must be >= 2000")
	ErrInvalidPaymentMethod = errors.New("payment method must be cash, transfer, or other")
	ErrInvalidUserID        = errors.New("user ID must be positive")
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentTransfer PaymentMethod = "transfer"
	PaymentOther    PaymentMethod = "other"
)

func (p PaymentMethod) Valid() bool {
	switch p {
	case PaymentCash, PaymentTransfer, PaymentOther:
		return true
	}
	return false
}

type Contribution struct {
	ID            int64
	ContributorID int64
	CategoryID    int64
	Amount        float64
	Month         int
	Year          int
	PaymentDate   time.Time
	PaymentMethod PaymentMethod
	UserID        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ContributionDetail is a read-only DTO returned by JOIN queries,
// enriching a Contribution with contributor and category info.
type ContributionDetail struct {
	Contribution
	HouseNumber     string
	ContributorName string
	Phone           string
	CategoryName    string
}

// New creates a Contribution enforcing domain invariants.
func New(
	userID int64,
	contributorID int64,
	categoryID int64,
	amount float64,
	month int,
	year int,
	paymentDate time.Time,
	paymentMethod PaymentMethod,
) (*Contribution, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if contributorID <= 0 {
		return nil, ErrInvalidContributorID
	}
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}
	if year < 2000 {
		return nil, ErrInvalidYear
	}
	if !paymentMethod.Valid() {
		return nil, ErrInvalidPaymentMethod
	}

	now := time.Now()
	return &Contribution{
		UserID:        userID,
		ContributorID: contributorID,
		CategoryID:    categoryID,
		Amount:        amount,
		Month:         month,
		Year:          year,
		PaymentDate:   paymentDate,
		PaymentMethod: paymentMethod,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
