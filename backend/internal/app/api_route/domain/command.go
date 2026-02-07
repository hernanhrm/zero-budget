package domain

import (
	"backend/infra/validation"
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateApiRoute struct {
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	PermissionID uuid.UUID `json:"permissionId"`
}

func (c CreateApiRoute) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.Method, validation.Required),
		validation.Field(&c.Path, validation.Required),
		validation.Field(&c.PermissionID, validation.Required, validation.IsUUID),
	)
}

type UpdateApiRoute struct {
	PermissionID null.String `json:"permissionId"`
}

func (u UpdateApiRoute) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.PermissionID, validation.NilOrNotEmpty, validation.IsUUID),
	)
}
