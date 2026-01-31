package router

import (
	"backend/app/workspace/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterWorkspaceRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/workspaces")

	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindOne)
}
