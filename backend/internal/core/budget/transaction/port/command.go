package port

import (
	"context"
	"time"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateTransaction struct {
	ID                      uuid.UUID   `json:"id"`
	OrganizationID          string      `json:"organizationId"`
	AccountID               uuid.UUID   `json:"accountId"`
	CategoryID              *uuid.UUID  `json:"categoryId"`
	SubcategoryID           *uuid.UUID  `json:"subcategoryId"`
	BudgetID                *uuid.UUID  `json:"budgetId"`
	Type                    string      `json:"type"`
	Amount                  int64       `json:"amount"`
	Description             null.String `json:"description"`
	ExternalReferenceNumber null.String `json:"externalReferenceNumber"`
	Date                    time.Time   `json:"date"`
}

func (c CreateTransaction) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.AccountID, validation.Required, validation.IsUUID),
		validation.Field(&c.Type, validation.Required, validation.Length(1, 20)),
		validation.Field(&c.Amount, validation.Required),
		validation.Field(&c.Date, validation.Required),
	)
}

type UpdateTransaction struct {
	CategoryID              *uuid.UUID  `json:"categoryId"`
	SubcategoryID           *uuid.UUID  `json:"subcategoryId"`
	BudgetID                *uuid.UUID  `json:"budgetId"`
	Type                    null.String `json:"type"`
	Amount                  null.Int    `json:"amount"`
	Description             null.String `json:"description"`
	ExternalReferenceNumber null.String `json:"externalReferenceNumber"`
	Date                    *time.Time  `json:"date"`
}

func (u UpdateTransaction) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Type, validation.NilOrNotEmpty, validation.Length(1, 20)),
	)
}
