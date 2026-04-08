package accounttype

import "strings"

// Canonical account type values stored in the API and database (SCREAMING_SNAKE).
const (
	Checking     = "CHECKING"
	Savings      = "SAVINGS"
	Cash         = "CASH"
	CreditCard   = "CREDIT_CARD"
	MoneyMarket  = "MONEY_MARKET"
	Investment   = "INVESTMENT"
	Retirement   = "RETIREMENT"
	Loan         = "LOAN"
	Mortgage     = "MORTGAGE"
	LineOfCredit = "LINE_OF_CREDIT"
	Other        = "OTHER"
)

// All lists every allowed account type in display order.
var All = []string{
	Checking,
	Savings,
	Cash,
	CreditCard,
	MoneyMarket,
	Investment,
	Retirement,
	Loan,
	Mortgage,
	LineOfCredit,
	Other,
}

var validSet = map[string]struct{}{
	Checking:     {},
	Savings:      {},
	Cash:         {},
	CreditCard:   {},
	MoneyMarket:  {},
	Investment:   {},
	Retirement:   {},
	Loan:         {},
	Mortgage:     {},
	LineOfCredit: {},
	Other:        {},
}

// IsValid reports whether s is an allowed canonical type (after Normalize).
func IsValid(s string) bool {
	_, ok := validSet[s]
	return ok
}

// Normalize trims whitespace and uppercases s for canonical storage and comparison.
func Normalize(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// AsAnySlice returns All as []any for validation.In(...).
func AsAnySlice() []any {
	out := make([]any, len(All))
	for i, v := range All {
		out[i] = v
	}
	return out
}
