package port

import (
	"context"

	"backend/adapter/validation"
	"backend/infra/money"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateOrganizationCurrency struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	CurrencyCode   string    `json:"currencyCode"`
	IsBase       bool               `json:"isBase"`
	Rate         money.ExchangeRate `json:"rate"`
}

func (c CreateOrganizationCurrency) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.CurrencyCode, validation.Required, validation.Length(3, 3)),
	)
}

type UpdateOrganizationCurrency struct {
	IsBase null.Bool               `json:"isBase"`
	Rate   money.NullExchangeRate `json:"rate"`
}

func (u UpdateOrganizationCurrency) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u)
}
