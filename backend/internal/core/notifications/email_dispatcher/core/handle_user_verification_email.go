package core

import (
	"context"
	"fmt"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"
)

func (s service) HandleUserVerificationEmail(ctx context.Context, event eventbusPort.Event) {
	s.logger.Info("HandleUserVerificationEmail invoked",
		"eventName", event.Name,
		"payloadType", fmt.Sprintf("%T", event.Payload),
	)

	var payload events.UserVerificationEmailPayload

	switch p := event.Payload.(type) {
	case events.UserVerificationEmailPayload:
		payload = p
	case map[string]any:
		payload = events.UserVerificationEmailPayload{
			UserID:          getString(p, "userId"),
			Email:           getString(p, "email"),
			Name:            getString(p, "name"),
			VerificationURL: getString(p, "verificationUrl"),
		}
	default:
		s.logger.Error("invalid payload for event", "event", event.Name, "payloadType", fmt.Sprintf("%T", event.Payload))
		return
	}

	s.logger.Info("parsed verification email payload",
		"userId", payload.UserID,
		"email", payload.Email,
		"name", payload.Name,
		"hasVerificationURL", payload.VerificationURL != "",
	)

	s.sendEmail(ctx, sendEmailInput{
		event:     events.UserVerificationEmail,
		recipient: payload.Email,
		data:      payload,
	})
}
