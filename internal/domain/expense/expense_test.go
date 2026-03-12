package expense_test

import (
	"testing"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
)

func TestNew_Valid(t *testing.T) {
	date := time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC)
	e, err := expense.New(1, "Groceries", 50.00, 1, date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if e.UserID != 1 {
		t.Errorf("userID = %d, want 1", e.UserID)
	}
	if e.Description != "Groceries" {
		t.Errorf("description = %q, want %q", e.Description, "Groceries")
	}
	if e.Amount != 50.00 {
		t.Errorf("amount = %v, want %v", e.Amount, 50.00)
	}
	if e.CategoryID != 1 {
		t.Errorf("categoryID = %d, want 1", e.CategoryID)
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

func TestNew_InvalidUserID(t *testing.T) {
	_, err := expense.New(0, "Coffee", 5.00, 1, time.Now())
	if err != expense.ErrInvalidUserID {
		t.Errorf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestNew_EmptyDescription(t *testing.T) {
	_, err := expense.New(1, "", 50.00, 1, time.Now())
	if err != expense.ErrEmptyDescription {
		t.Errorf("expected ErrEmptyDescription, got %v", err)
	}
}

func TestNew_ZeroAmount(t *testing.T) {
	_, err := expense.New(1, "Coffee", 0, 1, time.Now())
	if err != expense.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestNew_NegativeAmount(t *testing.T) {
	_, err := expense.New(1, "Coffee", -10.00, 1, time.Now())
	if err != expense.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestNew_InvalidCategoryID(t *testing.T) {
	_, err := expense.New(1, "Coffee", 5.00, 0, time.Now())
	if err != expense.ErrInvalidCategoryID {
		t.Errorf("expected ErrInvalidCategoryID, got %v", err)
	}
}
