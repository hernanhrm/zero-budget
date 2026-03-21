package handler

import (
	"time"

	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/infra/httpresponse"
	basedomain "backend/port"
	apperrors "backend/port/errors"

	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

type publishEventRequest struct {
	Event   string         `json:"event"`
	Payload map[string]any `json:"payload"`
}

type EventHandler struct {
	bus    eventbusPort.EventBus
	logger basedomain.Logger
}

func NewEventHandler(bus eventbusPort.EventBus, logger basedomain.Logger) EventHandler {
	return EventHandler{
		bus:    bus,
		logger: logger.With("component", "event.handler"),
	}
}

func (h EventHandler) Publish(c echo.Context) error {
	ctx := c.Request().Context()

	var req publishEventRequest
	if err := c.Bind(&req); err != nil {
		return oops.In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	if req.Event == "" {
		return httpresponse.BadRequest(c, "event name is required")
	}

	h.bus.Publish(ctx, eventbusPort.Event{
		Name:       req.Event,
		Payload:    req.Payload,
		OccurredAt: time.Now(),
	})

	h.logger.Info("event published", "event", req.Event)

	return httpresponse.NoContent(c)
}
