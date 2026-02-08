package router

import (
	"backend/core/api_route/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterApiRouteRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/v1/api-routes")

	g.POST("", h.Create)
	g.PUT("", h.Update)
	g.DELETE("", h.Delete)
	g.GET("", h.FindAll)
	g.GET("", h.FindOne)
}
