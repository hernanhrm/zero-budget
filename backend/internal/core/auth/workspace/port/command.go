package port

import (
	"context"

	"backend/adapter/validation"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateWorkspace struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
}

func (c CreateWorkspace) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.OrganizationID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Slug, validation.Required, validation.Length(2, 255)),
	)
}

type UpdateWorkspace struct {
	Name null.String `json:"name"`
	Slug null.String `json:"slug"`
}

func (u UpdateWorkspace) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 255)),
		validation.Field(&u.Slug, validation.NilOrNotEmpty, validation.Length(2, 255)),
	)
}
