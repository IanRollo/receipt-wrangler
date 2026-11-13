package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/shopspring/decimal"
)

// Multicurrency (fork feature): the receipt write path as a client sees it.

func receiptBody(groupId uint, userId uint, extra string) string {
	body := fmt.Sprintf(
		`{"name":"R","amount":"100.00","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"`,
		groupId, userId,
	)
	if extra != "" {
		body += "," + extra
	}
	return body + "}"
}

func postReceipt(userId uint, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(body))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(userId)))
	CreateReceipt(w, r)
	return w
}

func decodeReceipt(t *testing.T, w *httptest.ResponseRecorder) models.Receipt {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", w.Code, w.Body.String())
	}
	var receipt models.Receipt
	if err := json.Unmarshal(w.Body.Bytes(), &receipt); err != nil {
		t.Fatalf("decode receipt %q: %v", w.Body.String(), err)
	}
	return receipt
}

func validatorKeys(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var errors map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errors); err != nil {
		t.Fatalf("decode validator error %q: %v", w.Body.String(), err)
	}
	return errors
}

func setGroupDefaultCurrency(t *testing.T, groupId uint, code string) {
	t.Helper()
	db := repositories.GetDB()
	var settings models.GroupReceiptSettings
	if err := db.Where("group_id = ?", groupId).Limit(1).Find(&settings).Error; err != nil {
		t.Fatalf("load group settings: %v", err)
	}
	if settings.ID == 0 {
		settings = models.GroupReceiptSettings{GroupId: groupId, DefaultCurrency: code}
		if err := db.Create(&settings).Error; err != nil {
			t.Fatalf("create group settings: %v", err)
		}
		return
	}
	if err := db.Model(&settings).Update("default_currency", code).Error; err != nil {
		t.Fatalf("set group currency: %v", err)
	}
}

func storedReceipt(t *testing.T, id uint) models.Receipt {
	t.Helper()
	var receipt models.Receipt
	if err := repositories.GetDB().First(&receipt, id).Error; err != nil {
		t.Fatalf("load receipt: %v", err)
	}
	return receipt
}

func TestMulticurrency_OldClientCreateWithoutCurrencyGetsTheInstallDefault(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	receipt := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, "")))

	if receipt.Currency != "USD" {
		t.Errorf("currency = %q, want the install default USD", receipt.Currency)
	}
	if stored := storedReceipt(t, receipt.ID); stored.Currency != "USD" || stored.BaseAmount != nil || stored.ExchangeRate != nil {
		t.Errorf("stored %q base %v rate %v", stored.Currency, stored.BaseAmount, stored.ExchangeRate)
	}
}

func TestMulticurrency_ResolutionOrder_RequestThenGroupThenInstall(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	// Install default changed: a receipt with no currency takes it.
	repositories.GetDB().Create(&models.SystemSettings{DefaultCurrency: "GBP"})
	if got := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, ""))).Currency; got != "GBP" {
		t.Errorf("install default: %q", got)
	}

	// The group's default beats the install's.
	setGroupDefaultCurrency(t, groupId, "EUR")
	if got := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, ""))).Currency; got != "EUR" {
		t.Errorf("group default: %q", got)
	}

	// The request beats both, and the code is normalized.
	if got := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, `"currency":"jpy","amount":"1500"`))).Currency; got != "JPY" {
		t.Errorf("request currency: %q", got)
	}

	// An empty currency means "default", not an error.
	if got := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, `"currency":""`))).Currency; got != "EUR" {
		t.Errorf("empty currency: %q", got)
	}
}

func TestMulticurrency_UnknownCurrencyAndTooManyDecimalsAre400s(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	w := postReceipt(userId, receiptBody(groupId, userId, `"currency":"ZZZ"`))
	if w.Code != http.StatusBadRequest || validatorKeys(t, w)["currency"] == "" {
		t.Errorf("unknown code: %d %s", w.Code, w.Body.String())
	}

	// JPY is whole.
	w = postReceipt(userId, receiptBody(groupId, userId, `"currency":"JPY","amount":"1500.5"`))
	if w.Code != http.StatusBadRequest || validatorKeys(t, w)["amount"] == "" {
		t.Errorf("JPY decimals: %d %s", w.Code, w.Body.String())
	}

	// USD at most two places, and never silently rounded.
	w = postReceipt(userId, receiptBody(groupId, userId, `"amount":"10.005"`))
	if w.Code != http.StatusBadRequest || validatorKeys(t, w)["amount"] == "" {
		t.Errorf("USD decimals: %d %s", w.Code, w.Body.String())
	}

	// JOD takes three, and its items follow the receipt's currency.
	w = postReceipt(userId, receiptBody(groupId, userId, `"currency":"JOD","amount":"10.005"`))
	if w.Code != http.StatusOK {
		t.Errorf("JOD three places: %d %s", w.Code, w.Body.String())
	}
	w = postReceipt(userId, receiptBody(groupId, userId,
		`"amount":"10.00","receiptItems":[{"name":"i","amount":"1.005","status":"OPEN"}]`))
	if w.Code != http.StatusBadRequest || validatorKeys(t, w)["receiptItems.0.amount"] == "" {
		t.Errorf("USD item decimals: %d %s", w.Code, w.Body.String())
	}

	if count := receiptCount(t); count != 1 {
		t.Errorf("%d receipts stored, want only the JOD one", count)
	}
}

