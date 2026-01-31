package router

import (
	"backend/app/workspace_member/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterWorkspaceMemberRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/workspace-members")

	g.POST("", h.Create)
	g.PUT("", h.Update)
	g.DELETE("", h.Delete)
	g.GET("", h.FindAll)
	g.GET("", h.FindOne)
}
