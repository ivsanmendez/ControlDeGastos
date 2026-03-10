package contributor

import (
	"context"
	"time"
)

// Service orchestrates contributor use cases.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateContributor(ctx context.Context, callerID int64, houseNumber, name, phone string) (*Contributor, error) {
	c, err := New(callerID, houseNumber, name, phone)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetContributor(ctx context.Context, id int64) (*Contributor, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListContributors(ctx context.Context) ([]Contributor, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) UpdateContributor(ctx context.Context, id int64, name, phone string) (*Contributor, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, ErrEmptyName
	}

	c.Name = name
	c.Phone = phone
	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteContributor(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
