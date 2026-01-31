package domain

import (
	"context"

	"backend/infra/validation"

	"github.com/guregu/null/v6"
)

type CreateOrganization struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	OwnerID string `json:"ownerId"`
}

func (c CreateOrganization) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.Slug, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.OwnerID, validation.Required, validation.IsUUID),
	)
}

type UpdateOrganization struct {
	Name    null.String `json:"name"`
	Slug    null.String `json:"slug"`
	OwnerID null.String `json:"ownerId"`
}

func (u UpdateOrganization) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.Slug, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.OwnerID, validation.NilOrNotEmpty, validation.IsUUID),
	)
}
