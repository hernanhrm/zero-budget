package router

import (
	"backend/core/workspace/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterWorkspaceRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	workspaceGroup := e.Group("/v1/workspaces")

	workspaceGroup.POST("", h.Create)
	workspaceGroup.PUT("/:slug", h.Update)
	workspaceGroup.DELETE("/:slug", h.Delete)
	workspaceGroup.GET("", h.FindAll)
	workspaceGroup.GET("/:slug", h.FindOne)

	RegisterWorkspaceMemberRoutes(injector, workspaceGroup)
	RegisterRoleRoutes(injector, workspaceGroup)
}
