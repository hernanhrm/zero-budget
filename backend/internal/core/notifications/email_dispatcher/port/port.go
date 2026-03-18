package port

import (
	"context"

	eventbusPort "backend/core/notifications/eventbus/port"
)

type Service interface {
	HandleUserSignedUp(ctx context.Context, event eventbusPort.Event)
}
