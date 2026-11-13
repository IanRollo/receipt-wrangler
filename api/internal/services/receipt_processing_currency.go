package services

import (
	"encoding/json"

	"receipt-wrangler/api/internal/commands"
	ccy "receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/repositories"

	"gorm.io/gorm"
)

// aiConversionKeys are receipt keys a model must never set: a conversion is the user's own
// record of what they were charged, not something a scan can know.
var aiConversionKeys = []string{"baseAmount", "exchangeRate", "exchangeRateDate", "exchangeRateSource"}

// sanitizeAiMulticurrency makes the model's JSON safe to unmarshal into an
// UpsertReceiptCommand whatever it said about currency: a `currency` that is not a string is
// dropped (so the scan is never lost to a malformed one) and the conversion keys are removed.
// Anything that does not parse as a JSON object is returned untouched for the caller's own
// error handling.
func sanitizeAiMulticurrency(cleaned string) string {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(cleaned), &fields); err != nil {
		return cleaned
	}

	changed := false
	for _, key := range aiConversionKeys {
		if _, ok := fields[key]; ok {
			delete(fields, key)
			changed = true
		}
	}

	if raw, ok := fields["currency"]; ok {
		var code string
		if err := json.Unmarshal(raw, &code); err != nil {
			delete(fields, "currency")
			changed = true
		}
	}

	if !changed {
		return cleaned
	}

	rewritten, err := json.Marshal(fields)
	if err != nil {
		return cleaned
	}
	return string(rewritten)
}

// normalizeAiCurrency settles what a scan returned for currency. The model's code is kept
// when it is a valid ISO 4217 code and otherwise dropped, so the group's default (then the
// install default) applies when the receipt is stored. Amounts are then rounded half-even to
// the exponent of the currency the receipt will end up in, so a scan is never lost to a
// model that printed more decimals than the currency has (a hand-entered amount is refused
// instead). Conversion fields are cleared.
func normalizeAiCurrency(command *commands.UpsertReceiptCommand, groupId uint, tx *gorm.DB) {
	code := ccy.Normalize(command.Currency)
	if !ccy.IsValid(code) {
		code = ""
	}
	command.Currency = code

	command.BaseAmount = nil
	command.ExchangeRate = nil
	command.ExchangeRateDate = nil
	command.ExchangeRateSource = nil

	rounding := code
	if rounding == "" {
		if groupId > 0 {
			rounding = repositories.GroupBaseCurrency(tx, groupId)
		} else {
			rounding = repositories.InstallDefaultCurrency(tx)
		}
	}

	command.Amount = ccy.RoundToExponent(command.Amount, rounding)
	for index := range command.Items {
		command.Items[index].Amount = ccy.RoundToExponent(command.Items[index].Amount, rounding)
	}
}
