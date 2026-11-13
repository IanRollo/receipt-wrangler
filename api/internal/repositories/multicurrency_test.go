package repositories

import (
	"strings"
	"testing"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func stringPtr(value string) *string { return &value }

func seedSystemSettingsWithCurrency(t *testing.T, code string) {
	t.Helper()
	settings := models.SystemSettings{DefaultCurrency: code}
	if err := GetDB().Create(&settings).Error; err != nil {
		t.Fatalf("seed settings: %v", err)
	}
}

// --- resolution ---

func TestResolveInstallCurrency(t *testing.T) {
	cases := map[string]string{"": "USD", "usd": "USD", " eur ": "EUR", "nonsense": "USD", "JPY": "JPY"}
	for stored, want := range cases {
		if got := ResolveInstallCurrency(stored); got != want {
			t.Errorf("ResolveInstallCurrency(%q) = %q, want %q", stored, got, want)
		}
	}
}

func TestGroupBaseCurrency_GroupThenInstall(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "GBP")
	CreateTestGroupWithUsers()

	if got := GroupBaseCurrency(nil, 1); got != "GBP" {
		t.Errorf("no group settings row: %q, want the install default GBP", got)
	}

	GetDB().Create(&models.GroupReceiptSettings{GroupId: 1})
	if got := GroupBaseCurrency(nil, 1); got != "GBP" {
		t.Errorf("empty group default: %q", got)
	}

	GetDB().Model(&models.GroupReceiptSettings{}).Where("group_id = ?", 1).Update("default_currency", "EUR")
	if got := GroupBaseCurrency(nil, 1); got != "EUR" {
		t.Errorf("group default: %q", got)
	}

	bases := GroupBaseCurrencies(nil, []uint{1, 2})
	if bases[1] != "EUR" || bases[2] != "GBP" {
		t.Errorf("GroupBaseCurrencies = %v", bases)
	}
}

func TestResolveReceiptCurrency_OldClientUpdateKeepsAStoredEuro(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "USD")
	CreateTestGroupWithUsers()

	existing := models.Receipt{GroupId: 1}
	existing.Currency = "EUR"
	command := commands.UpsertReceiptCommand{GroupId: 1, Amount: decimal.RequireFromString("10.00")}

	if vErr := ResolveReceiptCurrency(nil, &command, &existing); len(vErr.Errors) != 0 {
		t.Fatalf("errors: %v", vErr.Errors)
	}
	if command.Currency != "EUR" {
		t.Errorf("an update with no currency became %q", command.Currency)
	}

	// The same command on a create is the group's default.
	command = commands.UpsertReceiptCommand{GroupId: 1, Amount: decimal.RequireFromString("10.00")}
	if ResolveReceiptCurrency(nil, &command, nil); command.Currency != "USD" {
		t.Errorf("a create with no currency became %q", command.Currency)
	}
}

func TestResolveReceiptCurrency_ConversionRules(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "USD")
	CreateTestGroupWithUsers()

	dec := func(value string) *decimal.Decimal {
		parsed := decimal.RequireFromString(value)
		return &parsed
	}

	// Rate only: the base amount is amount x rate, half-even to the base exponent.
	command := commands.UpsertReceiptCommand{GroupId: 1, Currency: "EUR", Amount: *dec("33.33"), ExchangeRate: dec("1.0855")}
	ResolveReceiptCurrency(nil, &command, nil)
	if command.BaseAmount == nil || command.BaseAmount.String() != "36.18" {
		t.Errorf("derived base: %v", command.BaseAmount)
	}

	// A base amount with more places than the base currency allows is refused, not rounded.
	command = commands.UpsertReceiptCommand{GroupId: 1, Currency: "EUR", Amount: *dec("10"), BaseAmount: dec("10.123")}
	if vErr := ResolveReceiptCurrency(nil, &command, nil); vErr.Errors["baseAmount"] == "" {
		t.Errorf("a 3-place USD base amount was accepted: %v", vErr.Errors)
	}

	// Equal to the base: everything is cleared.
	command = commands.UpsertReceiptCommand{GroupId: 1, Currency: "usd", Amount: *dec("10"), BaseAmount: dec("10"), ExchangeRate: dec("1"), ExchangeRateSource: stringPtr("x")}
	ResolveReceiptCurrency(nil, &command, nil)
	if command.BaseAmount != nil || command.ExchangeRate != nil || command.ExchangeRateSource != nil {
		t.Errorf("not cleared: %+v", command)
	}

	// Sending only a source replaces the set: no base amount, no rate, so no conversion.
	existing := models.Receipt{GroupId: 1, ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "EUR", BaseAmount: dec("11")}}
	command = commands.UpsertReceiptCommand{GroupId: 1, Currency: "EUR", Amount: *dec("10"), ExchangeRateSource: stringPtr("  ")}
	ResolveReceiptCurrency(nil, &command, &existing)
	if command.BaseAmount != nil || command.ExchangeRateSource != nil {
		t.Errorf("an explicit blank did not clear: %+v", command)
	}

	// Moving to another group without conversion keys does not carry the old conversion over.
	existing = models.Receipt{GroupId: 2, ReceiptMulticurrency: models.ReceiptMulticurrency{Currency: "EUR", BaseAmount: dec("11")}}
	command = commands.UpsertReceiptCommand{GroupId: 1, Currency: "EUR", Amount: *dec("10")}
	ResolveReceiptCurrency(nil, &command, &existing)
	if command.BaseAmount != nil {
		t.Errorf("a conversion relative to another group's base survived a move")
	}
}

