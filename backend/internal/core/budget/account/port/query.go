package port

import (
	"time"

	"backend/infra/money"
	"github.com/google/uuid"
)

type Account struct {
	ID             uuid.UUID   `json:"id"`
	OrganizationID string      `json:"organizationId"`
	Name           string      `json:"name"`
	Type           string      `json:"type"`
	Institution    string      `json:"institution"`
	AccountNumber  string      `json:"accountNumber"`
	CurrencyCode   string      `json:"currencyCode"`
	CurrentBalance money.Minor `json:"currentBalance"`
	IsActive       bool        `json:"isActive"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}
