// Package currency is the fork's multicurrency core: the table of active ISO 4217
// codes with their minor-unit exponents, amount validation by exponent, and the
// derivation of a recorded conversion. It has no database or HTTP dependency, so
// commands, repositories and services can all import it.
package currency

import (
	"sort"
	"strings"
)

// DefaultInstallCurrency is the install default currency before an administrator
// chooses one, and the value the multicurrency backfill labels existing receipts
// with on a fresh settings row.
const DefaultInstallCurrency = "USD"

// Info describes one currency. Exponent is ISO 4217's minor-unit exponent: the
// number of decimal places an amount in this currency may have (JPY 0, USD 2,
// JOD 3). Symbol is set only where the symbol identifies the currency; the dollar
// and yen signs are shared by several currencies and so are left to the code,
// except USD, the install's historical default.
type Info struct {
	Code     string
	Exponent int
	Symbol   string
}

var infos = buildTable()

func buildTable() map[string]Info {
	table := make(map[string]Info)
	add := func(exponent int, codes ...string) {
		for _, code := range codes {
			table[code] = Info{Code: code, Exponent: exponent}
		}
	}

	add(0, "BIF", "CLP", "DJF", "GNF", "ISK", "JPY", "KMF", "KRW", "PYG", "RWF", "UGX", "VND", "VUV", "XAF", "XOF", "XPF")
	add(3, "BHD", "IQD", "JOD", "KWD", "LYD", "OMR", "TND")
	// BGN is kept although Bulgaria joined the euro in 2026: a post's older receipts are still
	// in leva, and refusing the code would make them uneditable.
	add(2,
		"AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN",
		"BAM", "BBD", "BDT", "BGN", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BZD",
		"CAD", "CDF", "CHF", "CNY", "COP", "CRC", "CUP", "CVE", "CZK",
		"DKK", "DOP", "DZD",
		"EGP", "ERN", "ETB", "EUR",
		"FJD", "FKP",
		"GBP", "GEL", "GHS", "GIP", "GMD", "GTQ", "GYD",
		"HKD", "HNL", "HTG", "HUF",
		"IDR", "ILS", "INR", "IRR",
		"JMD",
		"KES", "KGS", "KHR", "KPW", "KYD", "KZT",
		"LAK", "LBP", "LKR", "LRD", "LSL",
		"MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN",
		"NAD", "NGN", "NIO", "NOK", "NPR", "NZD",
		"PAB", "PEN", "PGK", "PHP", "PKR", "PLN",
		"QAR",
		"RON", "RSD", "RUB",
		"SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLE", "SOS", "SRD", "SSP", "STN", "SVC", "SYP", "SZL",
		"THB", "TJS", "TMT", "TOP", "TRY", "TTD", "TWD", "TZS",
		"UAH", "USD", "UYU", "UZS",
		"VES",
		"WST",
		"XCD", "XCG",
		"YER",
		"ZAR", "ZMW", "ZWG",
	)

	symbols := map[string]string{
		"USD": "$", "EUR": "€", "GBP": "£", "INR": "₹", "KRW": "₩", "ILS": "₪", "VND": "₫",
		"PHP": "₱", "THB": "฿", "UAH": "₴", "NGN": "₦", "TRY": "₺", "KZT": "₸", "GEL": "₾",
		"AZN": "₼", "RUB": "₽", "PLN": "zł", "CZK": "Kč",
	}
	for code, symbol := range symbols {
		info := table[code]
		info.Symbol = symbol
		table[code] = info
	}

	return table
}

// Normalize trims and upper-cases a code, the form it is stored and compared in.
func Normalize(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

// Lookup returns the currency a code names, after Normalize.
func Lookup(code string) (Info, bool) {
	info, ok := infos[Normalize(code)]
	return info, ok
}

// IsValid reports whether code (after Normalize) is an active ISO 4217 code in the table.
func IsValid(code string) bool {
	_, ok := Lookup(code)
	return ok
}

// ExponentOf returns the number of decimal places a currency allows. An unknown code gets 2,
// the exponent of the install's historical single currency; callers validate codes first.
func ExponentOf(code string) int {
	if info, ok := Lookup(code); ok {
		return info.Exponent
	}
	return 2
}

// All returns every currency in the table, ordered by code.
func All() []Info {
	result := make([]Info, 0, len(infos))
	for _, info := range infos {
		result = append(result, info)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}
