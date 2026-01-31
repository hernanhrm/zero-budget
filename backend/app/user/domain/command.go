package domain

import (
	"context"

	"backend/infra/validation"

	"github.com/guregu/null/v6"
)

type CreateUser struct {
	Name  string
	Email string
}

func (c CreateUser) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.Email, validation.Required, validation.IsEmail),
	)
}

type UpdateUser struct {
	Name  null.String
	Email null.String
}

func (u UpdateUser) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.Email, validation.NilOrNotEmpty, validation.IsEmail),
	)
}
