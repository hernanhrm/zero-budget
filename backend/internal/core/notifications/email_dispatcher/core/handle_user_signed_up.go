package core

import (
	"context"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"
)

func (s service) HandleUserSignedUp(ctx context.Context, event eventbusPort.Event) {
	var payload events.UserSignedUpPayload

	switch p := event.Payload.(type) {
	case events.UserSignedUpPayload:
		payload = p
	case map[string]any:
		payload = events.UserSignedUpPayload{
			UserID:         getString(p, "userId"),
			Email:          getString(p, "email"),
			Name:           getString(p, "name"),
			OrganizationID: getString(p, "organizationId"),
		}
	default:
		s.logger.Error("invalid payload for event", "event", event.Name)
		return
	}

	s.sendEmail(ctx, sendEmailInput{
		event:          events.UserSignedUp,
		organizationID: payload.OrganizationID,
		recipient:      payload.Email,
		data:           payload,
	})
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
