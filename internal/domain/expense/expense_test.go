package expense_test

import (
	"testing"
	"time"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
)

func TestNew_Valid(t *testing.T) {
	date := time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC)
	e, err := expense.New("Groceries", 50.00, expense.CategoryFood, date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if e.Description != "Groceries" {
		t.Errorf("description = %q, want %q", e.Description, "Groceries")
	}
	if e.Amount != 50.00 {
		t.Errorf("amount = %v, want %v", e.Amount, 50.00)
	}
	if e.Category != expense.CategoryFood {
		t.Errorf("category = %v, want %v", e.Category, expense.CategoryFood)
	}
	if e.Date != date {
		t.Errorf("date = %v, want %v", e.Date, date)
	}
	if e.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if e.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestNew_EmptyDescription(t *testing.T) {
	_, err := expense.New("", 50.00, expense.CategoryFood, time.Now())
	if err != expense.ErrEmptyDescription {
		t.Errorf("expected ErrEmptyDescription, got %v", err)
	}
}

func TestNew_ZeroAmount(t *testing.T) {
	_, err := expense.New("Coffee", 0, expense.CategoryFood, time.Now())
	if err != expense.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestNew_NegativeAmount(t *testing.T) {
	_, err := expense.New("Coffee", -10.00, expense.CategoryFood, time.Now())
	if err != expense.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}
