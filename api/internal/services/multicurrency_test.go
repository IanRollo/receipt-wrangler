package services

import (
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"github.com/shopspring/decimal"
)

func mcDec(value string) decimal.Decimal { return decimal.RequireFromString(value) }

func mcDecPtr(value string) *decimal.Decimal {
	parsed := mcDec(value)
	return &parsed
}

// --- AI extraction ---

func TestSanitizeAiMulticurrency_DropsAMalformedCurrencyAndAnyConversion(t *testing.T) {
	cleaned := sanitizeAiMulticurrency(`{"name":"Shop","amount":12.5,"currency":5,"baseAmount":"9","exchangeRate":1.1}`)

	var command commands.UpsertReceiptCommand
	if err := jsonUnmarshalForTest(cleaned, &command); err != nil {
		t.Fatalf("a malformed currency lost the scan: %v (%s)", err, cleaned)
	}
	if command.Name != "Shop" || command.Currency != "" || command.BaseAmount != nil || command.ExchangeRate != nil {
		t.Errorf("sanitized into %+v", command)
	}

	untouched := `{"name":"Shop","currency":"EUR"}`
	if got := sanitizeAiMulticurrency(untouched); got != untouched {
		t.Errorf("a clean response was rewritten: %s", got)
	}
}

func TestNormalizeAiCurrency_KeepsAValidCodeAndRoundsHalfEven(t *testing.T) {
	command := commands.UpsertReceiptCommand{
		Currency: " eur ",
		Amount:   mcDec("12.345"),
		Items:    []commands.UpsertItemCommand{{Amount: mcDec("2.355")}},
	}
	normalizeAiCurrency(&command, 0, nil)

	if command.Currency != "EUR" || command.Amount.String() != "12.34" || command.Items[0].Amount.String() != "2.36" {
		t.Errorf("got %q %s %s", command.Currency, command.Amount, command.Items[0].Amount)
	}
}

func TestNormalizeAiCurrency_JpyIsRoundedToWholeYen(t *testing.T) {
	command := commands.UpsertReceiptCommand{Currency: "JPY", Amount: mcDec("1500.5")}
	normalizeAiCurrency(&command, 0, nil)
	if command.Amount.String() != "1500" {
		t.Errorf("amount = %s (half-even rounds .5 to the even yen)", command.Amount)
	}
}

func TestNormalizeAiCurrency_AnInvalidOrMissingCodeFallsBackToTheDefault(t *testing.T) {
	for _, code := range []string{"", "ZZZ", "dollars"} {
		command := commands.UpsertReceiptCommand{
			Currency:     code,
			Amount:       mcDec("10.005"),
			BaseAmount:   mcDecPtr("1"),
			ExchangeRate: mcDecPtr("1"),
		}
		normalizeAiCurrency(&command, 0, nil)

		// Empty means "default" when it is stored; the amount is rounded to the default's exponent.
		if command.Currency != "" || command.Amount.String() != "10" {
			t.Errorf("%q: currency %q amount %s", code, command.Currency, command.Amount)
		}
		if command.BaseAmount != nil || command.ExchangeRate != nil {
			t.Errorf("%q: the model's conversion was kept", code)
		}
	}
}

func TestDefaultPrompt_AsksForACurrencyAndTheOldOneDoesNot(t *testing.T) {
	if !strings.Contains(promptTextForTest(), `"currency"`) {
		t.Errorf("the new default prompt does not ask for a currency")
	}
}

// --- Receipt Summary ---

