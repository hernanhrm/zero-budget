package router

import (
	"backend/app/role/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterRoleRoutes(injector do.Injector, g *echo.Group) {
	h := di.MustInvoke[handler.HTTP](injector)

	rolesGroup := g.Group("/:slug/roles")

	rolesGroup.POST("", h.Create)
	rolesGroup.PUT("/:id", h.Update)
	rolesGroup.DELETE("/:id", h.Delete)
	rolesGroup.GET("", h.FindAll)
	rolesGroup.GET("/:id", h.FindOne)
}
