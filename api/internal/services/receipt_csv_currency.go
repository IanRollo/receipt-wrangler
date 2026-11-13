package services

import (
	"time"

	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
)

// optionalDecimalString prints a decimal that may be unset: "" when it is.
func optionalDecimalString(value *decimal.Decimal) string {
	if value == nil {
		return ""
	}
	return value.String()
}

// optionalDayString prints a calendar day that may be unset as yyyy-MM-dd: "" when it is.
func optionalDayString(day *models.CalendarDay) string {
	if day == nil || day.IsZero() {
		return ""
	}
	return day.UTC().Format(time.DateOnly)
}