// --- settings ---

func TestUpdateSystemSettings_DefaultCurrencyIsPointerAndNormalized(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "EUR")
	repository := NewSystemSettingsRepository(nil)

	command := buildSettingsCommand()
	updated, err := repository.UpdateSystemSettings(command)
	if err != nil {
		t.Fatal(err)
	}
	if updated.DefaultCurrency != "EUR" {
		t.Errorf("an omitted defaultCurrency changed the stored one to %q", updated.DefaultCurrency)
	}

	command.DefaultCurrency = stringPtr("jpy")
	if vErr := command.Validate(); len(vErr.Errors) != 0 {
		t.Fatalf("validate: %v", vErr.Errors)
	}
	if _, err := repository.UpdateSystemSettings(command); err != nil {
		t.Fatal(err)
	}
	stored, _ := repository.GetSystemSettings()
	if stored.DefaultCurrency != "JPY" {
		t.Errorf("stored %q", stored.DefaultCurrency)
	}

	command.DefaultCurrency = stringPtr("ZZZ")
	if vErr := command.Validate(); vErr.Errors["defaultCurrency"] == "" {
		t.Errorf("an unknown code was accepted")
	}
	command.DefaultCurrency = stringPtr("")
	if vErr := command.Validate(); vErr.Errors["defaultCurrency"] == "" {
		t.Errorf("an empty default was accepted")
	}
}

func TestUpdateGroupReceiptSettings_DefaultCurrencyIsPointerAndExposedResolved(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "GBP")
	CreateTestGroupWithUsers()
	repository := NewGroupReceiptSettingsRepository(nil)
	if _, err := repository.CreateGroupReceiptSettings(1); err != nil {
		t.Fatal(err)
	}

	read := func() models.GroupReceiptSettings {
		settings, err := repository.GetGroupReceiptSettingsByGroupId(1)
		if err != nil {
			t.Fatal(err)
		}
		return settings
	}

	if settings := read(); settings.DefaultCurrency != "" || settings.BaseCurrency != "GBP" {
		t.Errorf("fresh: default %q base %q", settings.DefaultCurrency, settings.BaseCurrency)
	}

	command := commands.UpdateGroupReceiptSettingsCommand{
		QuickScanPaidByEnabled: true, QuickScanPaidByRequired: true,
		QuickScanStatusEnabled: true, QuickScanStatusRequired: true,
		DefaultCurrency: stringPtr("eur"),
	}
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		t.Fatal(err)
	}
	if settings := read(); settings.DefaultCurrency != "EUR" || settings.BaseCurrency != "EUR" {
		t.Errorf("set: default %q base %q", settings.DefaultCurrency, settings.BaseCurrency)
	}

	// An omitted key leaves it alone.
	command.DefaultCurrency = nil
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		t.Fatal(err)
	}
	if settings := read(); settings.DefaultCurrency != "EUR" {
		t.Errorf("omitted key changed it to %q", settings.DefaultCurrency)
	}

	// An explicit empty clears it back to the install default.
	command.DefaultCurrency = stringPtr("")
	if _, err := repository.UpdateGroupReceiptSettings("1", command); err != nil {
		t.Fatal(err)
	}
	if settings := read(); settings.DefaultCurrency != "" || settings.BaseCurrency != "GBP" {
		t.Errorf("cleared: default %q base %q", settings.DefaultCurrency, settings.BaseCurrency)
	}

	command.DefaultCurrency = stringPtr("ZZZ")
	if vErr := command.Validate(); vErr.Errors["defaultCurrency"] == "" {
		t.Errorf("an unknown code was accepted")
	}
}

// --- migrations ---

