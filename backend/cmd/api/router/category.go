package router

import (
	"backend/core/budget/category/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterCategoryRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/v1/categories")

	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindOne)
}
