package repositories

import (
	"fmt"
	"strconv"
	"strings"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
)

// ResolveInstallCurrency turns the stored install default into the code in force: a valid
// code (normalized), else USD. It is what AppData.defaultCurrency carries, so it is never "".
func ResolveInstallCurrency(stored string) string {
	normalized := currency.Normalize(stored)
	if currency.IsValid(normalized) {
		return normalized
	}
	return currency.DefaultInstallCurrency
}

// InstallDefaultCurrency reads the install default currency live (no cache), like
// GetAppLocation. An unreadable settings row resolves to USD.
func InstallDefaultCurrency(tx *gorm.DB) string {
	settings, err := NewSystemSettingsRepository(tx).GetSystemSettings()
	if err != nil {
		return currency.DefaultInstallCurrency
	}
	return ResolveInstallCurrency(settings.DefaultCurrency)
}

// ResolveGroupCurrency turns a group's stored default into the code in force, which is also
// the group's base currency: its own valid default, else the install default.
func ResolveGroupCurrency(stored string, installDefault string) string {
	normalized := currency.Normalize(stored)
	if currency.IsValid(normalized) {
		return normalized
	}
	return installDefault
}

// GroupBaseCurrency is a group's resolved default currency, the currency its claims are
// reckoned in. A group with no settings row yet, or none set, uses the install default.
func GroupBaseCurrency(tx *gorm.DB, groupId uint) string {
	db := tx
	if db == nil {
		db = GetDB()
	}

	var stored string
	db.Model(&models.GroupReceiptSettings{}).Where("group_id = ?", groupId).Limit(1).Pluck("default_currency", &stored)

	return ResolveGroupCurrency(stored, InstallDefaultCurrency(tx))
}

// GroupBaseCurrencies resolves several groups at once (one query), for callers that fold
// receipts of many groups.
func GroupBaseCurrencies(tx *gorm.DB, groupIds []uint) map[uint]string {
	db := tx
	if db == nil {
		db = GetDB()
	}

	installDefault := InstallDefaultCurrency(tx)
	result := make(map[uint]string, len(groupIds))
	for _, groupId := range groupIds {
		result[groupId] = installDefault
	}
	if len(groupIds) == 0 {
		return result
	}

	var rows []struct {
		GroupId         uint
		DefaultCurrency string
	}
	db.Model(&models.GroupReceiptSettings{}).Select("group_id, default_currency").Where("group_id IN ?", groupIds).Scan(&rows)
	for _, row := range rows {
		result[row.GroupId] = ResolveGroupCurrency(row.DefaultCurrency, installDefault)
	}

	return result
}

// ResolveReceiptCurrency settles a receipt command's currency block, in place, before it is
// converted to a model. It is idempotent, so the handlers call it to turn problems into
// 400s and CreateReceipt/UpdateReceipt call it again as the chokepoint every writer
// (form, quick scan, email, AI) reaches.
//
// Currency: the request's, else (update) the stored one, else the group's default, else the
// install default; never empty afterwards. Every amount is checked against the currency's
// exponent and refused, not rounded, when it has too many places (errors keyed `amount` and
// `receiptItems.N.amount`).
//
// Conversion (meaningful only when the currency differs from the group's base currency):
//   - equal to the base: the conversion fields are cleared;
//   - none of the four keys sent: an update keeps the stored values while the currency and
//     group are unchanged, otherwise there is no conversion;
//   - any key sent: the request is authoritative for the whole set. baseAmount without a rate
//     derives the rate (base / amount, 6 dp), a rate without baseAmount derives the base amount
//     (amount x rate, half-even to the base currency's exponent), both are kept as given.
//
// existing is the stored receipt on an update, nil on a create.
func ResolveReceiptCurrency(tx *gorm.DB, command *commands.UpsertReceiptCommand, existing *models.Receipt) structs.ValidatorError {
	errors := command.ValidateMulticurrency()
	if len(errors) > 0 {
		return structs.ValidatorError{Errors: errors}
	}

	code := currency.Normalize(command.Currency)

	baseCurrency := GroupBaseCurrency(tx, command.GroupId)
	if code == "" && existing != nil && currency.IsValid(existing.Currency) {
		code = currency.Normalize(existing.Currency)
	}
	if code == "" {
		code = baseCurrency
	}
	command.Currency = code

	if msg := currency.ValidateAmount(command.Amount, code); msg != "" {
		errors["amount"] = msg
	}
	for index, item := range command.Items {
		if msg := currency.ValidateAmount(item.Amount, code); msg != "" {
			errors["receiptItems."+itoa(index)+".amount"] = msg
		}
	}

	resolveReceiptConversion(command, existing, code, baseCurrency, errors)

	return structs.ValidatorError{Errors: errors}
}

