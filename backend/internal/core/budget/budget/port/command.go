package port

import (
	"context"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateBudget struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Month          int16     `json:"month"`
	Year           int16     `json:"year"`
	CurrencyCode   string    `json:"currencyCode"`
	IsActive       bool      `json:"isActive"`
}

func (c CreateBudget) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Month, validation.Required),
		validation.Field(&c.Year, validation.Required),
		validation.Field(&c.CurrencyCode, validation.Required, validation.Length(3, 3)),
	)
}

type UpdateBudget struct {
	Name     null.String `json:"name"`
	IsActive null.Bool   `json:"isActive"`
}

func (u UpdateBudget) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 255)),
	)
}
