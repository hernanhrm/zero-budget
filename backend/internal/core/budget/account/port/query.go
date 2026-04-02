package port

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	CurrencyCode   string    `json:"currencyCode"`
	CurrentBalance int64     `json:"currentBalance"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
