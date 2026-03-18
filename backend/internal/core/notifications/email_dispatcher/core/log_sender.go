package core

import (
	"context"

	"backend/core/notifications/email_dispatcher/port"
	basedomain "backend/port"
)

type logSender struct {
	logger basedomain.Logger
}

func NewLogSender(logger basedomain.Logger) port.EmailSender {
	return logSender{
		logger: logger.With("component", "email_dispatcher.log_sender"),
	}
}

func (s logSender) Send(ctx context.Context, msg port.EmailMessage) error {
	s.logger.WithContext(ctx).Info("email sent (log only)",
		"to", msg.To,
		"subject", msg.Subject,
	)
	return nil
}