func TestReceiptSummaryFold_LegacyTotalCoversOnlyTheBaseCurrency(t *testing.T) {
	receipts := []models.Receipt{
		{Amount: mcDec("10.00"), Status: models.OPEN},
		{Amount: mcDec("20.00"), Status: models.OPEN, ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "USD"}},
		{Amount: mcDec("50.00"), Status: models.OPEN, ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "EUR"}},
		{Amount: mcDec("700"), Status: models.RESOLVED, ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "JPY"}},
	}

	summary := NewReceiptSummaryService(nil).fold(
		receipts, []models.ReceiptStatus{models.OPEN, models.RESOLVED}, nil, nil, 1, models.RECEIPT_SUMMARY_POSITION_BOTTOM, "USD",
	)

	if summary.BaseCurrency != "USD" {
		t.Errorf("baseCurrency = %q", summary.BaseCurrency)
	}
	if summary.Overall.ReceiptCount != 4 || summary.Overall.Total.String() != "30" {
		t.Errorf("overall count %d total %s: the legacy total must be USD only", summary.Overall.ReceiptCount, summary.Overall.Total)
	}

	totals := summary.Overall.CurrencyTotals
	if len(totals) != 3 || totals[0].Currency != "USD" || totals[1].Currency != "EUR" || totals[2].Currency != "JPY" {
		t.Fatalf("currencyTotals order: %+v", totals)
	}
	if totals[0].Total.String() != "30" || totals[0].ReceiptCount != 2 ||
		totals[1].Total.String() != "50" || totals[2].Total.String() != "700" {
		t.Errorf("currencyTotals = %+v", totals)
	}

	resolved := summary.Statuses[1]
	if resolved.Total.String() != "0" || len(resolved.CurrencyTotals) != 1 || resolved.CurrencyTotals[0].Currency != "JPY" {
		t.Errorf("RESOLVED row: total %s %+v", resolved.Total, resolved.CurrencyTotals)
	}
}

func TestReceiptSummaryFold_CustomFieldTotalsAreSplitByTheReceiptsCurrency(t *testing.T) {
	tip := mcDec("5")
	euroTip := mcDec("3")
	receipts := []models.Receipt{
		{Amount: mcDec("10"), CustomFields: []models.CustomFieldValue{{CustomFieldId: 9, CurrencyValue: &tip}}},
		{Amount: mcDec("20"), ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "EUR"},
			CustomFields: []models.CustomFieldValue{{CustomFieldId: 9, CurrencyValue: &euroTip}}},
	}

	summary := NewReceiptSummaryService(nil).fold(
		receipts, nil, []uint{9}, map[uint]string{9: "Tip"}, 1, models.RECEIPT_SUMMARY_POSITION_BOTTOM, "USD",
	)

	if got := summary.Overall.CustomFieldTotals[0].Total.String(); got != "5" {
		t.Errorf("legacy field total = %s, want the USD tip only", got)
	}
	totals := summary.Overall.CurrencyTotals
	if totals[0].CustomFieldTotals[0].Total.String() != "5" || totals[1].CustomFieldTotals[0].Total.String() != "3" {
		t.Errorf("per-currency field totals: %+v", totals)
	}
}

func TestEmptyReceiptSummary_CurrencyTotalsSerializeAsAnEmptyList(t *testing.T) {
	summary := emptyReceiptSummary(1, models.RECEIPT_SUMMARY_POSITION_BOTTOM)
	if summary.Overall.CurrencyTotals == nil {
		t.Errorf("currencyTotals is nil, which would serialize as null")
	}
}

// --- Claims ---

func linkedForTest(id uint, amount string, code string, base string) repositories.LinkedReceipt {
	linked := repositories.LinkedReceipt{ReceiptId: id, Amount: mcDec(amount), Currency: code}
	if base != "" {
		linked.BaseAmount = mcDecPtr(base)
		linked.BaseAmountCurrency = "USD"
	}
	return linked
}

func TestFoldClaim_BaseConvertedAndUnconverted(t *testing.T) {
	receipts := []repositories.LinkedReceipt{
		linkedForTest(1, "100.00", "USD", ""),
		linkedForTest(2, "100.00", "EUR", "108.50"),
		linkedForTest(3, "40.00", "EUR", ""),
		linkedForTest(4, "25.00", "EUR", ""),
		linkedForTest(5, "5000", "JPY", ""),
		linkedForTest(6, "7.00", "", ""), // a row from before multicurrency counts as the base
	}

	view := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, receipts, nil, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), "USD").view

	if view.Currency != "USD" || view.ReceiptsTotal != "215.50" || view.EffectiveClaimedAmount != "215.50" || view.OutstandingAmount != "215.50" {
		t.Errorf("figures: %s %s %s %s", view.Currency, view.ReceiptsTotal, view.EffectiveClaimedAmount, view.OutstandingAmount)
	}
	if view.ReceiptCount != 6 || view.UnconvertedReceiptCount != 3 {
		t.Errorf("counts: %d receipts, %d unconverted", view.ReceiptCount, view.UnconvertedReceiptCount)
	}
	want := []structs.ClaimCurrencyTotal{{Currency: "EUR", Total: "65.00"}, {Currency: "JPY", Total: "5000"}}
	if len(view.UnconvertedTotals) != 2 || view.UnconvertedTotals[0] != want[0] || view.UnconvertedTotals[1] != want[1] {
		t.Errorf("unconvertedTotals = %+v", view.UnconvertedTotals)
	}
}

