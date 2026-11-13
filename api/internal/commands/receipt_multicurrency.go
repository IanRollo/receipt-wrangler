package commands

import (
	"strings"

	"receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/models"
)

const (
	invalidCurrencyMessage = "Currency must be an active ISO 4217 code like USD, EUR or JPY"
	maxExchangeRateSource  = 255
)

// validateOptionalCurrencyCode checks a pointer-style currency setting: nil (omitted) and ""
// (clear back to the default) are valid, anything else must be a code in the table.
func validateOptionalCurrencyCode(code *string) string {
	if code == nil || strings.TrimSpace(*code) == "" {
		return ""
	}
	if !currency.IsValid(*code) {
		return invalidCurrencyMessage
	}
	return ""
}

// ValidateMulticurrency checks the parts of a receipt's currency block that need no database:
// the code's shape and the conversion fields' own constraints. It is not part of Validate, which
// answers a failure with a 500 on the legacy create: repositories.ResolveReceiptCurrency calls it
// first, so every currency problem is one 400 validator error. What depends on the resolved
// currency (amounts by exponent, the base currency's exponent) is checked there once the
// currency is known.
func (receipt *UpsertReceiptCommand) ValidateMulticurrency() map[string]string {
	errors := make(map[string]string)

	if msg := validateOptionalCurrencyCode(&receipt.Currency); len(msg) > 0 {
		errors["currency"] = msg
	}

	if receipt.BaseAmount != nil && receipt.BaseAmount.IsNegative() {
		errors["baseAmount"] = "Base amount cannot be negative"
	}

	if receipt.ExchangeRate != nil && !receipt.ExchangeRate.IsPositive() {
		errors["exchangeRate"] = "Exchange rate must be greater than zero"
	}

	if receipt.ExchangeRateDate != nil && strings.TrimSpace(*receipt.ExchangeRateDate) != "" {
		if _, ok := models.ParseCalendarDay(*receipt.ExchangeRateDate); !ok {
			errors["exchangeRateDate"] = "Exchange rate date must be a date like 2026-09-30"
		}
	}

	if receipt.ExchangeRateSource != nil && tooLong(strings.TrimSpace(*receipt.ExchangeRateSource), maxExchangeRateSource) {
		errors["exchangeRateSource"] = "Exchange rate source is too long"
	}

	return errors
}
