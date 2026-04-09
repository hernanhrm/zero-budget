package port

import (
	"time"

	"github.com/google/uuid"
)

// RelationCurrencies loads nested currency details on organization currency responses.
const RelationCurrencies = "currencies"

// OrganizationCurrencyCurrency is the nested currency payload for organization currency APIs.
// It is separate from currency/port.Currency so this module owns its contract shape.
type OrganizationCurrencyCurrency struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int16  `json:"decimalPlaces"`
}

type OrganizationCurrency struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	CurrencyCode   string    `json:"currencyCode"`
	IsBase         bool      `json:"isBase"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Currency       *OrganizationCurrencyCurrency `json:"currency,omitempty"`
}
