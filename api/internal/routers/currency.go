package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// BuildCurrencyRouter serves the ISO 4217 currency table (multicurrency, fork feature),
// mounted at /api/currency.
func BuildCurrencyRouter() *chi.Mux {
	currencyRouter := chi.NewRouter()

	currencyRouter.Use(middleware.UnifiedAuthMiddleware)
	currencyRouter.Get("/", handlers.GetCurrencies)

	return currencyRouter
}