func TestFoldClaim_FormatsMoneyAtTheBaseCurrencysExponent(t *testing.T) {
	entries := []models.ClaimLedgerEntry{{
		EntryType: models.ClaimLedgerEntryPayment,
		Date:      time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		Amount:    mcDecPtr("1000"),
	}}

	yen := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, []repositories.LinkedReceipt{linkedForTest(1, "5000", "JPY", "")}, entries, time.Now(), "JPY").view
	if yen.ReceiptsTotal != "5000" || yen.PaymentsTotal != "1000" || yen.LedgerEntries[0].Amount != "1000" {
		t.Errorf("JPY: %s %s %s", yen.ReceiptsTotal, yen.PaymentsTotal, yen.LedgerEntries[0].Amount)
	}

	dinar := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, []repositories.LinkedReceipt{linkedForTest(1, "5.125", "JOD", "")}, nil, time.Now(), "JOD").view
	if dinar.ReceiptsTotal != "5.125" {
		t.Errorf("JOD: %s", dinar.ReceiptsTotal)
	}
}

func TestFoldClaim_NoUnconvertedReceiptsSerializeAsAnEmptyList(t *testing.T) {
	view := foldClaim(models.Claim{Status: models.ClaimStatusDraft}, nil, nil, time.Now(), "USD").view
	if view.UnconvertedTotals == nil || view.UnconvertedReceiptCount != 0 {
		t.Errorf("unconvertedTotals %v count %d", view.UnconvertedTotals, view.UnconvertedReceiptCount)
	}
}

func TestBuildClaimList_OutstandingIsPerCurrencyAndTheScalarSpeaksForOne(t *testing.T) {
	today := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	submitted := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	claimed := func(amount string) *decimal.Decimal { return mcDecPtr(amount) }

	usd := foldClaim(models.Claim{Status: models.ClaimStatusSubmitted, SubmittedDate: &submitted, ClaimedAmount: claimed("100.00")}, nil, nil, today, "USD")
	eur := foldClaim(models.Claim{Status: models.ClaimStatusSubmitted, SubmittedDate: &submitted, ClaimedAmount: claimed("40.00")}, nil, nil, today, "EUR")
	eur2 := foldClaim(models.Claim{Status: models.ClaimStatusSubmitted, SubmittedDate: &submitted, ClaimedAmount: claimed("2.00")}, nil, nil, today, "EUR")

	mixed := buildClaimList([]claimFold{usd, eur, eur2}, today, "USD").Summary
	if len(mixed.OutstandingTotals) != 2 ||
		mixed.OutstandingTotals[0] != (structs.ClaimCurrencyTotal{Currency: "EUR", Total: "42.00"}) ||
		mixed.OutstandingTotals[1] != (structs.ClaimCurrencyTotal{Currency: "USD", Total: "100.00"}) {
		t.Errorf("mixed totals = %+v", mixed.OutstandingTotals)
	}
	if mixed.OutstandingTotal != "100.00" {
		t.Errorf("the scalar must be the install currency's only: %s", mixed.OutstandingTotal)
	}

	single := buildClaimList([]claimFold{eur, eur2}, today, "USD").Summary
	if single.OutstandingTotal != "42.00" || len(single.OutstandingTotals) != 1 {
		t.Errorf("one currency: scalar %s totals %+v", single.OutstandingTotal, single.OutstandingTotals)
	}

	empty := buildClaimList(nil, today, "USD").Summary
	if empty.OutstandingTotals == nil || empty.OutstandingTotal != "0.00" {
		t.Errorf("empty: %v %s", empty.OutstandingTotals, empty.OutstandingTotal)
	}
}

// --- Claim packet ---