func resolveReceiptConversion(
	command *commands.UpsertReceiptCommand,
	existing *models.Receipt,
	code string,
	baseCurrency string,
	errors map[string]string,
) {
	supplied := command.BaseAmount != nil || command.ExchangeRate != nil ||
		command.ExchangeRateDate != nil || command.ExchangeRateSource != nil

	if code == baseCurrency {
		clearConversion(command)
		return
	}

	if !supplied {
		keepStored := existing != nil && existing.GroupId == command.GroupId && currency.Normalize(existing.Currency) == code
		if keepStored {
			command.BaseAmount = existing.BaseAmount
			command.ExchangeRate = existing.ExchangeRate
			if existing.ExchangeRateDate != nil {
				day := existing.ExchangeRateDate.UTC().Format("2006-01-02")
				command.ExchangeRateDate = &day
			}
			if existing.ExchangeRateSource != "" {
				source := existing.ExchangeRateSource
				command.ExchangeRateSource = &source
			}
		}
		return
	}

	if command.ExchangeRate != nil {
		rate := command.ExchangeRate.Round(6)
		if !rate.IsPositive() {
			errors["exchangeRate"] = "Exchange rate must be greater than zero"
			return
		}
		command.ExchangeRate = &rate
	}

	if command.BaseAmount != nil {
		if msg := currency.ValidateAmount(*command.BaseAmount, baseCurrency); msg != "" {
			errors["baseAmount"] = msg
			return
		}
	}

	command.BaseAmount, command.ExchangeRate = currency.DeriveConversion(
		command.Amount, baseCurrency, command.BaseAmount, command.ExchangeRate,
	)

	if command.ExchangeRateDate != nil && strings.TrimSpace(*command.ExchangeRateDate) == "" {
		command.ExchangeRateDate = nil
	}
	if command.ExchangeRateSource != nil {
		source := strings.TrimSpace(*command.ExchangeRateSource)
		if source == "" {
			command.ExchangeRateSource = nil
		} else {
			command.ExchangeRateSource = &source
		}
	}
}

func clearConversion(command *commands.UpsertReceiptCommand) {
	command.BaseAmount = nil
	command.ExchangeRate = nil
	command.ExchangeRateDate = nil
	command.ExchangeRateSource = nil
}

func itoa(value int) string {
	return strconv.Itoa(value)
}

// ResolveReceiptCurrencyForWrite is the repository-side entry point: it loads the stored
// receipt for an update (receiptId > 0), runs ResolveReceiptCurrency and turns any
// validation problem into an error, since by this point the handlers have already answered
// the caller's mistakes with a 400.
func ResolveReceiptCurrencyForWrite(tx *gorm.DB, receiptId uint, command *commands.UpsertReceiptCommand) error {
	var existing *models.Receipt
	if receiptId > 0 {
		var stored models.Receipt
		db := tx
		if db == nil {
			db = GetDB()
		}
		err := db.Select("id", "group_id", "currency", "base_amount", "exchange_rate", "exchange_rate_date", "exchange_rate_source").
			Limit(1).Find(&stored, receiptId).Error
		if err != nil {
			return err
		}
		existing = &stored
	}

	vErr := ResolveReceiptCurrency(tx, command, existing)
	for key, message := range vErr.Errors {
		return fmt.Errorf("invalid receipt %s: %s", key, message)
	}
	return nil
}
