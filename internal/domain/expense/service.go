package expense

import (
	"context"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
)

// Repository is the outbound port for expense persistence.
type Repository interface {
	Save(ctx context.Context, e *Expense) error
	FindByID(ctx context.Context, id int64) (*Expense, error)
	FindAll(ctx context.Context) ([]Expense, error)
	FindAllByUser(ctx context.Context, userID int64) ([]Expense, error)
	FindAllDetailed(ctx context.Context) ([]ExpenseDetail, error)
	FindAllDetailedByUser(ctx context.Context, userID int64) ([]ExpenseDetail, error)
	Delete(ctx context.Context, id int64) error
}

// EventPublisher is the outbound port for domain event dispatch.
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

// Service orchestrates expense use cases.
type Service struct {
	repo   Repository
	events EventPublisher
}

func NewService(repo Repository, events EventPublisher) *Service {
	return &Service{repo: repo, events: events}
}

func (s *Service) CreateExpense(ctx context.Context, callerID int64, description string, amount float64, categoryID int64, date time.Time) (*Expense, error) {
	e, err := New(callerID, description, amount, categoryID, date)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, e); err != nil {
		return nil, err
	}
	_ = s.events.Publish(ctx, Event{
		Type:       EventCreated,
		Expense:    *e,
		OccurredAt: time.Now(),
	})
	return e, nil
}

func (s *Service) GetExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64) (*Expense, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if callerRole != user.RoleAdmin && e.UserID != callerID {
		return nil, ErrForbidden
	}
	return e, nil
}

func (s *Service) ListExpenses(ctx context.Context, callerID int64, callerRole user.Role) ([]ExpenseDetail, error) {
	if callerRole == user.RoleAdmin {
		return s.repo.FindAllDetailed(ctx)
	}
	return s.repo.FindAllDetailedByUser(ctx, callerID)
}

func (s *Service) DeleteExpense(ctx context.Context, callerID int64, callerRole user.Role, id int64) error {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if callerRole != user.RoleAdmin && e.UserID != callerID {
		return ErrForbidden
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.events.Publish(ctx, Event{
		Type:       EventDeleted,
		Expense:    *e,
		OccurredAt: time.Now(),
	})
	return nil
}
