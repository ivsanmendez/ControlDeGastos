package expense_category

import (
	"context"
	"time"
)

// Repository is the outbound port for expense category persistence.
type Repository interface {
	Save(ctx context.Context, c *ExpenseCategory) error
	FindByID(ctx context.Context, id int64) (*ExpenseCategory, error)
	FindAll(ctx context.Context) ([]ExpenseCategory, error)
	FindActive(ctx context.Context) ([]ExpenseCategory, error)
	Update(ctx context.Context, c *ExpenseCategory) error
	Delete(ctx context.Context, id int64) error
}

// Service orchestrates expense category use cases.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, callerID int64, name, description string) (*ExpenseCategory, error) {
	c, err := New(callerID, name, description)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*ExpenseCategory, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListCategories(ctx context.Context) ([]ExpenseCategory, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) ListActiveCategories(ctx context.Context) ([]ExpenseCategory, error) {
	return s.repo.FindActive(ctx)
}

func (s *Service) UpdateCategory(ctx context.Context, id int64, name, description string, isActive bool) (*ExpenseCategory, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	c.Name = name
	c.Description = description
	c.IsActive = isActive
	c.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
