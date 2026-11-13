package structs

import "time"

// The wire shapes of the reimbursement claims feature (swagger: Claim,
// ClaimLedgerEntry, ClaimList, ClaimListSummary, ClaimReceiptsResult). The file is
// not named claims.go because that is the JWT Claims type next door.
//
// Calendar days are bare yyyy-MM-dd strings and money is a decimal string with two
// places, so neither depends on how a client parses a timestamp or a float. Every
// array is non-nil, so it serializes as [] rather than null: the generated Dart
// client has no null guard.

type ClaimView struct {
	Id                     uint                   `json:"id"`
	GroupId                uint                   `json:"groupId"`
	Name                   string                 `json:"name"`
	ReferenceNumber        string                 `json:"referenceNumber,omitempty"`
	Type                   string                 `json:"type,omitempty"`
	Status                 string                 `json:"status"`
	AwaitingSettlement     bool                   `json:"awaitingSettlement"`
	SubmittedDate          string                 `json:"submittedDate,omitempty"`
	DueDate                string                 `json:"dueDate,omitempty"`
	ClaimedAmount          string                 `json:"claimedAmount,omitempty"`
	Notes                  string                 `json:"notes,omitempty"`
	CreatedAt              time.Time              `json:"createdAt"`
	UpdatedAt              time.Time              `json:"updatedAt"`
	CreatedBy              *uint                  `json:"createdBy,omitempty"`
	ReceiptIds             []uint                 `json:"receiptIds"`
	ReceiptCount           int                    `json:"receiptCount"`
	ReceiptsTotal          string                 `json:"receiptsTotal"`
	EffectiveClaimedAmount string                 `json:"effectiveClaimedAmount"`
	PaymentsTotal          string                 `json:"paymentsTotal"`
	RecoupmentsTotal       string                 `json:"recoupmentsTotal"`
	PaidAmount             string                 `json:"paidAmount"`
	OutstandingAmount      string                 `json:"outstandingAmount"`
	DaysWaiting            *int                   `json:"daysWaiting,omitempty"`
	LedgerEntries          []ClaimLedgerEntryView `json:"ledgerEntries"`
	// Multicurrency (fork feature). Currency is the claim's group's base currency, the one every
	// figure above is reckoned in. A linked receipt counts at its amount when it is in that
	// currency, at its recorded baseAmount when it was converted, and otherwise is UNCONVERTED:
	// left out of the figures and reported here instead. UnconvertedTotals is [] when empty.
	Currency                string               `json:"currency"`
	UnconvertedReceiptCount int                  `json:"unconvertedReceiptCount"`
	UnconvertedTotals       []ClaimCurrencyTotal `json:"unconvertedTotals"`
}

// ClaimCurrencyTotal is an amount in one currency, money as a decimal string at that
// currency's own number of decimal places.
type ClaimCurrencyTotal struct {
	Currency string `json:"currency"`
	Total    string `json:"total"`
}

type ClaimLedgerEntryView struct {
	Id        uint      `json:"id"`
	ClaimId   uint      `json:"claimId"`
	Type      string    `json:"type"`
	Date      string    `json:"date"`
	Amount    string    `json:"amount,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy *uint     `json:"createdBy,omitempty"`
}

type ClaimListView struct {
	Claims  []ClaimView          `json:"claims"`
	Summary ClaimListSummaryView `json:"summary"`
}

type ClaimListSummaryView struct {
	AsOfDate         string `json:"asOfDate"`
	ClaimCount       int    `json:"claimCount"`
	AwaitingCount    int    `json:"awaitingCount"`
	OutstandingTotal string `json:"outstandingTotal"`
	// OutstandingTotals (multicurrency, fork feature) is the outstanding amount per claim
	// currency, [] when nothing is awaiting. Claims of one group share its base currency, so
	// only the All group can show several. OutstandingTotal stays the sum of the claims in one
	// currency: the only one present, else the install default.
	OutstandingTotals []ClaimCurrencyTotal `json:"outstandingTotals"`
	OldestDaysWaiting *int                 `json:"oldestDaysWaiting,omitempty"`
}

// ClaimReceiptsResult reports, per requested receipt, whether an add or remove
// took effect. Build it with NewClaimReceiptsResult so both lists are [].
type ClaimReceiptsResult struct {
	ReceiptIds []uint                 `json:"receiptIds"`
	Failures   []ClaimReceiptsFailure `json:"failures"`
}

type ClaimReceiptsFailure struct {
	ReceiptId uint   `json:"receiptId"`
	Message   string `json:"message"`
}

func NewClaimReceiptsResult() ClaimReceiptsResult {
	return ClaimReceiptsResult{
		ReceiptIds: []uint{},
		Failures:   []ClaimReceiptsFailure{},
	}
}