func TestBackfillReceiptCurrency_LabelsOldReceiptsWithTheInstallDefaultOnly(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "EUR")
	CreateTestGroupWithUsers()

	old := models.Receipt{Name: "old", GroupId: 1, PaidByUserID: 1, Status: models.OPEN, Amount: decimal.NewFromInt(1)}
	foreign := models.Receipt{Name: "jpy", GroupId: 1, PaidByUserID: 1, Status: models.OPEN, Amount: decimal.NewFromInt(1)}
	foreign.Currency = "JPY"
	GetDB().Select("*").Create(&old)
	GetDB().Create(&foreign)
	GetDB().Model(&models.Receipt{}).Where("id = ?", old.ID).UpdateColumn("currency", nil)

	if err := backfillReceiptCurrency(GetDB()); err != nil {
		t.Fatal(err)
	}

	var reloaded models.Receipt
	GetDB().First(&reloaded, old.ID)
	if reloaded.Currency != "EUR" {
		t.Errorf("old receipt = %q, want the install default", reloaded.Currency)
	}
	var labelled models.Receipt
	GetDB().First(&labelled, foreign.ID)
	if labelled.Currency != "JPY" {
		t.Errorf("an already-labelled receipt became %q", labelled.Currency)
	}

	// Changing the default later never re-labels what the backfill stored.
	GetDB().Model(&models.SystemSettings{}).Where("1 = 1").Update("default_currency", "GBP")
	backfillReceiptCurrency(GetDB())
	var again models.Receipt
	GetDB().First(&again, old.ID)
	reloaded = again
	if reloaded.Currency != "EUR" {
		t.Errorf("a later default re-labelled the receipt: %q", reloaded.Currency)
	}
}

func TestBackfillReceiptCurrency_FreshSettingsRowMeansUSD(t *testing.T) {
	defer TruncateTestDb()
	GetDB().Create(&models.SystemSettings{})
	CreateTestGroupWithUsers()

	receipt := models.Receipt{Name: "old", GroupId: 1, PaidByUserID: 1, Status: models.OPEN, Amount: decimal.NewFromInt(1)}
	GetDB().Select("*").Create(&receipt)
	GetDB().Model(&models.Receipt{}).Where("id = ?", receipt.ID).UpdateColumn("currency", "")

	if err := backfillReceiptCurrency(GetDB()); err != nil {
		t.Fatal(err)
	}
	var reloaded models.Receipt
	GetDB().First(&reloaded, receipt.ID)
	if reloaded.Currency != "USD" {
		t.Errorf("currency = %q", reloaded.Currency)
	}
}

func TestUpgradeDefaultPromptCurrency_UpdatesAnUntouchedPromptOnly(t *testing.T) {
	defer TruncateTestDb()

	// Stored with the line endings a Windows-edited database could carry.
	stored := models.Prompt{
		Name:   constants.DefaultPromptName,
		Prompt: strings.ReplaceAll(constants.PreviousDefaultPromptText, "\n", "\r\n"),
	}
	GetDB().Create(&stored)

	if err := upgradeDefaultPromptCurrency(GetDB()); err != nil {
		t.Fatal(err)
	}
	var reloaded models.Prompt
	GetDB().First(&reloaded, stored.ID)
	if reloaded.Prompt != constants.DefaultPromptText || !strings.Contains(reloaded.Prompt, `"currency"`) {
		t.Errorf("an unedited default prompt was not upgraded")
	}
}

func TestUpgradeDefaultPromptCurrency_LeavesAnEditedPromptAlone(t *testing.T) {
	defer TruncateTestDb()

	edited := constants.PreviousDefaultPromptText + "\nAlways answer in pirate.\n"
	stored := models.Prompt{Name: constants.DefaultPromptName, Prompt: edited}
	GetDB().Create(&stored)

	if err := upgradeDefaultPromptCurrency(GetDB()); err != nil {
		t.Fatal(err)
	}
	var reloaded models.Prompt
	GetDB().First(&reloaded, stored.ID)
	if reloaded.Prompt != edited {
		t.Errorf("an edited prompt was overwritten")
	}

	// Nothing stored at all is fine too.
	TruncateTestDb()
	if err := upgradeDefaultPromptCurrency(GetDB()); err != nil {
		t.Errorf("no prompt: %v", err)
	}
}

func TestDefaultPromptTexts(t *testing.T) {
	if strings.Contains(constants.PreviousDefaultPromptText, `"currency"`) {
		t.Errorf("the previous prompt must be the text before multicurrency")
	}
	if !strings.Contains(constants.DefaultPromptText, `"currency"`) {
		t.Errorf("the new prompt does not ask for a currency")
	}
	if !strings.Contains(constants.DefaultPromptText, "@ocrText") || !strings.Contains(constants.DefaultPromptText, "@categories") {
		t.Errorf("the new prompt lost a placeholder")
	}
}

