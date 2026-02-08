package router

import (
	"backend/core/email_log/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterEmailLogRoutes(injector do.Injector, g *echo.Group) {
	h := di.MustInvoke[handler.HTTP](injector)

	emailLogsGroup := g.Group("/:id/logs")
	emailLogsGroup.GET("", h.FindByTemplate)
}
