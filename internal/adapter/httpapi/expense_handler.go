package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
	"github.com/ivsanmendez/ControlDeGastos/internal/port"
)

type ExpenseHandler struct {
	svc port.ExpenseService
}

type createExpenseRequest struct {
	Description string           `json:"description"`
	Amount      float64          `json:"amount"`
	Category    expense.Category `json:"category"`
	Date        time.Time        `json:"date"`
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	e, err := h.svc.CreateExpense(r.Context(), req.Description, req.Amount, req.Category, req.Date)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.svc.ListExpenses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	e, err := h.svc.GetExpense(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteExpense(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}