package repositories

import (
	"fmt"
	"strings"

	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"

	"gorm.io/gorm"
)

// Names of the multicurrency data migrations (fork feature).
const (
	backfillReceiptCurrencyMigration      = "backfill-receipt-currency"
	upgradeDefaultPromptCurrencyMigration = "upgrade-default-prompt-currency"
)

// backfillReceiptCurrency labels every receipt that has no currency with the install default.
// Existing amounts were always in the install's single currency; storing it explicitly means
// a later change of the default never re-labels old receipts.
func backfillReceiptCurrency(tx *gorm.DB) error {
	installDefault := InstallDefaultCurrency(tx)

	return tx.Model(&models.Receipt{}).
		Where("currency IS NULL OR currency = ''").
		UpdateColumn("currency", installDefault).Error
}

// upgradeDefaultPromptCurrency moves the stored default prompt to the version that asks for a
// currency, but only when it is byte-for-byte the previous default (ignoring line endings and
// surrounding whitespace). A prompt the user edited is never touched; it is logged once so the
// operator knows scans will rely on the group default for currency.
func upgradeDefaultPromptCurrency(tx *gorm.DB) error {
	var prompt models.Prompt
	err := tx.Where("name = ?", constants.DefaultPromptName).Limit(1).Find(&prompt).Error
	if err != nil {
		return err
	}
	if prompt.ID == 0 {
		return nil
	}

	if normalizePromptText(prompt.Prompt) == normalizePromptText(constants.PreviousDefaultPromptText) {
		return tx.Model(&prompt).UpdateColumn("prompt", constants.DefaultPromptText).Error
	}

	if !strings.Contains(prompt.Prompt, `"currency"`) {
		logging.LogStd(
			logging.LOG_LEVEL_INFO,
			fmt.Sprintf("The %q prompt was edited and does not ask for a currency; scanned receipts will use their group's default currency. Add a \"currency\" key to it to read the currency from the receipt.", constants.DefaultPromptName),
		)
	}

	return nil
}

func normalizePromptText(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
}
