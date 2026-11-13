package services

import (
	"time"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

// currencyBlock records changes to a receipt's currency and its recorded conversion
// (multicurrency, fork feature), right after the amount in the fixed field order. A
// conversion is evidence, so a change to any part of it is a history entry like any other.
func (diff *receiptDiff) currencyBlock(from models.Receipt, after models.Receipt) {
	diff.compare(structs.ReceiptHistoryFieldCurrency, "", textSide(from.Currency), textSide(after.Currency))
	diff.compare(structs.ReceiptHistoryFieldBaseAmount, "", currencySide(from.BaseAmount), currencySide(after.BaseAmount))
	diff.compare(structs.ReceiptHistoryFieldExchangeRate, "", textSide(optionalDecimalString(from.ExchangeRate)), textSide(optionalDecimalString(after.ExchangeRate)))
	diff.compare(structs.ReceiptHistoryFieldExchangeRateDate, "", calendarDaySide(dayAsTime(from.ExchangeRateDate)), calendarDaySide(dayAsTime(after.ExchangeRateDate)))
	diff.compare(structs.ReceiptHistoryFieldExchangeRateSource, "", textSide(from.ExchangeRateSource), textSide(after.ExchangeRateSource))
}

// dayAsTime views an optional calendar day as the *time.Time calendarDaySide compares.
func dayAsTime(day *models.CalendarDay) *time.Time {
	if day == nil {
		return nil
	}
	return &day.Time
}
