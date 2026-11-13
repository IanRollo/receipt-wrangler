package models

import (
	"receipt-wrangler/api/internal/currency"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ReceiptMulticurrency is the fork's currency block of a receipt, embedded in Receipt so
// the upstream-hot model file carries one line. GORM and encoding/json both flatten it:
// the columns live on receipts and the keys sit at the top level of the receipt.
//
// Currency is an ISO 4217 code, plain string on the wire (never an enum), and is never
// empty once a receipt is stored: every create and update resolves it (request, then the
// group's default, then the install default). Items have no currency of their own and
// inherit this one.
//
// The conversion fields are the user's record of what they were actually charged, kept as
// evidence. They are meaningful only when Currency differs from the group's base currency
// (GroupReceiptSettings.DefaultCurrency, resolved) and are cleared when it does not:
//   - BaseAmount is in the base currency;
//   - ExchangeRate is base units per one unit of Currency;
//   - ExchangeRateDate is a calendar day;
//   - ExchangeRateSource is free text (e.g. "Card statement").
//
// ConversionNote is a transient, read-only warning (gorm:"-") raised when a stored base
// amount and rate disagree about the amount by more than 1%.
type ReceiptMulticurrency struct {
	Currency           string           `gorm:"type:varchar(3)" json:"currency"`
	BaseAmount         *decimal.Decimal `gorm:"type:decimal(13,3)" json:"baseAmount,omitempty"`
	ExchangeRate       *decimal.Decimal `gorm:"type:decimal(20,6)" json:"exchangeRate,omitempty"`
	ExchangeRateDate   *CalendarDay     `json:"exchangeRateDate,omitempty"`
	ExchangeRateSource string           `gorm:"size:255" json:"exchangeRateSource,omitempty"`
	ConversionNote     string           `gorm:"-" json:"conversionNote,omitempty"`
}

// HasConversion reports whether the receipt carries a recorded conversion.
func (receipt Receipt) HasConversion() bool {
	return receipt.BaseAmount != nil || receipt.ExchangeRate != nil
}

// FillConversionNote recomputes the transient conversionNote from the stored figures.
func (receipt *Receipt) FillConversionNote() {
	receipt.ConversionNote = currency.ConversionNote(receipt.Amount, receipt.BaseAmount, receipt.ExchangeRate)
}

// AfterFind fills the transient conversionNote on every read path (single, list, preload,
// history snapshots), since it is a pure function of columns already loaded.
func (receipt *Receipt) AfterFind(_ *gorm.DB) error {
	receipt.FillConversionNote()
	return nil
}
