package core

import (
	"context"
	"fmt"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"
)

func (s service) HandleOrganizationInvitationCreated(ctx context.Context, event eventbusPort.Event) {
	s.logger.Info("HandleOrganizationInvitationCreated invoked",
		"eventName", event.Name,
		"payloadType", fmt.Sprintf("%T", event.Payload),
	)

	var payload events.OrganizationInvitationCreatedPayload

	switch p := event.Payload.(type) {
	case events.OrganizationInvitationCreatedPayload:
		payload = p
	case map[string]any:
		payload = events.OrganizationInvitationCreatedPayload{
			Email:            getString(p, "email"),
			InviterName:      getString(p, "inviterName"),
			InviterEmail:     getString(p, "inviterEmail"),
			InviterInitial:   getString(p, "inviterInitial"),
			OrganizationName: getString(p, "organizationName"),
			AcceptURL:        getString(p, "acceptUrl"),
			DeclineURL:       getString(p, "declineUrl"),
		}
	default:
		s.logger.Error("invalid payload for event", "event", event.Name, "payloadType", fmt.Sprintf("%T", event.Payload))
		return
	}

	s.logger.Info("parsed organization invitation payload",
		"email", payload.Email,
		"inviterName", payload.InviterName,
		"organizationName", payload.OrganizationName,
		"hasAcceptURL", payload.AcceptURL != "",
	)

	s.sendEmail(ctx, sendEmailInput{
		event:     events.OrganizationInvitationCreated,
		recipient: payload.Email,
		data:      payload,
	})
}
