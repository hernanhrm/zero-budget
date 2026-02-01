package router

import (
	"backend/app/workspace/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterWorkspaceRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	workspaceGroup := e.Group("/workspaces")

	workspaceGroup.POST("", h.Create)
	workspaceGroup.PUT("/:slug", h.Update)
	workspaceGroup.DELETE("/:slug", h.Delete)
	workspaceGroup.GET("", h.FindAll)
	workspaceGroup.GET("/:slug", h.FindOne)

	RegisterWorkspaceMemberRoutes(injector, workspaceGroup)
}
