package domain

import (
	"context"

	"backend/infra/validation"

	"github.com/guregu/null/v6"
)

type CreateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (c CreateUser) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.FirstName, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.LastName, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.Email, validation.Required, validation.IsEmail),
	)
}

type UpdateUser struct {
	FirstName null.String `json:"firstName"`
	LastName  null.String `json:"lastName"`
	Email     null.String `json:"email"`
}

func (u UpdateUser) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.FirstName, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.LastName, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.Email, validation.NilOrNotEmpty, validation.IsEmail),
	)
}