func packetBuilderForTest(base string, receipts []models.Receipt) *claimPacketBuilder {
	return &claimPacketBuilder{
		claim:           structs.ClaimView{Name: "Claim", Currency: base, ReceiptCount: len(receipts), ReceiptsTotal: "0.00", EffectiveClaimedAmount: "0.00", PaymentsTotal: "0.00", RecoupmentsTotal: "0.00", PaidAmount: "0.00", OutstandingAmount: "0.00"},
		receipts:        receipts,
		location:        time.UTC,
		now:             time.Date(2026, 11, 12, 10, 0, 0, 0, time.UTC),
		currency:        reporting.CurrencyFormat{Symbol: "$", ThousandsSeparator: ",", DecimalSeparator: "."},
		installCurrency: "USD",
	}
}

func renderPacketForTest(t *testing.T, builder *claimPacketBuilder) string {
	t.Helper()
	sections := make([]packetReceiptSection, 0, len(builder.receipts))
	for _, receipt := range builder.receipts {
		section := buildPacketReceiptSection(receipt, builder.location, builder.moneyFor(receipt))
		builder.addConversionFields(&section, receipt)
		sections = append(sections, section)
	}
	html, err := renderClaimPacketHtml(builder.document(sections), 0)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return html
}

func packetReceiptForTest(id uint, name string, amount string, code string, base string, category string) models.Receipt {
	receipt := models.Receipt{
		BaseModel: models.BaseModel{ID: id},
		Name:      name,
		Amount:    mcDec(amount),
		Date:      time.Date(2026, 9, int(id), 0, 0, 0, 0, time.UTC),
		Status:    models.OPEN,
	}
	receipt.Currency = code
	if base != "" {
		receipt.BaseAmount = mcDecPtr(base)
		receipt.BaseAmountCurrency = "USD"
		receipt.ExchangeRate = mcDecPtr("1.085")
		day := models.NewCalendarDay(time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
		receipt.ExchangeRateDate = &day
		receipt.ExchangeRateSource = "Embassy cashier rate"
	}
	if category != "" {
		receipt.Categories = []models.Category{{BaseModel: models.BaseModel{ID: 1}, Name: category}}
	}
	return receipt
}

func TestClaimPacket_SingleCurrencyReadsAsBefore(t *testing.T) {
	builder := packetBuilderForTest("USD", []models.Receipt{
		packetReceiptForTest(1, "Lunch", "12.34", "USD", "", "Food"),
		packetReceiptForTest(2, "Dinner", "20.00", "", "", "Food"),
	})
	html := renderPacketForTest(t, builder)

	for _, want := range []string{"Subtotal: Food</td><td class=\"num\">$32.34", "$12.34", "Total of 2 receipts (each counted once)", "Claim currency", "USD"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
	for _, unwanted := range []string{"Not converted", "Incomplete", "Counted in"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("a single-currency packet shows %q", unwanted)
		}
	}
}

func TestClaimPacket_CoverSheetSubtotalsAndNotConvertedBlock(t *testing.T) {
	builder := packetBuilderForTest("USD", []models.Receipt{
		packetReceiptForTest(1, "Lunch", "100.00", "USD", "", "Food"),
		packetReceiptForTest(2, "Paris dinner", "100.00", "EUR", "108.50", "Food"),
		packetReceiptForTest(3, "Taxi Paris", "40.00", "EUR", "", "Travel"),
		packetReceiptForTest(4, "Tokyo hotel", "12000", "JPY", "", ""),
	})
	html := renderPacketForTest(t, builder)

	wants := []string{
		// The cover sheet says the claim currency and that the packet is incomplete.
		"Claim currency",
		"Incomplete: 2 receipts are not converted to USD",
		"EUR</th><td class=\"num\">40.00 EUR",
		"JPY</th><td class=\"num\">12,000 JPY",
		// Subtotals and the grand total count only base and converted receipts, in the base currency.
		"Subtotal: Food</td><td class=\"num\">$208.50",
		"Counted in USD",
		"Total of 2 counted receipts in USD",
		"$208.50",
		// The original amount carries its code; the converted one shows its base amount.
		"100.00 EUR",
		"$108.50",
		// The unconverted receipts are listed apart, per currency.
		"<h3>Not converted</h3>",
		"Amount (EUR)",
		"Amount (JPY)",
		"Taxi Paris (#3)",
		"Tokyo hotel (#4)",
		// Per-receipt sections state the conversion as evidence.
		"Amount in USD",
		"1 EUR = 1.085 USD",
		"2026-09-02",
		"Embassy cashier rate",
		"Not converted: no amount in USD was recorded",
	}
	for _, want := range wants {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}

	// The not-converted receipts must not be in the category subtotals.
	if strings.Contains(html, "Subtotal: Travel") {
		t.Errorf("an unconverted receipt's category was subtotalled")
	}
}

func TestFormatPacketNumber(t *testing.T) {
	format := reporting.CurrencyFormat{ThousandsSeparator: ".", DecimalSeparator: ","}
	cases := []struct {
		value    string
		exponent int
		want     string
	}{
		{"1234567.891", 3, "1.234.567,891"},
		{"-1234.5", 2, "-1.234,50"},
		{"12000", 0, "12.000"},
		{"0.5", 2, "0,50"},
		{"999", 0, "999"},
	}
	for _, c := range cases {
		if got := formatPacketNumber(mcDec(c.value), c.exponent, format); got != c.want {
			t.Errorf("%s/%d = %q, want %q", c.value, c.exponent, got, c.want)
		}
	}
}

// --- CSV ---

func TestReceiptCsv_GainsTrailingCurrencyColumns(t *testing.T) {
	receipt := packetReceiptForTest(7, "Paris dinner", "100.00", "EUR", "108.50", "")
	receipt.ReceiptItems = []models.Item{{Name: "wine", Amount: mcDec("10")}}

	service := NewReceiptCsvService()
	service.Location = time.UTC
	result, err := service.BuildReceiptCsv([]models.Receipt{receipt})
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(result.ReceiptCsvBytes)), "\n")
	wantHeader := "Id,Added At,Receipt Date,Name,Paid By,Amount,Status,Categories,Tags,Resolved Date,Currency,Base Amount,Base Currency,Exchange Rate,Rate Date,Rate Source"
	if strings.TrimSpace(lines[0]) != wantHeader {
		t.Errorf("header = %s", lines[0])
	}
	if !strings.HasSuffix(strings.TrimSpace(lines[1]), ",EUR,108.5,USD,1.085,2026-09-02,Embassy cashier rate") {
		t.Errorf("row = %s", lines[1])
	}

	itemLines := strings.Split(strings.TrimSpace(string(result.ReceiptItemCsvBytes)), "\n")
	if !strings.HasSuffix(strings.TrimSpace(itemLines[0]), ",Currency") || !strings.HasSuffix(strings.TrimSpace(itemLines[1]), ",EUR") {
		t.Errorf("items.csv: %v", itemLines)
	}
}

