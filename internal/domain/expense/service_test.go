package expense_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
)

// fakeRepo is an in-memory implementation of expense.Repository.
type fakeRepo struct {
	data    map[int64]*expense.Expense
	nextID  int64
	saveErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{data: make(map[int64]*expense.Expense), nextID: 1}
}

func (r *fakeRepo) Save(_ context.Context, e *expense.Expense) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	e.ID = r.nextID
	r.nextID++
	cp := *e
	r.data[e.ID] = &cp
	return nil
}

func (r *fakeRepo) FindByID(_ context.Context, id int64) (*expense.Expense, error) {
	e, ok := r.data[id]
	if !ok {
		return nil, expense.ErrNotFound
	}
	cp := *e
	return &cp, nil
}

func (r *fakeRepo) FindAll(_ context.Context) ([]expense.Expense, error) {
	result := make([]expense.Expense, 0, len(r.data))
	for _, e := range r.data {
		result = append(result, *e)
	}
	return result, nil
}

func (r *fakeRepo) Delete(_ context.Context, id int64) error {
	if _, ok := r.data[id]; !ok {
		return expense.ErrNotFound
	}
	delete(r.data, id)
	return nil
}

// fakePublisher records published events.
type fakePublisher struct {
	events []expense.Event
}

func (p *fakePublisher) Publish(_ context.Context, e expense.Event) error {
	p.events = append(p.events, e)
	return nil
}

func newService() (*expense.Service, *fakeRepo, *fakePublisher) {
	repo := newFakeRepo()
	pub := &fakePublisher{}
	svc := expense.NewService(repo, pub)
	return svc, repo, pub
}

var (
	ctx      = context.Background()
	testDate = time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC)
)

func TestCreateExpense_HappyPath(t *testing.T) {
	svc, repo, pub := newService()

	e, err := svc.CreateExpense(ctx, "Lunch", 12.50, expense.CategoryFood, testDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.ID == 0 {
		t.Error("expected ID to be set after save")
	}
	if _, ok := repo.data[e.ID]; !ok {
		t.Error("expense not found in repo after create")
	}
	if len(pub.events) != 1 || pub.events[0].Type != expense.EventCreated {
		t.Errorf("expected one EventCreated, got %v", pub.events)
	}
}

func TestCreateExpense_InvalidInput(t *testing.T) {
	svc, _, pub := newService()

	_, err := svc.CreateExpense(ctx, "", 12.50, expense.CategoryFood, testDate)
	if !errors.Is(err, expense.ErrEmptyDescription) {
		t.Errorf("expected ErrEmptyDescription, got %v", err)
	}
	if len(pub.events) != 0 {
		t.Error("no event should be published on invalid input")
	}
}

func TestCreateExpense_RepoError(t *testing.T) {
	svc, repo, _ := newService()
	repo.saveErr = errors.New("db unavailable")

	_, err := svc.CreateExpense(ctx, "Taxi", 8.00, expense.CategoryTransport, testDate)
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

func TestGetExpense_Found(t *testing.T) {
	svc, _, _ := newService()
	created, _ := svc.CreateExpense(ctx, "Bus", 2.50, expense.CategoryTransport, testDate)

	got, err := svc.GetExpense(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("id = %d, want %d", got.ID, created.ID)
	}
}

func TestGetExpense_NotFound(t *testing.T) {
	svc, _, _ := newService()

	_, err := svc.GetExpense(ctx, 999)
	if !errors.Is(err, expense.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListExpenses(t *testing.T) {
	svc, _, _ := newService()
	svc.CreateExpense(ctx, "Coffee", 3.00, expense.CategoryFood, testDate)
	svc.CreateExpense(ctx, "Metro", 1.50, expense.CategoryTransport, testDate)

	list, err := svc.ListExpenses(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 expenses, got %d", len(list))
	}
}

func TestDeleteExpense_HappyPath(t *testing.T) {
	svc, repo, pub := newService()
	created, _ := svc.CreateExpense(ctx, "Dinner", 30.00, expense.CategoryFood, testDate)
	pub.events = nil // reset after create

	err := svc.DeleteExpense(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := repo.data[created.ID]; ok {
		t.Error("expense should be removed from repo after delete")
	}
	if len(pub.events) != 1 || pub.events[0].Type != expense.EventDeleted {
		t.Errorf("expected one EventDeleted, got %v", pub.events)
	}
}

func TestDeleteExpense_NotFound(t *testing.T) {
	svc, _, _ := newService()

	err := svc.DeleteExpense(ctx, 999)
	if !errors.Is(err, expense.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
