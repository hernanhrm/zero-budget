package port

import (
	"context"
	"errors"

	"backend/adapter/validation"
	"backend/core/budget/account/accounttype"
	"backend/infra/money"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateAccount struct {
	ID             uuid.UUID         `json:"id"`
	OrganizationID string            `json:"organizationId"`
	Name           string            `json:"name"`
	Type           string            `json:"type"`
	Institution    string            `json:"institution"`
	AccountNumber  string            `json:"accountNumber"`
	CurrencyCode   string            `json:"currencyCode"`
	CurrentBalance money.MajorAmount `json:"currentBalance"`
	IsActive       bool              `json:"isActive"`
}

func (c CreateAccount) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Type, validation.Required, validation.Length(1, 50), validation.In(accounttype.AsAnySlice()...)),
		validation.Field(&c.Institution, validation.Length(0, 255)),
		validation.Field(&c.AccountNumber, validation.Length(0, 64)),
		validation.Field(&c.CurrencyCode, validation.Required, validation.Length(3, 3)),
	)
}

type UpdateAccount struct {
	Name           null.String `json:"name"`
	Type           null.String `json:"type"`
	Institution    null.String `json:"institution"`
	AccountNumber  null.String `json:"accountNumber"`
	CurrencyCode   null.String `json:"currencyCode"`
	CurrentBalance null.Int    `json:"currentBalance"`
	IsActive       null.Bool   `json:"isActive"`
}

func (u UpdateAccount) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 255)),
		validation.Field(&u.Type, validation.By(validateUpdateAccountType)),
		validation.Field(&u.Institution, validation.When(u.Institution.Valid, validation.Length(0, 255))),
		validation.Field(&u.AccountNumber, validation.When(u.AccountNumber.Valid, validation.Length(0, 64))),
		validation.Field(&u.CurrencyCode, validation.NilOrNotEmpty, validation.Length(3, 3)),
	)
}

func validateUpdateAccountType(v any) error {
	ns, ok := v.(null.String)
	if !ok || !ns.Valid {
		return nil
	}
	t := accounttype.Normalize(ns.String)
	if t == "" {
		return errors.New("cannot be blank")
	}
	if !accounttype.IsValid(t) {
		return errors.New("must be a valid account type")
	}
	return nil
}
