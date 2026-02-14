package eventbus

import (
	"context"
	"log"
	"sync"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
	"github.com/ivsanmendez/ControlDeGastos/internal/port"
)

// InMemBus is a synchronous in-memory event bus.
// Implements expense.EventPublisher.
type InMemBus struct {
	mu          sync.RWMutex
	subscribers []port.EventSubscriber
}

func New() *InMemBus {
	return &InMemBus{}
}

func (b *InMemBus) Subscribe(sub port.EventSubscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, sub)
}

func (b *InMemBus) Publish(ctx context.Context, event expense.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subscribers {
		if err := sub.Handle(ctx, event); err != nil {
			log.Printf("event subscriber error: %v", err)
		}
	}
	return nil
}