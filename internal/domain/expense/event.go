package expense

import "time"

type EventType string

const (
	EventCreated EventType = "expense.created"
	EventUpdated EventType = "expense.updated"
	EventDeleted EventType = "expense.deleted"
)

type Event struct {
	Type       EventType
	Expense    Expense
	OccurredAt time.Time
}