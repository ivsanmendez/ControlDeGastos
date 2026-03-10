package httpapi

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

type AuthHandler struct {
	svc port.AuthService
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userResponse struct {
	ID    int64     `json:"id"`
	Email string    `json:"email"`
	Role  user.Role `json:"role"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.svc.Register(r.Context(), req.Email, req.Password, auditInfoFromRequest(r))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailTaken):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrWeakPassword):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	writeJSON(w, http.StatusCreated, userResponse{ID: u.ID, Email: u.Email, Role: u.Role})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, pair, err := h.svc.Login(r.Context(), req.Email, req.Password, auditInfoFromRequest(r))
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pair, err := h.svc.RefreshToken(r.Context(), req.RefreshToken, auditInfoFromRequest(r))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrTokenRevoked):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, user.ErrTokenExpired), errors.Is(err, user.ErrTokenNotFound):
			writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		default:
			writeError(w, http.StatusInternalServerError, "token refresh failed")
		}
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Logout(r.Context(), req.RefreshToken, auditInfoFromRequest(r)); err != nil {
		if errors.Is(err, user.ErrTokenNotFound) {
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
		} else {
			writeError(w, http.StatusInternalServerError, "logout failed")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "no claims in context")
		return
	}

	u, err := h.svc.GetUser(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, userResponse{ID: u.ID, Email: u.Email, Role: u.Role})
}

// auditInfoFromRequest extracts IP and User-Agent from the HTTP request.
func auditInfoFromRequest(r *http.Request) user.AuditInfo {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	// X-Forwarded-For may contain multiple IPs; take the first.
	if idx := strings.Index(ip, ","); idx != -1 {
		ip = strings.TrimSpace(ip[:idx])
	}

	return user.AuditInfo{
		IP:        ip,
		UserAgent: r.Header.Get("User-Agent"),
	}
}
