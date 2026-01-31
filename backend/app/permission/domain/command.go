package domain

import (
	"context"

	"backend/infra/validation"

	"github.com/guregu/null/v6"
)

type CreatePermission struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func (c CreatePermission) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.Slug, validation.Required, validation.Length(2, 100)),
	)
}

type UpdatePermission struct {
	Description null.String `json:"description"`
}

func (u UpdatePermission) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Description, validation.NilOrNotEmpty),
	)
}
