package services

import (
	"sort"

	ccy "receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"github.com/shopspring/decimal"
)

// claimReceiptFigures is what a claim's linked receipts add up to in the claim's base
// currency. A receipt counts at its amount when it is in the base currency (an empty
// currency, which only a pre-multicurrency row can have, counts as the base), at its
// recorded baseAmount when it was converted, and otherwise is unconverted: it is left out
// of counted and reported per currency instead, never summed into it.
type claimReceiptFigures struct {
	counted          decimal.Decimal
	unconvertedCount int
	unconverted      map[string]decimal.Decimal
}

func foldClaimReceipts(receipts []repositories.LinkedReceipt, baseCurrency string) claimReceiptFigures {
	figures := claimReceiptFigures{counted: decimal.Zero, unconverted: map[string]decimal.Decimal{}}

	for _, receipt := range receipts {
		switch {
		case receipt.Currency == "" || receipt.Currency == baseCurrency:
			figures.counted = figures.counted.Add(receipt.Amount)
		case receipt.BaseAmount != nil:
			figures.counted = figures.counted.Add(*receipt.BaseAmount)
		default:
			figures.unconvertedCount++
			figures.unconverted[receipt.Currency] = figures.unconverted[receipt.Currency].Add(receipt.Amount)
		}
	}

	return figures
}

// unconvertedTotalViews lists the unconverted totals by currency code, as a non-nil slice.
func (figures claimReceiptFigures) unconvertedTotalViews() []structs.ClaimCurrencyTotal {
	return currencyTotalViews(figures.unconverted)
}

// currencyTotalViews renders per-currency totals ordered by code, each at its own currency's
// number of decimal places. Always non-nil, so it serializes as [].
func currencyTotalViews(totals map[string]decimal.Decimal) []structs.ClaimCurrencyTotal {
	codes := make([]string, 0, len(totals))
	for code := range totals {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	views := make([]structs.ClaimCurrencyTotal, 0, len(codes))
	for _, code := range codes {
		views = append(views, structs.ClaimCurrencyTotal{Currency: code, Total: formatMoneyIn(totals[code], code)})
	}
	return views
}

// formatMoneyIn renders money with the currency's own number of decimal places (JPY none,
// USD two, JOD three).
func formatMoneyIn(value decimal.Decimal, code string) string {
	return value.StringFixed(int32(ccy.ExponentOf(code)))
}
