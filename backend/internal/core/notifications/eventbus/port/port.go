package port

import (
	"context"
	"time"
)

type Event struct {
	Name       string
	Payload    any
	OccurredAt time.Time
	Metadata   map[string]string
}

type HandlerFunc func(ctx context.Context, event Event)

type EventBus interface {
	Publish(ctx context.Context, event Event)
	Subscribe(eventName string, handler HandlerFunc)
	Start(ctx context.Context)
	Shutdown()
}
