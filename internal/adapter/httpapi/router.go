package httpapi

import (
	"net/http"

	"github.com/ivsanmendez/ControlDeGastos/internal/port"
)

// RegisterRoutes wires all HTTP routes onto the given mux.
func RegisterRoutes(mux *http.ServeMux, expenseSvc port.ExpenseService) {
	h := &ExpenseHandler{svc: expenseSvc}

	mux.HandleFunc("GET /health", Health)

	mux.HandleFunc("POST /expenses", h.Create)
	mux.HandleFunc("GET /expenses", h.List)
	mux.HandleFunc("GET /expenses/{id}", h.GetByID)
	mux.HandleFunc("DELETE /expenses/{id}", h.Delete)
}