func TestMigrationsAreRegistered(t *testing.T) {
	names := map[string]bool{}
	for _, migration := range dataMigrations {
		names[migration.name] = true
	}
	if !names[backfillReceiptCurrencyMigration] || !names[upgradeDefaultPromptCurrencyMigration] {
		t.Errorf("multicurrency migrations missing from the registry: %v", names)
	}
}

// --- column widening ---

// An existing install has receipts.amount as decimal(10,2); AutoMigrate must move it to the new
// type without losing rows. SQLite keeps no precision, so this proves the migration runs and the
// data survives there; the Postgres/MySQL ALTER is documented, not exercised (api/CLAUDE.md).
func TestAutoMigrateWidensReceiptAmountWithoutLosingRows(t *testing.T) {
	legacy, err := gorm.Open(sqlite.Open("file:widen_receipts?mode=memory&cache=shared&_pragma=foreign_keys(0)"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Current schema first (it creates the tables receipts points at), then the receipts table
	// is put back to what a pre-multicurrency install has: no currency columns, decimal(10,2).
	if err := legacy.AutoMigrate(&models.Receipt{}, &models.Item{}); err != nil {
		t.Fatalf("initial AutoMigrate: %v", err)
	}
	statements := []string{
		`DROP TABLE receipts`,
		`CREATE TABLE receipts (id integer primary key autoincrement, created_at datetime, updated_at datetime, created_by integer, created_by_string text, name text not null, amount decimal(10,2) not null, date datetime not null, resolved_date datetime, paid_by_user_id integer, status text not null default 'OPEN', group_id integer not null, email_inbox text, email_sender text)`,
		`INSERT INTO receipts (name, amount, date, group_id) VALUES ('old', 12.34, '2026-01-01', 1)`,
	}
	for _, statement := range statements {
		if err := legacy.Exec(statement).Error; err != nil {
			t.Fatalf("legacy schema: %v", err)
		}
	}

	if err := legacy.AutoMigrate(&models.Receipt{}); err != nil {
		t.Fatalf("AutoMigrate over the legacy table: %v", err)
	}

	var receipt models.Receipt
	if err := legacy.First(&receipt).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if receipt.Name != "old" || receipt.Amount.String() != "12.34" {
		t.Errorf("row after migrate: %q %s", receipt.Name, receipt.Amount)
	}
	for _, column := range []string{"currency", "base_amount", "exchange_rate", "exchange_rate_date", "exchange_rate_source"} {
		if !legacy.Migrator().HasColumn(&models.Receipt{}, column) {
			t.Errorf("column %s was not added", column)
		}
	}

	// A three-place amount round-trips (SQLite keeps the declared decimal(10,2) text but
	// stores any number as given).
	if err := legacy.Model(&receipt).Update("amount", decimal.RequireFromString("5.125")).Error; err != nil {
		t.Fatal(err)
	}
	var reloaded models.Receipt
	legacy.First(&reloaded, receipt.ID)
	if reloaded.Amount.String() != "5.125" {
		t.Errorf("three-place amount stored as %s", reloaded.Amount)
	}
}

func TestStampConversionCurrency_StampsOnlyRowsWithAnUnstampedConversion(t *testing.T) {
	defer TruncateTestDb()
	seedSystemSettingsWithCurrency(t, "EUR")
	CreateTestGroupWithUsers()

	base := decimal.NewFromInt(5)
	converted := models.Receipt{Name: "c", GroupId: 1, PaidByUserID: 1, Status: models.OPEN, Amount: decimal.NewFromInt(1)}
	converted.Currency = "GBP"
	converted.BaseAmount = &base
	plain := models.Receipt{Name: "p", GroupId: 1, PaidByUserID: 1, Status: models.OPEN, Amount: decimal.NewFromInt(1)}
	plain.Currency = "EUR"
	GetDB().Create(&converted)
	GetDB().Create(&plain)

	if err := stampConversionCurrency(GetDB()); err != nil {
		t.Fatal(err)
	}
	var got, other models.Receipt
	GetDB().First(&got, converted.ID)
	GetDB().First(&other, plain.ID)
	if got.BaseAmountCurrency != "EUR" || other.BaseAmountCurrency != "" {
		t.Errorf("stamps: converted %q plain %q", got.BaseAmountCurrency, other.BaseAmountCurrency)
	}
}
