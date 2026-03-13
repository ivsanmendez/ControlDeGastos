package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/i18n"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

type ExpenseHandler struct {
	svc port.ExpenseService
	tr  *i18n.Translator
}

type createExpenseRequest struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CategoryID  int64     `json:"category_id"`
	Date        time.Time `json:"date"`
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}
	e, err := h.svc.CreateExpense(r.Context(), claims.UserID, req.Description, req.Amount, req.CategoryID, req.Date)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	expenses, err := h.svc.ListExpenses(r.Context(), claims.UserID, claims.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}
	e, err := h.svc.GetExpense(r.Context(), claims.UserID, claims.Role, id)
	if err != nil {
		if errors.Is(err, expense.ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
		} else {
			writeError(w, http.StatusNotFound, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, e)
}

type updateExpenseRequest struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CategoryID  int64     `json:"category_id"`
	Date        time.Time `json:"date"`
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	var req updateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}

	e, err := h.svc.UpdateExpense(r.Context(), claims.UserID, claims.Role, id, req.Description, req.Amount, req.CategoryID, req.Date)
	if err != nil {
		if errors.Is(err, expense.ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
		} else if errors.Is(err, expense.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "expense_not_found")
		} else {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}
	if err := h.svc.DeleteExpense(r.Context(), claims.UserID, claims.Role, id); err != nil {
		if errors.Is(err, expense.ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
		} else {
			writeError(w, http.StatusNotFound, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
