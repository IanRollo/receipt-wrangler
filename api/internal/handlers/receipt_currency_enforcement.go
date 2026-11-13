package handlers

import (
	"net/http"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
)

// enforceReceiptCurrency answers a receipt write whose currency block is wrong with a 400
// validator error (keyed currency / amount / receiptItems.N.amount / baseAmount /
// exchangeRate), writing it itself, and returns false when it did. It resolves a copy of the
// command, so nothing is changed here: CreateReceipt/UpdateReceipt resolve the real one again.
// existing is the stored receipt on an update, nil on a create.
func enforceReceiptCurrency(w http.ResponseWriter, command commands.UpsertReceiptCommand, existing *models.Receipt) bool {
	vErr := repositories.ResolveReceiptCurrency(nil, &command, existing)
	if len(vErr.Errors) == 0 {
		return true
	}

	structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
	return false
}
