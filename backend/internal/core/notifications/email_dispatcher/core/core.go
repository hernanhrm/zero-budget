package core

import (
	"context"

	emailLogPort "backend/core/notifications/email_log/port"
	emailTemplatePort "backend/core/notifications/email_template/port"
	"backend/core/notifications/email_dispatcher/port"
	"backend/infra/dafi"
	basedomain "backend/port"
	apperrors "backend/port/errors"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/samber/oops"
)

type service struct {
	emailTemplateSvc emailTemplatePort.Service
	emailLogSvc      emailLogPort.Service
	emailSender      port.EmailSender
	logger           basedomain.Logger
}

func New(
	emailTemplateSvc emailTemplatePort.Service,
	emailLogSvc emailLogPort.Service,
	emailSender port.EmailSender,
	logger basedomain.Logger,
) port.Service {
	return service{
		emailTemplateSvc: emailTemplateSvc,
		emailLogSvc:      emailLogSvc,
		emailSender:      emailSender,
		logger:           logger.With("component", "email_dispatcher.service"),
	}
}

type sendEmailInput struct {
	event          string
	organizationID string
	recipient      string
	data           any
}

func (s service) sendEmail(ctx context.Context, input sendEmailInput) {
	criteria := dafi.Where("event", dafi.Equal, input.event).
		And("isActive", dafi.Equal, true)

	if input.organizationID != "" {
		criteria = criteria.And("organizationId", dafi.Equal, input.organizationID)
	} else {
		criteria = criteria.And("organizationId", dafi.IsNull, nil)
	}

	tmpl, err := s.emailTemplateSvc.FindOne(ctx, criteria)
	if err != nil {
		s.logger.Error("failed to find email template",
			"event", input.event,
			"organizationId", input.organizationID,
			"error", err,
		)
		return
	}

	parsed, err := s.emailTemplateSvc.ParseTemplate(tmpl, input.data)
	if err != nil {
		s.logger.Error("failed to parse email template",
			"event", input.event,
			"templateId", tmpl.ID,
			"error", err,
		)
		return
	}

	logID := uuid.New()

	sendErr := s.emailSender.Send(ctx, port.EmailMessage{
		To:      input.recipient,
		Subject: parsed.Subject,
		Content: parsed.Content,
	})

	status := emailLogPort.EmailLogStatusSent
	var errorMessage null.String
	if sendErr != nil {
		status = emailLogPort.EmailLogStatusFailed
		errorMessage = null.StringFrom(sendErr.Error())
		s.logger.Error("failed to send email",
			"event", input.event,
			"recipient", input.recipient,
			"error", sendErr,
		)
	}

	createLog := emailLogPort.CreateEmailLog{
		ID:             logID,
		TemplateID:     tmpl.ID,
		OrganizationID: null.NewString(input.organizationID, input.organizationID != ""),
		RecipientEmail: input.recipient,
		Event:          input.event,
		Subject:        parsed.Subject,
		Content:        parsed.Content,
		Status:         status,
		ErrorMessage:   errorMessage,
	}

	if err := s.emailLogSvc.Create(ctx, createLog); err != nil {
		s.logger.Error("failed to create email log",
			"event", input.event,
			"error", oops.In(apperrors.LayerService).Wrap(err),
		)
	}
}
