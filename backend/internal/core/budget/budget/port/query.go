package port

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Month          int16     `json:"month"`
	Year           int16     `json:"year"`
	CurrencyCode   string    `json:"currencyCode"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
