package port

import (
	"context"

	eventbusPort "backend/core/notifications/eventbus/port"
)

type Service interface {
	HandleUserSignedUp(ctx context.Context, event eventbusPort.Event)
	HandleUserVerificationEmail(ctx context.Context, event eventbusPort.Event)
	HandleUserPasswordReset(ctx context.Context, event eventbusPort.Event)
	HandleOrganizationInvitationCreated(ctx context.Context, event eventbusPort.Event)
}
