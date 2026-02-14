package port

import (
	"context"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
)

// ExpenseRepository is the driven port for expense persistence.
type ExpenseRepository = expense.Repository

// EventPublisher is the driven port for publishing domain events.
type EventPublisher = expense.EventPublisher

// EventSubscriber reacts to domain events.
type EventSubscriber interface {
	Handle(ctx context.Context, event expense.Event) error
}