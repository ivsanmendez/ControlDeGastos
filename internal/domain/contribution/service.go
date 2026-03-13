package contribution

import (
	"context"
	"time"
)

// Repository is the outbound port for contribution persistence.
type Repository interface {
	Save(ctx context.Context, c *Contribution) error
	Update(ctx context.Context, c *Contribution) error
	FindByID(ctx context.Context, id int64) (*Contribution, error)
	FindAll(ctx context.Context) ([]Contribution, error)
	FindByContributorAndYear(ctx context.Context, contributorID int64, year int) ([]Contribution, error)
	Delete(ctx context.Context, id int64) error

	// Detailed variants return ContributionDetail with contributor info via JOIN.
	FindDetailedByID(ctx context.Context, id int64) (*ContributionDetail, error)
	FindAllDetailed(ctx context.Context) ([]ContributionDetail, error)
	FindDetailedByContributorAndYear(ctx context.Context, contributorID int64, year int) ([]ContributionDetail, error)
}

// Service orchestrates contribution use cases.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateContribution(
	ctx context.Context,
	callerID int64,
	contributorID int64,
	categoryID int64,
	amount float64,
	month int,
	year int,
	paymentDate time.Time,
	paymentMethod PaymentMethod,
) (*Contribution, error) {
	c, err := New(callerID, contributorID, categoryID, amount, month, year, paymentDate, paymentMethod)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetContribution(ctx context.Context, id int64) (*ContributionDetail, error) {
	return s.repo.FindDetailedByID(ctx, id)
}

func (s *Service) ListContributions(ctx context.Context, contributorID int64, year int) ([]ContributionDetail, error) {
	if contributorID > 0 && year > 0 {
		return s.repo.FindDetailedByContributorAndYear(ctx, contributorID, year)
	}
	return s.repo.FindAllDetailed(ctx)
}

func (s *Service) UpdateContribution(
	ctx context.Context,
	id int64,
	contributorID int64,
	categoryID int64,
	amount float64,
	month int,
	year int,
	paymentDate time.Time,
	paymentMethod PaymentMethod,
) (*Contribution, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
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

	existing.ContributorID = contributorID
	existing.CategoryID = categoryID
	existing.Amount = amount
	existing.Month = month
	existing.Year = year
	existing.PaymentDate = paymentDate
	existing.PaymentMethod = paymentMethod
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteContribution(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
