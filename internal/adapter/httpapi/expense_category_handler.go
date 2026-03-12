package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/i18n"
	ec "github.com/ivsanmendez/ControlDeContabilidad/internal/domain/expense_category"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

type ExpenseCategoryHandler struct {
	svc port.ExpenseCategoryService
	tr  *i18n.Translator
}

type createExpenseCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateExpenseCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func (h *ExpenseCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	var req createExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}

	c, err := h.svc.CreateCategory(r.Context(), claims.UserID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, ec.ErrDuplicate) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *ExpenseCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

func (h *ExpenseCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	c, err := h.svc.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, ec.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "expense_category_not_found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ExpenseCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	var req updateExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}

	c, err := h.svc.UpdateCategory(r.Context(), id, req.Name, req.Description, req.IsActive)
	if err != nil {
		if errors.Is(err, ec.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "expense_category_not_found")
		} else if errors.Is(err, ec.ErrDuplicate) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ExpenseCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	if err := h.svc.DeleteCategory(r.Context(), id); err != nil {
		if errors.Is(err, ec.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "expense_category_not_found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
