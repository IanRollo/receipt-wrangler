package services

import (
	"sort"
	"strings"
	"time"

	ccy "receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/reporting/render"

	"github.com/shopspring/decimal"
)

// Multicurrency in the claim packet (fork feature). The claim is reckoned in its base
// currency; a receipt in another currency counts toward the figures only through the
// conversion the user recorded, and one with none is listed apart as "Not converted" so the
// packet never looks complete when it is not.

// packetCurrencyOf is a receipt's currency code; an empty one (a row from before
// multicurrency that was never backfilled) counts as the claim's base currency.
func (builder *claimPacketBuilder) packetCurrencyOf(receipt models.Receipt) string {
	if receipt.Currency == "" {
		return builder.claim.Currency
	}
	return receipt.Currency
}

// countedAmount is what a receipt adds to the claim's figures, in the base currency: its
// amount when it is in the base currency, its recorded base amount when converted, and
// nothing (false) when it is unconverted.
func (builder *claimPacketBuilder) countedAmount(receipt models.Receipt) (decimal.Decimal, bool) {
	if builder.packetCurrencyOf(receipt) == builder.claim.Currency {
		return receipt.Amount, true
	}
	if receipt.ConvertedTo(builder.claim.Currency) {
		return *receipt.BaseAmount, true
	}
	return decimal.Zero, false
}

// countedAmountString is a receipt's counted amount in the base currency, "" when it is not
// counted.
func (builder *claimPacketBuilder) countedAmountString(receipt models.Receipt) string {
	amount, counted := builder.countedAmount(receipt)
	if !counted {
		return ""
	}
	return builder.moneyIn(amount, builder.claim.Currency)
}

// moneyIn formats an amount in a currency. The install's own currency keeps the format the
// app has always printed (its configured symbol and separators, two places), so a
// single-currency packet reads exactly as before; any other currency, and any currency whose
// exponent is not two, prints as a number at its own number of places followed by its code.
func (builder *claimPacketBuilder) moneyIn(value decimal.Decimal, code string) string {
	if code == "" {
		code = builder.claim.Currency
	}

	exponent := ccy.ExponentOf(code)
	if code == builder.installCurrency && exponent == 2 {
		return render.FormatCurrency(value, builder.currency)
	}

	return formatPacketNumber(value, exponent, builder.currency) + " " + code
}

// moneyFor formats amounts of one receipt: in the receipt's own currency.
func (builder *claimPacketBuilder) moneyFor(receipt models.Receipt) func(decimal.Decimal) string {
	code := builder.packetCurrencyOf(receipt)
	return func(value decimal.Decimal) string { return builder.moneyIn(value, code) }
}

// formatPacketNumber groups and separates a number per the app's configured separators, at
// exactly `exponent` decimal places.
func formatPacketNumber(value decimal.Decimal, exponent int, format reporting.CurrencyFormat) string {
	thousands := format.ThousandsSeparator
	if thousands == "" {
		thousands = ","
	}
	decimalSeparator := format.DecimalSeparator
	if decimalSeparator == "" {
		decimalSeparator = "."
	}

	fixed := value.StringFixed(int32(exponent))
	sign := ""
	if strings.HasPrefix(fixed, "-") {
		sign = "-"
		fixed = fixed[1:]
	}

	integerPart, fraction, _ := strings.Cut(fixed, ".")

	var grouped strings.Builder
	for index, digit := range integerPart {
		if index > 0 && (len(integerPart)-index)%3 == 0 {
			grouped.WriteString(thousands)
		}
		grouped.WriteRune(digit)
	}

	if exponent == 0 {
		return sign + grouped.String()
	}
	return sign + grouped.String() + decimalSeparator + fraction
}

// addConversionFields puts the receipt's currency and, when it is not in the base currency,
// its recorded conversion (or the lack of one) right after its Amount in the per-receipt
// section.
func (builder *claimPacketBuilder) addConversionFields(section *packetReceiptSection, receipt models.Receipt) {
	code := builder.packetCurrencyOf(receipt)
	base := builder.claim.Currency

	extra := []packetField{{"Currency", code}}
	if code != base {
		if receipt.BaseAmount == nil {
			extra = append(extra, packetField{"Conversion", "Not converted: no amount in " + base + " was recorded, so this receipt is not counted in the claim's figures"})
		} else {
			// The conversion is shown in the currency it was recorded in, which is the claim's
			// unless the group's base currency changed afterwards.
			recorded := receipt.BaseAmountCurrency
			if recorded == "" {
				recorded = base
			}
			if recorded != base {
				extra = append(extra, packetField{"Conversion", "Not converted to " + base + ": the recorded conversion is to " + recorded + ", so this receipt is not counted in the claim's figures"})
			}
			extra = append(extra, packetField{"Amount in " + recorded, builder.moneyIn(*receipt.BaseAmount, recorded)})
			if receipt.ExchangeRate != nil {
				extra = append(extra, packetField{"Exchange rate", "1 " + code + " = " + receipt.ExchangeRate.String() + " " + recorded})
			}
			if receipt.ExchangeRateDate != nil && !receipt.ExchangeRateDate.IsZero() {
				extra = append(extra, packetField{"Rate date", receipt.ExchangeRateDate.UTC().Format(time.DateOnly)})
			}
			if receipt.ExchangeRateSource != "" {
				extra = append(extra, packetField{"Rate source", receipt.ExchangeRateSource})
			}
			if receipt.ConversionNote != "" {
				extra = append(extra, packetField{"Conversion note", receipt.ConversionNote})
			}
		}
	}

	fields := make([]packetField, 0, len(section.Fields)+len(extra))
	for _, field := range section.Fields {
		fields = append(fields, field)
		if field.Label == "Amount" {
			fields = append(fields, extra...)
			extra = nil
		}
	}
	section.Fields = append(fields, extra...)
}

// packetNotConvertedGroup lists the receipts of one currency that have no recorded conversion.
type packetNotConvertedGroup struct {
	Currency string
	Rows     []packetReceiptRow
	Total    string
}

// notConverted groups the unconverted receipts by currency, ordered by code.
func (builder *claimPacketBuilder) notConverted(sections []packetReceiptSection) []packetNotConvertedGroup {
	byCurrency := map[string]*packetNotConvertedGroup{}
	totals := map[string]decimal.Decimal{}

	for index, receipt := range builder.receipts {
		if _, counted := builder.countedAmount(receipt); counted {
			continue
		}

		code := builder.packetCurrencyOf(receipt)
		group, ok := byCurrency[code]
		if !ok {
			group = &packetNotConvertedGroup{Currency: code}
			byCurrency[code] = group
		}
		group.Rows = append(group.Rows, sections[index].Row)
		totals[code] = totals[code].Add(receipt.Amount)
	}

	codes := make([]string, 0, len(byCurrency))
	for code := range byCurrency {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	groups := make([]packetNotConvertedGroup, 0, len(codes))
	for _, code := range codes {
		group := *byCurrency[code]
		group.Total = builder.moneyIn(totals[code], code)
		groups = append(groups, group)
	}
	return groups
}
