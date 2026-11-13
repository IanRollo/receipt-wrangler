package handlers

import (
	"net/http"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/stacktrace"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// currencyView is one entry of GET /currency (swagger: CurrencyInfo).
type currencyView struct {
	Code     string `json:"code"`
	Exponent int    `json:"exponent"`
	Symbol   string `json:"symbol,omitempty"`
}

// GetCurrencies serves the table of active ISO 4217 codes the server accepts, with the
// number of decimal places each allows, so a client's pickers and inputs never carry a copy
// of the table. Any signed-in user may read it.
func GetCurrencies(w http.ResponseWriter, r *http.Request) {
	HandleRequest(structs.Handler{
		ErrorMessage: "Error getting currencies",
		Writer:       w,
		Request:      r,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			all := currency.All()
			views := make([]currencyView, 0, len(all))
			for _, info := range all {
				views = append(views, currencyView{Code: info.Code, Exponent: info.Exponent, Symbol: info.Symbol})
			}

			bytes, err := utils.MarshalResponseData(views)
			if err != nil {
				return http.StatusInternalServerError, stacktrace.WithStack(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
			return 0, nil
		},
	})
}
