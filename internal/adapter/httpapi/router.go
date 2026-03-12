package httpapi

import (
	"net/http"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/i18n"
	jwtadapter "github.com/ivsanmendez/ControlDeContabilidad/internal/adapter/jwt"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/user"
	"github.com/ivsanmendez/ControlDeContabilidad/internal/port"
)

// RegisterRoutes wires all HTTP routes onto the given mux.
func RegisterRoutes(mux *http.ServeMux, expenseSvc port.ExpenseService, authSvc port.AuthService, contribSvc port.ContributionService, contributorSvc port.ContributorService, categorySvc port.CategoryService, expCatSvc port.ExpenseCategoryService, receiptSvc port.ReceiptFolioService, reportSvc port.ReportService, jwtIssuer *jwtadapter.Issuer, signer port.ReceiptSigner, tr *i18n.Translator) {
	auth := RequireAuth(jwtIssuer, tr)
	expH := &ExpenseHandler{svc: expenseSvc, tr: tr}
	authH := &AuthHandler{svc: authSvc, tr: tr}
	contribH := &ContributionHandler{svc: contribSvc, tr: tr}
	contributorH := &ContributorHandler{svc: contributorSvc, tr: tr}
	categoryH := &CategoryHandler{svc: categorySvc, tr: tr}
	expCatH := &ExpenseCategoryHandler{svc: expCatSvc, tr: tr}
	receiptH := &ReceiptHandler{contribSvc: contribSvc, contributorSvc: contributorSvc, receiptSvc: receiptSvc, signer: signer, tr: tr}
	reportH := &ReportHandler{svc: reportSvc, tr: tr}

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
		auth, RequirePermission(user.PermExpenseCreate, tr),
	))
	mux.Handle("GET /expenses", Chain(
		http.HandlerFunc(expH.List),
		auth, RequirePermission(user.PermExpenseReadOwn, tr),
	))
	mux.Handle("GET /expenses/{id}", Chain(
		http.HandlerFunc(expH.GetByID),
		auth, RequirePermission(user.PermExpenseReadOwn, tr),
	))
	mux.Handle("DELETE /expenses/{id}", Chain(
		http.HandlerFunc(expH.Delete),
		auth, RequirePermission(user.PermExpenseDeleteOwn, tr),
	))

	// Protected contributor routes
	mux.Handle("POST /contributors", Chain(
		http.HandlerFunc(contributorH.Create),
		auth, RequirePermission(user.PermContributorCreate, tr),
	))
	mux.Handle("GET /contributors", Chain(
		http.HandlerFunc(contributorH.List),
		auth, RequirePermission(user.PermContributorRead, tr),
	))
	mux.Handle("GET /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.GetByID),
		auth, RequirePermission(user.PermContributorRead, tr),
	))
	mux.Handle("PUT /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.Update),
		auth, RequirePermission(user.PermContributorUpdate, tr),
	))
	mux.Handle("DELETE /contributors/{id}", Chain(
		http.HandlerFunc(contributorH.Delete),
		auth, RequirePermission(user.PermContributorDelete, tr),
	))

	// Protected contribution category routes
	mux.Handle("POST /contribution-categories", Chain(
		http.HandlerFunc(categoryH.Create),
		auth, RequirePermission(user.PermCategoryCreate, tr),
	))
	mux.Handle("GET /contribution-categories", Chain(
		http.HandlerFunc(categoryH.List),
		auth, RequirePermission(user.PermCategoryRead, tr),
	))
	mux.Handle("GET /contribution-categories/{id}", Chain(
		http.HandlerFunc(categoryH.GetByID),
		auth, RequirePermission(user.PermCategoryRead, tr),
	))
	mux.Handle("PUT /contribution-categories/{id}", Chain(
		http.HandlerFunc(categoryH.Update),
		auth, RequirePermission(user.PermCategoryUpdate, tr),
	))
	mux.Handle("DELETE /contribution-categories/{id}", Chain(
		http.HandlerFunc(categoryH.Delete),
		auth, RequirePermission(user.PermCategoryDelete, tr),
	))

	// Protected expense category routes
	mux.Handle("POST /expense-categories", Chain(
		http.HandlerFunc(expCatH.Create),
		auth, RequirePermission(user.PermExpenseCategoryCreate, tr),
	))
	mux.Handle("GET /expense-categories", Chain(
		http.HandlerFunc(expCatH.List),
		auth, RequirePermission(user.PermExpenseCategoryRead, tr),
	))
	mux.Handle("GET /expense-categories/{id}", Chain(
		http.HandlerFunc(expCatH.GetByID),
		auth, RequirePermission(user.PermExpenseCategoryRead, tr),
	))
	mux.Handle("PUT /expense-categories/{id}", Chain(
		http.HandlerFunc(expCatH.Update),
		auth, RequirePermission(user.PermExpenseCategoryUpdate, tr),
	))
	mux.Handle("DELETE /expense-categories/{id}", Chain(
		http.HandlerFunc(expCatH.Delete),
		auth, RequirePermission(user.PermExpenseCategoryDelete, tr),
	))

	// Protected contribution routes
	mux.Handle("POST /contributions", Chain(
		http.HandlerFunc(contribH.Create),
		auth, RequirePermission(user.PermContributionCreate, tr),
	))
	mux.Handle("GET /contributions", Chain(
		http.HandlerFunc(contribH.List),
		auth, RequirePermission(user.PermContributionRead, tr),
	))
	mux.Handle("GET /contributions/{id}", Chain(
		http.HandlerFunc(contribH.GetByID),
		auth, RequirePermission(user.PermContributionRead, tr),
	))
	mux.Handle("DELETE /contributions/{id}", Chain(
		http.HandlerFunc(contribH.Delete),
		auth, RequirePermission(user.PermContributionDelete, tr),
	))

	// Receipt digital signature (POST: requires password for key decryption)
	mux.Handle("POST /contributions/receipt-signature", Chain(
		http.HandlerFunc(receiptH.ReceiptSignature),
		auth, RequirePermission(user.PermContributionRead, tr),
	))

	// Receipt folio verification
	mux.Handle("GET /receipts/verify/{folio}", Chain(
		http.HandlerFunc(receiptH.VerifyReceipt),
		auth, RequirePermission(user.PermReceiptVerify, tr),
	))

	// Reports
	mux.Handle("GET /reports/monthly-balance", Chain(
		http.HandlerFunc(reportH.MonthlyBalance),
		auth, RequirePermission(user.PermReportRead, tr),
	))
}