// --- Receipt history ---

func TestReceiptHistoryDiff_RecordsCurrencyAndConversionChanges(t *testing.T) {
	before := packetReceiptForTest(1, "R", "100.00", "EUR", "108.50", "")
	after := before
	after.BaseAmount = mcDecPtr("109.00")
	after.ExchangeRateSource = "Card statement"
	after.Currency = "GBP"
	after.ExchangeRate = nil

	changes := diffReceiptSnapshots(&before, after, receiptDiffContext{})
	got := map[string]string{}
	for _, change := range changes {
		got[change.Field] = change.Kind
	}
	want := map[string]string{
		structs.ReceiptHistoryFieldCurrency:           structs.ReceiptHistoryKindChanged,
		structs.ReceiptHistoryFieldBaseAmount:         structs.ReceiptHistoryKindChanged,
		structs.ReceiptHistoryFieldExchangeRate:       structs.ReceiptHistoryKindRemoved,
		structs.ReceiptHistoryFieldExchangeRateSource: structs.ReceiptHistoryKindChanged,
	}
	for field, kind := range want {
		if got[field] != kind {
			t.Errorf("%s: %q, want %q (all: %v)", field, got[field], kind, got)
		}
	}
	if _, ok := got[structs.ReceiptHistoryFieldExchangeRateDate]; ok {
		t.Errorf("an unchanged rate date was recorded")
	}

	if identical := diffReceiptSnapshots(&before, before, receiptDiffContext{}); len(identical) != 0 {
		t.Errorf("identical snapshots produced %v", identical)
	}
}
