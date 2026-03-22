package core

import (
	"context"
	"fmt"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"
)

func (s service) HandleUserPasswordReset(ctx context.Context, event eventbusPort.Event) {
	s.logger.Info("HandleUserPasswordReset invoked",
		"eventName", event.Name,
		"payloadType", fmt.Sprintf("%T", event.Payload),
	)

	var payload events.UserPasswordResetPayload

	switch p := event.Payload.(type) {
	case events.UserPasswordResetPayload:
		payload = p
	case map[string]any:
		payload = events.UserPasswordResetPayload{
			UserID:   getString(p, "userId"),
			Email:    getString(p, "email"),
			Name:     getString(p, "name"),
			ResetURL: getString(p, "resetUrl"),
		}
	default:
		s.logger.Error("invalid payload for event", "event", event.Name, "payloadType", fmt.Sprintf("%T", event.Payload))
		return
	}

	s.logger.Info("parsed password reset payload",
		"userId", payload.UserID,
		"email", payload.Email,
		"name", payload.Name,
		"hasResetURL", payload.ResetURL != "",
	)

	s.sendEmail(ctx, sendEmailInput{
		event:     events.UserPasswordReset,
		recipient: payload.Email,
		data:      payload,
	})
}
