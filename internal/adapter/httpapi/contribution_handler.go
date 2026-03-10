package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/contribution"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

type ContributionHandler struct {
	svc port.ContributionService
}

type createContributionRequest struct {
	ContributorID int64                      `json:"contributor_id"`
	Amount        float64                    `json:"amount"`
	Month         int                        `json:"month"`
	Year          int                        `json:"year"`
	PaymentDate   string                     `json:"payment_date"`
	PaymentMethod contribution.PaymentMethod `json:"payment_method"`
}

func (h *ContributionHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "no claims in context")
		return
	}

	var req createContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid payment_date format, expected YYYY-MM-DD")
		return
	}

	c, err := h.svc.CreateContribution(
		r.Context(),
		claims.UserID,
		req.ContributorID,
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
			writeError(w, http.StatusBadRequest, "invalid contributor_id")
			return
		}
	}

	var year int
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid year")
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
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	c, err := h.svc.GetContribution(r.Context(), id)
	if err != nil {
		if errors.Is(err, contribution.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ContributionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.DeleteContribution(r.Context(), id); err != nil {
		if errors.Is(err, contribution.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
