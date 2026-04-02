package port

import (
	"context"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateAccount struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	CurrencyCode   string    `json:"currencyCode"`
	CurrentBalance int64     `json:"currentBalance"`
	IsActive       bool      `json:"isActive"`
}

func (c CreateAccount) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Type, validation.Required, validation.Length(1, 50)),
		validation.Field(&c.CurrencyCode, validation.Required, validation.Length(3, 3)),
	)
}

type UpdateAccount struct {
	Name           null.String `json:"name"`
	Type           null.String `json:"type"`
	CurrencyCode   null.String `json:"currencyCode"`
	CurrentBalance null.Int    `json:"currentBalance"`
	IsActive       null.Bool   `json:"isActive"`
}

func (u UpdateAccount) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 255)),
		validation.Field(&u.Type, validation.NilOrNotEmpty, validation.Length(1, 50)),
		validation.Field(&u.CurrencyCode, validation.NilOrNotEmpty, validation.Length(3, 3)),
	)
}
