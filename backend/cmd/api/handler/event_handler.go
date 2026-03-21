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

	h.logger.Info("received publish request",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
	)

	var req publishEventRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		return oops.In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	h.logger.Info("parsed event request",
		"event", req.Event,
		"payloadKeys", payloadKeys(req.Payload),
	)

	if req.Event == "" {
		h.logger.Warn("event name is empty, rejecting request")
		return httpresponse.BadRequest(c, "event name is required")
	}

	h.bus.Publish(ctx, eventbusPort.Event{
		Name:       req.Event,
		Payload:    req.Payload,
		OccurredAt: time.Now(),
	})

	h.logger.Info("event published successfully", "event", req.Event)

	return httpresponse.NoContent(c)
}

func payloadKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