func receiptCount(t *testing.T) int64 {
	t.Helper()
	var count int64
	repositories.GetDB().Model(&models.Receipt{}).Count(&count)
	return count
}

func TestMulticurrency_ThreePlaceAmountRoundTrips(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	receipt := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, `"currency":"KWD","amount":"12.345"`)))
	if got := storedReceipt(t, receipt.ID).Amount.String(); got != "12.345" {
		t.Errorf("stored amount = %s", got)
	}
}

func TestMulticurrency_ConversionIsDerivedBothWays(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	// Only the amount charged: the rate is derived.
	receipt := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","baseAmount":"108.50","exchangeRateDate":"2026-09-01","exchangeRateSource":"Card statement"`)))
	if receipt.BaseAmount == nil || receipt.BaseAmount.String() != "108.5" ||
		receipt.ExchangeRate == nil || receipt.ExchangeRate.String() != "1.085" {
		t.Errorf("from base: %v %v", receipt.BaseAmount, receipt.ExchangeRate)
	}
	stored := storedReceipt(t, receipt.ID)
	if stored.ExchangeRateDate == nil || stored.ExchangeRateDate.UTC().Format(time.DateOnly) != "2026-09-01" || stored.ExchangeRateSource != "Card statement" {
		t.Errorf("stored date/source: %v %q", stored.ExchangeRateDate, stored.ExchangeRateSource)
	}
	wire, _ := json.Marshal(stored)
	if !strings.Contains(string(wire), `"exchangeRateDate":"2026-09-01"`) {
		t.Errorf("the rate date is not a bare day on the wire: %s", wire)
	}

	// Only the rate: the base amount is derived and rounded to the base exponent.
	receipt = decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","amount":"33.33","exchangeRate":"1.0855"`)))
	if receipt.BaseAmount == nil || receipt.BaseAmount.String() != "36.18" {
		t.Errorf("from rate: %v", receipt.BaseAmount)
	}
	if receipt.ConversionNote != "" {
		t.Errorf("derived figures raised a note: %q", receipt.ConversionNote)
	}
}

func TestMulticurrency_BothGivenAreKeptAndDisagreementRaisesANote(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	receipt := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","baseAmount":"120.00","exchangeRate":"1.085"`)))
	if receipt.BaseAmount.String() != "120" || receipt.ExchangeRate.String() != "1.085" {
		t.Errorf("figures were altered: %v %v", receipt.BaseAmount, receipt.ExchangeRate)
	}

	// The note is computed on read, from the stored figures.
	reloaded := storedReceipt(t, receipt.ID)
	if reloaded.ConversionNote == "" {
		t.Errorf("a 10%% disagreement raised no conversionNote")
	}

	// A fee-sized gap does not.
	receipt = decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","baseAmount":"109.00","exchangeRate":"1.085"`)))
	reloaded = storedReceipt(t, receipt.ID)
	if reloaded.ConversionNote != "" {
		t.Errorf("a 0.5%% gap raised %q", reloaded.ConversionNote)
	}
}

