package currency

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func dec(value string) decimal.Decimal { return decimal.RequireFromString(value) }

func decPtr(value string) *decimal.Decimal {
	parsed := dec(value)
	return &parsed
}

func TestLookup_NormalizesAndKnowsExponents(t *testing.T) {
	cases := map[string]int{"usd": 2, " EUR ": 2, "JPY": 0, "krw": 0, "JOD": 3, "KWD": 3, "BHD": 3, "OMR": 3, "TND": 3}
	for code, exponent := range cases {
		info, ok := Lookup(code)
		if !ok {
			t.Fatalf("%q is not in the table", code)
		}
		if info.Exponent != exponent {
			t.Errorf("%q exponent = %d, want %d", code, info.Exponent, exponent)
		}
		if info.Code != Normalize(code) {
			t.Errorf("%q code = %q", code, info.Code)
		}
	}
}

func TestIsValid_RejectsUnknownCodes(t *testing.T) {
	for _, code := range []string{"", "US", "USDD", "XXX", "ZZZ", "$", "1234", "BTC", "HRK", "ZWL"} {
		if IsValid(code) {
			t.Errorf("%q should not be a valid code", code)
		}
	}
	for _, code := range []string{"USD", "eur", "GBP", "JPY", "JOD"} {
		if !IsValid(code) {
			t.Errorf("%q should be valid", code)
		}
	}
}

func TestTable_EveryCodeIsThreeUpperLetters(t *testing.T) {
	all := All()
	if len(all) < 150 {
		t.Errorf("table has only %d currencies", len(all))
	}
	for _, info := range all {
		if len(info.Code) != 3 || info.Code != strings.ToUpper(info.Code) {
			t.Errorf("bad code %q", info.Code)
		}
		if info.Exponent < 0 || info.Exponent > 3 {
			t.Errorf("%s exponent %d is outside what the amount columns hold", info.Code, info.Exponent)
		}
	}
}

func TestValidateAmount_ByExponent(t *testing.T) {
	cases := []struct {
		amount string
		code   string
		ok     bool
	}{
		{"1500", "JPY", true},
		{"1500.0", "JPY", true},
		{"1500.5", "JPY", false},
		{"12.34", "USD", true},
		{"12.3", "USD", true},
		{"12.345", "USD", false},
		{"12.345", "JOD", true},
		{"12.3456", "JOD", false},
		{"-5.25", "EUR", true},
		{"-5.255", "EUR", false},
		{"99999999999", "USD", false},
	}
	for _, c := range cases {
		message := ValidateAmount(dec(c.amount), c.code)
		if (message == "") != c.ok {
			t.Errorf("%s %s: message %q, want ok=%v", c.amount, c.code, message, c.ok)
		}
	}

	if got := ValidateAmount(dec("1.5"), "JPY"); !strings.Contains(got, "whole") {
		t.Errorf("JPY message = %q", got)
	}
}

func TestRoundToExponent_IsHalfEven(t *testing.T) {
	cases := []struct {
		in, code, want string
	}{
		{"2.345", "USD", "2.34"},
		{"2.355", "USD", "2.36"},
		{"0.5", "JPY", "0"},
		{"1.5", "JPY", "2"},
		{"2.5", "JPY", "2"},
		{"1.2345", "JOD", "1.234"},
	}
	for _, c := range cases {
		if got := RoundToExponent(dec(c.in), c.code).String(); got != c.want {
			t.Errorf("round(%s, %s) = %s, want %s", c.in, c.code, got, c.want)
		}
	}
}

func TestDeriveConversion(t *testing.T) {
	amount := dec("100.00")

	base, rate := DeriveConversion(amount, "USD", decPtr("108.50"), nil)
	if base.String() != "108.5" || rate == nil || rate.String() != "1.085" {
		t.Errorf("from base: %v %v", base, rate)
	}

	base, rate = DeriveConversion(amount, "USD", nil, decPtr("1.0855"))
	if base == nil || base.String() != "108.55" || rate.String() != "1.0855" {
		t.Errorf("from rate: %v %v", base, rate)
	}

	// Half-even to the base currency exponent: 33.325 -> 33.32 (the even neighbour).
	base, _ = DeriveConversion(dec("33.325"), "USD", nil, decPtr("1"))
	if base.String() != "33.32" {
		t.Errorf("half-even: %v", base)
	}

	// A whole-yen base rounds to yen.
	base, _ = DeriveConversion(dec("12.34"), "JPY", nil, decPtr("151.5"))
	if base.String() != "1870" {
		t.Errorf("JPY base: %v", base)
	}

	// Both given: kept exactly, even though they disagree.
	base, rate = DeriveConversion(amount, "USD", decPtr("120"), decPtr("1.0855"))
	if base.String() != "120" || rate.String() != "1.0855" {
		t.Errorf("both given: %v %v", base, rate)
	}

	base, rate = DeriveConversion(amount, "USD", nil, nil)
	if base != nil || rate != nil {
		t.Errorf("neither: %v %v", base, rate)
	}

	// A zero amount has no rate to derive.
	base, rate = DeriveConversion(decimal.Zero, "USD", decPtr("5"), nil)
	if base == nil || rate != nil {
		t.Errorf("zero amount: %v %v", base, rate)
	}

	// The rate is kept to six places.
	_, rate = DeriveConversion(dec("3"), "USD", decPtr("10"), nil)
	if rate.String() != "3.333333" {
		t.Errorf("rate precision: %v", rate)
	}
}

func TestConversionNote(t *testing.T) {
	amount := dec("100")

	if note := ConversionNote(amount, decPtr("108.5"), decPtr("1.085")); note != "" {
		t.Errorf("agreeing figures raised %q", note)
	}
	// A card fee of under 1% is normal.
	if note := ConversionNote(amount, decPtr("109.3"), decPtr("1.085")); note != "" {
		t.Errorf("0.7%% gap raised %q", note)
	}
	if note := ConversionNote(amount, decPtr("120"), decPtr("1.085")); note == "" {
		t.Errorf("a 10%% gap raised nothing")
	}
	if note := ConversionNote(amount, nil, decPtr("1.085")); note != "" {
		t.Errorf("a missing base amount raised %q", note)
	}
	if note := ConversionNote(amount, decPtr("108.5"), nil); note != "" {
		t.Errorf("a missing rate raised %q", note)
	}
}
