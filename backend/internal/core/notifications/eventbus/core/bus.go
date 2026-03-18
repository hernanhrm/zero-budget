package core

import (
	"context"
	"sync"
	"time"

	"backend/core/notifications/eventbus/port"
	basedomain "backend/port"
)

const defaultBufferSize = 256

type bus struct {
	handlers map[string][]port.HandlerFunc
	mu       sync.RWMutex
	events   chan port.Event
	wg       sync.WaitGroup
	logger   basedomain.Logger
	done     chan struct{}
}

func New(logger basedomain.Logger) port.EventBus {
	return &bus{
		handlers: make(map[string][]port.HandlerFunc),
		events:   make(chan port.Event, defaultBufferSize),
		logger:   logger.With("component", "eventbus"),
		done:     make(chan struct{}),
	}
}

func (b *bus) Subscribe(eventName string, handler port.HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *bus) Publish(ctx context.Context, event port.Event) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}

	select {
	case b.events <- event:
		b.logger.Info("event published", "event", event.Name)
	default:
		b.logger.Error("event bus buffer full, dropping event", "event", event.Name)
	}
}

func (b *bus) Start(ctx context.Context) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case event := <-b.events:
				b.dispatch(ctx, event)
			case <-b.done:
				// Drain remaining events
				for {
					select {
					case event := <-b.events:
						b.dispatch(ctx, event)
					default:
						return
					}
				}
			}
		}
	}()
}

func (b *bus) Shutdown() {
	close(b.done)
	b.wg.Wait()
}

func (b *bus) dispatch(ctx context.Context, event port.Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Name]
	b.mu.RUnlock()

	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					b.logger.Error("event handler panicked", "event", event.Name, "panic", r)
				}
			}()
			handler(ctx, event)
		}()
	}
}
