package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/i18n"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

type ContributionHandler struct {
	svc port.ContributionService
	tr  *i18n.Translator
}

type createContributionRequest struct {
	ContributorID int64                      `json:"contributor_id"`
	CategoryID    int64                      `json:"category_id"`
	Amount        float64                    `json:"amount"`
	Month         int                        `json:"month"`
	Year          int                        `json:"year"`
	PaymentDate   string                     `json:"payment_date"`
	PaymentMethod contribution.PaymentMethod `json:"payment_method"`
}

func (h *ContributionHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeErrorT(w, r, h.tr, http.StatusUnauthorized, "no_claims_in_context")
		return
	}

	var req createContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_payment_date_format")
		return
	}

	c, err := h.svc.CreateContribution(
		r.Context(),
		claims.UserID,
		req.ContributorID,
		req.CategoryID,
		req.Amount,
		req.Month,
		req.Year,
		paymentDate,
		req.PaymentMethod,
	)
	if err != nil {
		if errors.Is(err, contribution.ErrDuplicate) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *ContributionHandler) List(w http.ResponseWriter, r *http.Request) {
	contributorIDStr := r.URL.Query().Get("contributor_id")
	yearStr := r.URL.Query().Get("year")

	var contributorID int64
	if contributorIDStr != "" {
		var err error
		contributorID, err = strconv.ParseInt(contributorIDStr, 10, 64)
		if err != nil {
			writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_contributor_id")
			return
		}
	}

	var year int
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_year")
			return
		}
	}

	contributions, err := h.svc.ListContributions(r.Context(), contributorID, year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, contributions)
}

func (h *ContributionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	c, err := h.svc.GetContribution(r.Context(), id)
	if err != nil {
		if errors.Is(err, contribution.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "contribution_not_found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type updateContributionRequest struct {
	ContributorID int64                      `json:"contributor_id"`
	CategoryID    int64                      `json:"category_id"`
	Amount        float64                    `json:"amount"`
	Month         int                        `json:"month"`
	Year          int                        `json:"year"`
	PaymentDate   string                     `json:"payment_date"`
	PaymentMethod contribution.PaymentMethod `json:"payment_method"`
}

func (h *ContributionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	var req updateContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_request_body")
		return
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_payment_date_format")
		return
	}

	c, err := h.svc.UpdateContribution(
		r.Context(),
		id,
		req.ContributorID,
		req.CategoryID,
		req.Amount,
		req.Month,
		req.Year,
		paymentDate,
		req.PaymentMethod,
	)
	if err != nil {
		if errors.Is(err, contribution.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "contribution_not_found")
		} else if errors.Is(err, contribution.ErrDuplicate) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ContributionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorT(w, r, h.tr, http.StatusBadRequest, "invalid_id")
		return
	}

	if err := h.svc.DeleteContribution(r.Context(), id); err != nil {
		if errors.Is(err, contribution.ErrNotFound) {
			writeErrorT(w, r, h.tr, http.StatusNotFound, "contribution_not_found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
