package handlers

import (
	"net/http"

	"receipt-wrangler/api/internal/currency"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"github.com/shopspring/decimal"
)

// enforceClaimAmountExponent refuses (400, keyed `key`) a claim or ledger amount with more
// decimal places than its claim's base currency allows: claim figures are reckoned in the
// group's base currency, so a JPY claim takes whole yen and a JOD claim up to three places.
// A nil amount is nothing to check. It writes the rejection itself and returns false when
// it did.
func enforceClaimAmountExponent(w http.ResponseWriter, groupId uint, key string, amount *decimal.Decimal) bool {
	if amount == nil {
		return true
	}

	baseCurrency := repositories.GroupBaseCurrency(nil, groupId)
	message := currency.ValidateAmount(*amount, baseCurrency)
	if message == "" {
		return true
	}

	structs.WriteValidatorErrorResponse(w, structs.ValidatorError{Errors: map[string]string{key: message}}, http.StatusBadRequest)
	return false
}

// claimGroupId is the group a claim lives in, for the checks that need its base currency.
func claimGroupId(claimId uint) (uint, error) {
	claim, err := repositories.NewClaimRepository(nil).GetClaimById(claimId)
	return claim.GroupId, err
}

// enforceClaimAmountExponentForClaim is enforceClaimAmountExponent for a handler that has the
// claim's id rather than its group. A claim that cannot be read is left to the service, which
// answers it the standard way.
func enforceClaimAmountExponentForClaim(w http.ResponseWriter, claimId uint, key string, amount *decimal.Decimal) bool {
	groupId, err := claimGroupId(claimId)
	if err != nil {
		return true
	}
	return enforceClaimAmountExponent(w, groupId, key, amount)
}
