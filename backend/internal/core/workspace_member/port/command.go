package port

import (
	"context"

	"backend/adapter/validation"

	"github.com/google/uuid"
)

type CreateWorkspaceMember struct {
	WorkspaceID uuid.UUID `json:"workspaceId"`
	UserID      uuid.UUID `json:"userId"`
	RoleID      uuid.UUID `json:"roleId"`
}

func (c CreateWorkspaceMember) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.WorkspaceID, validation.Required),
		validation.Field(&c.UserID, validation.Required),
		validation.Field(&c.RoleID, validation.Required),
	)
}

type UpdateWorkspaceMember struct {
	RoleID uuid.UUID `json:"roleId"`
}

func (u UpdateWorkspaceMember) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.RoleID, validation.Required),
	)
}
