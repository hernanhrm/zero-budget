package port

import (
	"context"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type CreateEmailTemplate struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Event       string    `json:"event"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Content     string    `json:"content"`
	IsActive    bool      `json:"isActive"`
	Locale      string    `json:"locale"`
}

func (c CreateEmailTemplate) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Event, validation.Required, validation.Length(1, 100)),
		validation.Field(&c.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&c.Description, validation.Length(0, 1000)),
		validation.Field(&c.Subject, validation.Required, validation.Length(1, 500)),
		validation.Field(&c.Content, validation.Required),
		validation.Field(&c.Locale, validation.Length(2, 10)),
	)
}

type UpdateEmailTemplate struct {
	Event       null.String `json:"event"`
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	Subject     null.String `json:"subject"`
	Content     null.String `json:"content"`
	IsActive    null.Bool   `json:"isActive"`
	Locale      null.String `json:"locale"`
}

func (u UpdateEmailTemplate) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Event, validation.NilOrNotEmpty, validation.Length(1, 100)),
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&u.Description, validation.Length(0, 1000)),
		validation.Field(&u.Subject, validation.NilOrNotEmpty, validation.Length(1, 500)),
		validation.Field(&u.Content, validation.NilOrNotEmpty),
		validation.Field(&u.Locale, validation.NilOrNotEmpty, validation.Length(2, 10)),
	)
}
