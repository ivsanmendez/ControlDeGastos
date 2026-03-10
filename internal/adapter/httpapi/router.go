package httpapi

import (
	"net/http"

	jwtadapter "github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/jwt"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

// RegisterRoutes wires all HTTP routes onto the given mux.
func RegisterRoutes(mux *http.ServeMux, expenseSvc port.ExpenseService, authSvc port.AuthService, contribSvc port.ContributionService, contributorSvc port.ContributorService, jwtIssuer *jwtadapter.Issuer, signer port.ReceiptSigner) {
	auth := RequireAuth(jwtIssuer)
	expH := &ExpenseHandler{svc: expenseSvc}
	authH := &AuthHandler{svc: authSvc}
	contribH := &ContributionHandler{svc: contribSvc}
	contributorH := &ContributorHandler{svc: contributorSvc}
	receiptH := &ReceiptHandler{contribSvc: contribSvc, contributorSvc: contributorSvc, signer: signer}

	// Public routes
	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.HandleFunc("POST /auth/refresh", authH.Refresh)

	// Protected auth routes
	mux.Handle("POST /auth/logout", Chain(http.HandlerFunc(authH.Logout), auth))
	mux.Handle("GET /auth/me", Chain(http.HandlerFunc(authH.Me), auth))

	// Protected expense routes
	mux.Handle("POST /expenses", Chain(
		http.HandlerFunc(expH.Create),
		auth, RequirePermission(user.PermExpenseCreate),
	))
	mux.Handle("GET /expenses", Chain(
		http.HandlerFunc(expH.List),
		auth, RequirePermission(user.PermExpenseReadOwn),
	))
	mux.Handle("GET /expenses/{id}", Chain(
		http.HandlerFunc(expH.GetByID),
		auth, RequirePermission(user.PermExpenseReadOwn),
	))
	mux.Handle("DELETE /expenses/{id}", Chain(
		http.HandlerFunc(expH.Delete),
		auth, RequirePermission(user.PermExpenseDeleteOwn),
	))

	// Protected contributor routes
	mux.Handle("POST /contributors", Chain(
		http.HandlerFunc(contributorH.Create),
		auth, RequirePermission(user.PermContributorCreate),
	))
	mux.Handle("GET /contributors", Chain(
		http.HandlerFunc(contributorH.List),
		auth, RequirePermission(user.PermContributorRead),
	))
	mux.Handle("GET /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.GetByID),
		auth, RequirePermission(user.PermContributorRead),
	))
	mux.Handle("PUT /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.Update),
		auth, RequirePermission(user.PermContributorUpdate),
	))
	mux.Handle("DELETE /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.Delete),
		auth, RequirePermission(user.PermContributorDelete),
	))

	// Protected contribution routes
	mux.Handle("POST /contributions", Chain(
		http.HandlerFunc(contribH.Create),
		auth, RequirePermission(user.PermContributionCreate),
	))
	mux.Handle("GET /contributions", Chain(
		http.HandlerFunc(contribH.List),
		auth, RequirePermission(user.PermContributionRead),
	))
	mux.Handle("GET /contributions/{id}", Chain(
		http.HandlerFunc(contribH.GetByID),
		auth, RequirePermission(user.PermContributionRead),
	))
	mux.Handle("DELETE /contributions/{id}", Chain(
		http.HandlerFunc(contribH.Delete),
		auth, RequirePermission(user.PermContributionDelete),
	))

	// Receipt digital signature (POST: requires password for key decryption)
	mux.Handle("POST /contributions/receipt-signature", Chain(
		http.HandlerFunc(receiptH.ReceiptSignature),
		auth, RequirePermission(user.PermContributionRead),
	))
}
