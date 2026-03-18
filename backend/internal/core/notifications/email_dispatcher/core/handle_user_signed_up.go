package core

import (
	"context"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"

	"github.com/google/uuid"
)

func (s service) HandleUserSignedUp(ctx context.Context, event eventbusPort.Event) {
	payload, ok := event.Payload.(events.UserSignedUpPayload)
	if !ok {
		s.logger.Error("invalid payload for event", "event", event.Name)
		return
	}

	workspaceID, err := uuid.Parse(payload.WorkspaceID)
	if err != nil {
		s.logger.Error("invalid workspace ID", "workspaceId", payload.WorkspaceID, "error", err)
		return
	}

	s.sendEmail(ctx, sendEmailInput{
		event:       events.UserSignedUp,
		workspaceID: workspaceID,
		recipient:   payload.Email,
		data:        payload,
	})
}
