package router

import (
	"api/handler"
	"api/middleware"

	"backend/adapter/di"
	"backend/adapter/localconfig"
	eventbusPort "backend/core/notifications/eventbus/port"
	basedomain "backend/port"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterEventRoutes(injector do.Injector, e *echo.Echo) {
	bus := di.MustInvoke[eventbusPort.EventBus](injector)
	logger := di.MustInvoke[basedomain.Logger](injector)
	cfg := di.MustInvoke[localconfig.LocalConfig](injector)

	h := handler.NewEventHandler(bus, logger)

	eventsGroup := e.Group("/v1/events", middleware.APIKeyAuth(cfg.Identity.InternalAPIKey))
	eventsGroup.POST("/publish", h.Publish)
}
