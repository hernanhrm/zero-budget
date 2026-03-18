package port

import (
	"context"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateRole struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (c CreateRole) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.WorkspaceID, validation.Required),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&c.Description, validation.Length(0, 500)),
	)
}

type UpdateRole struct {
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
}

func (u UpdateRole) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&u.Description, validation.Length(0, 500)),
	)
}
