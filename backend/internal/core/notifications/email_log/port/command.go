package port

import (
	"context"
	"errors"

	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type EmailLogStatus string

const (
	EmailLogStatusPending   EmailLogStatus = "pending"
	EmailLogStatusSent      EmailLogStatus = "sent"
	EmailLogStatusFailed    EmailLogStatus = "failed"
	EmailLogStatusDelivered EmailLogStatus = "delivered"
)

func (s EmailLogStatus) Validate() error {
	switch s {
	case EmailLogStatusPending, EmailLogStatusSent, EmailLogStatusFailed, EmailLogStatusDelivered:
		return nil
	default:
		return errors.New("invalid email log status")
	}
}

type CreateEmailLog struct {
	ID             uuid.UUID      `json:"id"`
	TemplateID     uuid.UUID      `json:"templateId"`
	WorkspaceID    uuid.UUID      `json:"workspaceId"`
	RecipientEmail string         `json:"recipientEmail"`
	Event          string         `json:"event"`
	Subject        string         `json:"subject"`
	Content        string         `json:"content"`
	Status         EmailLogStatus `json:"status"`
	ErrorMessage   null.String    `json:"errorMessage"`
}

func (c CreateEmailLog) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.TemplateID, validation.Required, validation.IsUUID),
		validation.Field(&c.WorkspaceID, validation.Required, validation.IsUUID),
		validation.Field(&c.RecipientEmail, validation.Required, validation.IsEmail),
		validation.Field(&c.Event, validation.Required, validation.Length(1, 100)),
		validation.Field(&c.Subject, validation.Required, validation.Length(1, 500)),
		validation.Field(&c.Content, validation.Required),
		validation.Field(&c.Status, validation.Required),
	)
}

type UpdateEmailLog struct {
	Status       EmailLogStatus `json:"status"`
	ErrorMessage null.String    `json:"errorMessage"`
}

func (u UpdateEmailLog) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Status, validation.Required),
	)
}
