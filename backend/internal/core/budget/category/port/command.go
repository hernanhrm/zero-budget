package port

import (
	"context"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateCategory struct {
	ID             uuid.UUID   `json:"id"`
	OrganizationID string      `json:"organizationId"`
	ParentID       *uuid.UUID  `json:"parentId"`
	Name           string      `json:"name"`
	Icon           null.String `json:"icon"`
	Color          null.String `json:"color"`
	IsActive       bool        `json:"isActive"`
}

func (c CreateCategory) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 255)),
		validation.Field(&c.Icon, validation.NilOrNotEmpty, validation.Length(1, 50)),
		validation.Field(&c.Color, validation.NilOrNotEmpty, validation.Length(4, 7)),
	)
}

type UpdateCategory struct {
	ParentID *uuid.UUID  `json:"parentId"`
	Name     null.String `json:"name"`
	Icon     null.String `json:"icon"`
	Color    null.String `json:"color"`
	IsActive null.Bool   `json:"isActive"`
}

func (u UpdateCategory) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 255)),
		validation.Field(&u.Icon, validation.NilOrNotEmpty, validation.Length(1, 50)),
		validation.Field(&u.Color, validation.NilOrNotEmpty, validation.Length(4, 7)),
	)
}
