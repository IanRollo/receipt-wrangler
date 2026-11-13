package currency

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// receiptAmountLimit is the first value a decimal(13,3) column cannot hold.
var receiptAmountLimit = decimal.New(1, 10)

// ValidateAmount reports what is wrong with an amount in the given currency, or "" when it is
// acceptable. An amount with more decimal places than the currency's exponent is refused
// rather than rounded: the user's number is the evidence.
func ValidateAmount(amount decimal.Decimal, code string) string {
	exponent := ExponentOf(code)
	if !amount.Equal(amount.Round(int32(exponent))) {
		return AmountPlacesMessage(code)
	}
	if amount.Abs().GreaterThanOrEqual(receiptAmountLimit) {
		return "Amount is too large"
	}
	return ""
}

// AmountPlacesMessage is the validator message for an amount with too many decimal places.
func AmountPlacesMessage(code string) string {
	exponent := ExponentOf(code)
	switch exponent {
	case 0:
		return fmt.Sprintf("%s amounts must be whole numbers", Normalize(code))
	case 1:
		return fmt.Sprintf("%s amounts can have at most 1 decimal place", Normalize(code))
	}
	return fmt.Sprintf("%s amounts can have at most %d decimal places", Normalize(code), exponent)
}

// RoundToExponent rounds half-even (banker's rounding) to the currency's exponent.
func RoundToExponent(amount decimal.Decimal, code string) decimal.Decimal {
	return amount.RoundBank(int32(ExponentOf(code)))
}

// rateFractionDigits is the precision a derived exchange rate is kept to.
const rateFractionDigits = 6

// noteThreshold is the relative gap between a stated base amount and amount x rate beyond
// which a conversionNote is raised.
var noteThreshold = decimal.RequireFromString("0.01")

// DeriveConversion completes a conversion from what the user gave. With only a base amount
// the rate is base / amount to six places (left unset for a zero amount, which has no rate);
// with only a rate the base amount is amount x rate rounded half-even to baseCurrency's
// exponent; with both, both are returned as given. Neither returns nil, nil.
func DeriveConversion(amount decimal.Decimal, baseCurrency string, baseAmount, rate *decimal.Decimal) (*decimal.Decimal, *decimal.Decimal) {
	switch {
	case baseAmount != nil && rate == nil:
		if amount.IsZero() {
			return baseAmount, nil
		}
		derived := baseAmount.Div(amount).Round(rateFractionDigits)
		return baseAmount, &derived
	case baseAmount == nil && rate != nil:
		derived := RoundToExponent(amount.Mul(*rate), baseCurrency)
		return &derived, rate
	}
	return baseAmount, rate
}

// ConversionNote returns a non-blocking warning when a stated base amount and exchange rate
// disagree about the amount by more than 1%, or "" when they agree (or either is missing).
// Card fees make small gaps normal; a large one is usually a typo.
func ConversionNote(amount decimal.Decimal, baseAmount, rate *decimal.Decimal) string {
	if baseAmount == nil || rate == nil {
		return ""
	}

	expected := amount.Mul(*rate)
	if expected.IsZero() {
		if baseAmount.IsZero() {
			return ""
		}
		return "Base amount does not match the amount and exchange rate"
	}

	gap := baseAmount.Sub(expected).Abs().Div(expected.Abs())
	if gap.GreaterThan(noteThreshold) {
		return fmt.Sprintf(
			"Base amount differs from amount x exchange rate (%s) by %s%%",
			expected.RoundBank(2).String(),
			gap.Mul(decimal.NewFromInt(100)).RoundBank(1).String(),
		)
	}
	return ""
}