func TestMulticurrency_ConversionIsClearedWhenTheCurrencyIsTheBase(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	receipt := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"USD","baseAmount":"100.00","exchangeRate":"1","exchangeRateSource":"x"`)))
	stored := storedReceipt(t, receipt.ID)
	if stored.BaseAmount != nil || stored.ExchangeRate != nil || stored.ExchangeRateDate != nil || stored.ExchangeRateSource != "" {
		t.Errorf("a USD receipt in a USD group kept a conversion: %+v", stored.ReceiptMulticurrency)
	}
}

func TestMulticurrency_OldClientUpdateKeepsTheStoredCurrencyAndConversion(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	created := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","baseAmount":"108.50","exchangeRateDate":"2026-09-01","exchangeRateSource":"Statement"`)))

	// The old client edits the name and knows nothing about currency or conversion.
	body := fmt.Sprintf(`{"name":"Renamed","amount":"100.00","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`, groupId, userId)
	w, r := updateReceiptRequest(created.ID, userId, body)
	UpdateReceipt(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("update: %d %s", w.Code, w.Body.String())
	}

	stored := storedReceipt(t, created.ID)
	if stored.Name != "Renamed" || stored.Currency != "EUR" {
		t.Errorf("name %q currency %q: a EUR receipt must stay EUR", stored.Name, stored.Currency)
	}
	if stored.BaseAmount == nil || stored.BaseAmount.String() != "108.5" || stored.ExchangeRate == nil ||
		stored.ExchangeRateDate == nil || stored.ExchangeRateSource != "Statement" {
		t.Errorf("the conversion did not survive: %+v", stored.ReceiptMulticurrency)
	}
}

func TestMulticurrency_UpdateReplacesOrClearsTheConversion(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	created := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId,
		`"currency":"EUR","baseAmount":"108.50","exchangeRateSource":"Statement"`)))

	// Sending any conversion key replaces the whole conversion.
	body := fmt.Sprintf(`{"name":"R","amount":"100.00","currency":"EUR","exchangeRate":"1.1","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`, groupId, userId)
	w, r := updateReceiptRequest(created.ID, userId, body)
	UpdateReceipt(w, r)
	stored := storedReceipt(t, created.ID)
	if w.Code != http.StatusOK || stored.BaseAmount.String() != "110" || stored.ExchangeRate.String() != "1.1" || stored.ExchangeRateSource != "" {
		t.Errorf("replace: %d %+v", w.Code, stored.ReceiptMulticurrency)
	}

	// Moving the receipt into the base currency clears it, explicitly.
	body = fmt.Sprintf(`{"name":"R","amount":"100.00","currency":"USD","date":"2026-09-01T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"OPEN"}`, groupId, userId)
	w, r = updateReceiptRequest(created.ID, userId, body)
	UpdateReceipt(w, r)
	stored = storedReceipt(t, created.ID)
	if w.Code != http.StatusOK || stored.Currency != "USD" || stored.BaseAmount != nil || stored.ExchangeRate != nil {
		t.Errorf("clear: %d %+v", w.Code, stored.ReceiptMulticurrency)
	}
}

func TestMulticurrency_OldShapedRequestsStillCreateAndUpdateAsBefore(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, groupId := seedRestrictedReceiptCreator(t, nil)

	created := decodeReceipt(t, postReceipt(userId, receiptBody(groupId, userId, "")))
	if created.Amount.String() != "100" || created.Name != "R" || created.Status != models.OPEN {
		t.Errorf("created %+v", created)
	}

	body := fmt.Sprintf(`{"name":"R2","amount":"55.50","date":"2026-09-02T00:00:00Z","groupId":%d,"paidByUserId":%d,"status":"RESOLVED"}`, groupId, userId)
	w, r := updateReceiptRequest(created.ID, userId, body)
	UpdateReceipt(w, r)
	stored := storedReceipt(t, created.ID)
	if w.Code != http.StatusOK || stored.Name != "R2" || stored.Amount.String() != "55.5" || stored.Currency != "USD" {
		t.Errorf("updated: %d %+v", w.Code, stored)
	}

	// The wire still carries every old key.
	wire, _ := json.Marshal(created)
	for _, key := range []string{`"name"`, `"amount"`, `"date"`, `"groupId"`, `"status"`, `"receiptItems"`, `"customFields"`} {
		if !strings.Contains(string(wire), key) {
			t.Errorf("response lost %s", key)
		}
	}
}

