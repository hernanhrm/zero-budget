package core

import (
	"context"

	"backend/core/notifications/email_dispatcher/port"
	basedomain "backend/port"

	"github.com/resend/resend-go/v3"
	"github.com/samber/oops"
)

type resendSender struct {
	client      *resend.Client
	fromAddress string
	logger      basedomain.Logger
}

func NewResendSender(apiKey, fromAddress string, logger basedomain.Logger) port.EmailSender {
	return resendSender{
		client:      resend.NewClient(apiKey),
		fromAddress: fromAddress,
		logger:      logger.With("component", "email_dispatcher.resend_sender"),
	}
}

func (s resendSender) Send(ctx context.Context, msg port.EmailMessage) error {
	params := &resend.SendEmailRequest{
		From:    s.fromAddress,
		To:      []string{msg.To},
		Subject: msg.Subject,
		Html:    msg.Content,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return oops.Wrapf(err, "failed to send email via Resend to %s", msg.To)
	}

	s.logger.WithContext(ctx).Info("email sent via Resend",
		"to", msg.To,
		"subject", msg.Subject,
		"message_id", sent.Id,
	)

	return nil
}
