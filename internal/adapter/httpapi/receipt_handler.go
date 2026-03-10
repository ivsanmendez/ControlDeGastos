package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

// ReceiptHandler serves digitally signed receipt payloads.
type ReceiptHandler struct {
	contribSvc     port.ContributionService
	contributorSvc port.ContributorService
	signer         port.ReceiptSigner
}

type receiptSignRequest struct {
	ContributorID int64  `json:"contributor_id"`
	Year          int    `json:"year"`
	Password      string `json:"password"`
	SignerName    string `json:"signer_name"`
}

type receiptPayment struct {
	Month  int     `json:"month"`
	Amount float64 `json:"amount"`
}

type receiptData struct {
	ContributorID   int64            `json:"contributor_id"`
	HouseNumber     string           `json:"house_number"`
	ContributorName string           `json:"contributor_name"`
	Year            int              `json:"year"`
	Payments        []receiptPayment `json:"payments"`
	Total           float64          `json:"total"`
	SignerName      string           `json:"signer_name"`
	GeneratedAt     time.Time        `json:"generated_at"`
}

type receiptSignatureResponse struct {
	Data        receiptData `json:"data"`
	Signature   string      `json:"signature"`
	Certificate string      `json:"certificate"`
}

// ReceiptSignature handles POST /contributions/receipt-signature.
func (h *ReceiptHandler) ReceiptSignature(w http.ResponseWriter, r *http.Request) {
	if !h.signer.Available() {
		writeError(w, http.StatusServiceUnavailable, "receipt signing is not configured")
		return
	}

	var req receiptSignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ContributorID == 0 || req.Year == 0 {
		writeError(w, http.StatusBadRequest, "contributor_id and year are required")
		return
	}
	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}
	if req.SignerName == "" {
		writeError(w, http.StatusBadRequest, "signer_name is required")
		return
	}

	// Fetch contributor info
	contrib, err := h.contributorSvc.GetContributor(r.Context(), req.ContributorID)
	if err != nil {
		writeError(w, http.StatusNotFound, "contributor not found")
		return
	}

	// Fetch contributions for that year
	contributions, err := h.contribSvc.ListContributions(r.Context(), req.ContributorID, req.Year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load contributions")
		return
	}

	// Build receipt data
	var payments []receiptPayment
	var total float64
	for _, c := range contributions {
		payments = append(payments, receiptPayment{Month: c.Month, Amount: c.Amount})
		total += c.Amount
	}

	data := receiptData{
		ContributorID:   req.ContributorID,
		HouseNumber:     contrib.HouseNumber,
		ContributorName: contrib.Name,
		Year:            req.Year,
		Payments:        payments,
		Total:           total,
		SignerName:      req.SignerName,
		GeneratedAt:     time.Now().UTC(),
	}

	// Build canonical JSON for signing
	canonical, err := json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to serialize receipt data")
		return
	}

	sig, err := h.signer.Sign(canonical, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "decrypt") {
			writeError(w, http.StatusUnauthorized, "invalid certificate password")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to sign receipt")
		return
	}

	resp := receiptSignatureResponse{
		Data:        data,
		Signature:   base64.StdEncoding.EncodeToString(sig),
		Certificate: base64.StdEncoding.EncodeToString(h.signer.Certificate()),
	}

	writeJSON(w, http.StatusOK, resp)
}
