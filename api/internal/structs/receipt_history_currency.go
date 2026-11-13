package structs

// Receipt history field codes added by multicurrency (fork feature). Like every enum-like
// value on the history contract they are plain strings on the wire; a client must tolerate a
// code it does not know.
const (
	ReceiptHistoryFieldCurrency           = "CURRENCY"
	ReceiptHistoryFieldBaseAmount         = "BASE_AMOUNT"
	ReceiptHistoryFieldExchangeRate       = "EXCHANGE_RATE"
	ReceiptHistoryFieldExchangeRateDate   = "EXCHANGE_RATE_DATE"
	ReceiptHistoryFieldExchangeRateSource = "EXCHANGE_RATE_SOURCE"
)
