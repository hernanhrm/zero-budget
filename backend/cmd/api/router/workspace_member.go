package router

import (
	"backend/core/workspace_member/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterWorkspaceMemberRoutes(injector do.Injector, g *echo.Group) {
	h := di.MustInvoke[handler.HTTP](injector)

	membersGroup := g.Group("/:slug/members")

	membersGroup.GET("", h.FindAll)
	membersGroup.POST("", h.Create)
	membersGroup.PUT("", h.Update)
	membersGroup.DELETE("", h.Delete)
}