func TestMulticurrency_CurrencyTableEndpoint(t *testing.T) {
	defer repositories.TruncateTestDb()
	userId, _ := seedRestrictedReceiptCreator(t, nil)

	w, r := requestForUser(userId)
	GetCurrencies(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var table []struct {
		Code     string `json:"code"`
		Exponent int    `json:"exponent"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &table); err != nil {
		t.Fatal(err)
	}
	exponents := map[string]int{}
	for _, entry := range table {
		exponents[entry.Code] = entry.Exponent
	}
	if exponents["JPY"] != 0 || exponents["USD"] != 2 || exponents["JOD"] != 3 {
		t.Errorf("exponents JPY %d USD %d JOD %d", exponents["JPY"], exponents["USD"], exponents["JOD"])
	}
}

func seedForeignReceipt(t *testing.T, groupId uint, userId uint, amount string, code string, baseAmount string) uint {
	t.Helper()
	receipt := models.Receipt{
		Name:         "Foreign",
		Amount:       decimal.RequireFromString(amount),
		Date:         time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		GroupId:      groupId,
		PaidByUserID: userId,
		Status:       models.OPEN,
	}
	receipt.Currency = code
	if baseAmount != "" {
		base := decimal.RequireFromString(baseAmount)
		receipt.BaseAmount = &base
	}
	if err := repositories.GetDB().Create(&receipt).Error; err != nil {
		t.Fatalf("seed receipt: %v", err)
	}
	return receipt.ID
}

func TestMulticurrency_ClaimFiguresCountBaseConvertedAndSetAsideUnconverted(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	claim := seedClaimRow(t, 1, 1, "PCS")
	base := seedForeignReceipt(t, 1, 1, "100.00", "USD", "")
	converted := seedForeignReceipt(t, 1, 1, "100.00", "EUR", "108.50")
	unconvertedEur := seedForeignReceipt(t, 1, 1, "40.00", "EUR", "")
	unconvertedJpy := seedForeignReceipt(t, 1, 1, "5000", "JPY", "")
	linkReceipts(t, claim.ID, base, converted, unconvertedEur, unconvertedJpy)

	view := getClaim(t, 1, claim.ID)
	if view.Currency != "USD" {
		t.Errorf("claim currency = %q", view.Currency)
	}
	if view.ReceiptsTotal != "208.50" || view.OutstandingAmount != "208.50" {
		t.Errorf("receiptsTotal %s outstanding %s, want 100.00 + 108.50", view.ReceiptsTotal, view.OutstandingAmount)
	}
	if view.ReceiptCount != 4 || view.UnconvertedReceiptCount != 2 {
		t.Errorf("count %d unconverted %d", view.ReceiptCount, view.UnconvertedReceiptCount)
	}
	if len(view.UnconvertedTotals) != 2 || view.UnconvertedTotals[0] != (structs.ClaimCurrencyTotal{Currency: "EUR", Total: "40.00"}) ||
		view.UnconvertedTotals[1] != (structs.ClaimCurrencyTotal{Currency: "JPY", Total: "5000"}) {
		t.Errorf("unconvertedTotals = %+v", view.UnconvertedTotals)
	}
}

func TestMulticurrency_ClaimInAJpyGroupIsReckonedInWholeYen(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)
	setGroupDefaultCurrency(t, 1, "JPY")

	claim := seedClaimRow(t, 1, 1, "Tokyo")
	yen := seedForeignReceipt(t, 1, 1, "12000", "JPY", "")
	usd := seedForeignReceipt(t, 1, 1, "10.00", "USD", "1500")
	linkReceipts(t, claim.ID, yen, usd)

	view := getClaim(t, 1, claim.ID)
	if view.Currency != "JPY" || view.ReceiptsTotal != "13500" {
		t.Errorf("currency %s total %s", view.Currency, view.ReceiptsTotal)
	}

	// The ledger is validated by the base currency's exponent: no fractions of a yen.
	w, r := claimRequest(http.MethodPost, 1, `{"type":"PAYMENT","date":"2026-10-01","amount":"100.50"}`, claimParams(claim.ID))
	CreateClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusBadRequest)

	w, r = claimRequest(http.MethodPost, 1, `{"type":"PAYMENT","date":"2026-10-01","amount":"1000"}`, claimParams(claim.ID))
	CreateClaimLedgerEntry(w, r)
	assertStatus(t, w, http.StatusOK)
	view = getClaim(t, 1, claim.ID)
	if view.PaymentsTotal != "1000" || view.OutstandingAmount != "12500" {
		t.Errorf("payments %s outstanding %s", view.PaymentsTotal, view.OutstandingAmount)
	}
}

func TestMulticurrency_AllGroupClaimSummaryIsPerCurrency(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)

	// A second real group in EUR.
	second := models.Group{Name: "Second"}
	repositories.GetDB().Create(&second)
	grantGroupPerms(t, 1, second.ID, permissions.GroupReceiptsRead, permissions.GroupReceiptsUpdate)
	setGroupDefaultCurrency(t, second.ID, "EUR")

	submitted := func(groupId uint, name string, claimed string) {
		claim := seedClaimRow(t, groupId, 1, name)
		amount := decimal.RequireFromString(claimed)
		day := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		repositories.GetDB().Model(&claim).Updates(map[string]interface{}{
			"status": models.ClaimStatusSubmitted, "claimed_amount": amount, "submitted_date": day,
		})
	}
	submitted(1, "usd claim", "100.00")
	submitted(second.ID, "eur claim", "40.00")

	allGroup := models.Group{Name: "All", IsAllGroup: true}
	repositories.GetDB().Create(&allGroup)
	grantGroupPerms(t, 1, allGroup.ID, permissions.GroupReceiptsRead)

	list := listClaims(t, 1, allGroup.ID)
	summary := list.Summary
	if len(summary.OutstandingTotals) != 2 ||
		summary.OutstandingTotals[0] != (structs.ClaimCurrencyTotal{Currency: "EUR", Total: "40.00"}) ||
		summary.OutstandingTotals[1] != (structs.ClaimCurrencyTotal{Currency: "USD", Total: "100.00"}) {
		t.Errorf("outstandingTotals = %+v", summary.OutstandingTotals)
	}
	// The legacy scalar only ever speaks for the install default currency here.
	if summary.OutstandingTotal != "100.00" {
		t.Errorf("outstandingTotal = %s, want the USD total only", summary.OutstandingTotal)
	}
}

func TestMulticurrency_ClaimPacketEndToEnd(t *testing.T) {
	defer tearDownGenericHandlerTest()
	stubPacketPdf(t)
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)

	claim := seedClaimRow(t, 1, 1, "Paris trip")
	usd := seedForeignReceipt(t, 1, 1, "50.00", "USD", "")
	converted := seedForeignReceipt(t, 1, 1, "100.00", "EUR", "108.50")
	unconverted := seedForeignReceipt(t, 1, 1, "40.00", "EUR", "")
	linkReceipts(t, claim.ID, usd, converted, unconverted)

	packet := downloadPacket(t, 1, claim.ID)
	html := packet.html()
	for _, want := range []string{
		"Claim currency",
		"Incomplete: 1 receipt is not converted to USD",
		"40.00 EUR",
		"Counted in USD",
		"$158.50", // 50.00 + 108.50, the grand total
		"<h3>Not converted</h3>",
		"Amount in USD",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the PDF does not print %q", want)
		}
	}

	rows := readCsv(t, packet.files["receipts.csv"])
	header := strings.Join(rows[0], ",")
	if !strings.Contains(header, "Currency,Base Amount,Exchange Rate,Rate Date,Rate Source,Images") {
		t.Errorf("receipts.csv header = %s", header)
	}

	var manifest struct {
		Claim structs.ClaimView `json:"claim"`
	}
	if err := json.Unmarshal(packet.files["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Claim.Currency != "USD" || manifest.Claim.UnconvertedReceiptCount != 1 {
		t.Errorf("manifest claim: %+v", manifest.Claim)
	}
}

func TestMulticurrency_SummaryLegacyTotalIsBaseCurrencyOnly(t *testing.T) {
	defer tearDownGenericHandlerTest()
	repositories.CreateTestGroupWithUsers()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead)
	repositories.GetDB().Create(&models.GroupReceiptSettings{GroupId: 1, ReceiptSummaryEnabled: true})

	seedForeignReceipt(t, 1, 1, "10.00", "USD", "")
	seedForeignReceipt(t, 1, 1, "90.00", "EUR", "")

	summary, err := services.NewReceiptSummaryService(nil).GetReceiptSummary(1, "1", commands.ReceiptSummaryCommand{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.BaseCurrency != "USD" || summary.Overall.Total.String() != "10" || summary.Overall.ReceiptCount != 2 {
		t.Errorf("overall: base %s total %s count %d", summary.BaseCurrency, summary.Overall.Total, summary.Overall.ReceiptCount)
	}
	if len(summary.Overall.CurrencyTotals) != 2 || summary.Overall.CurrencyTotals[1].Currency != "EUR" ||
		summary.Overall.CurrencyTotals[1].Total.String() != "90" {
		t.Errorf("currencyTotals: %+v", summary.Overall.CurrencyTotals)
	}

	wire, _ := json.Marshal(summary)
	if !strings.Contains(string(wire), `"currencyTotals":[{`) || strings.Contains(string(wire), `"currencyTotals":null`) {
		t.Errorf("wire: %s", wire)
	}
}
