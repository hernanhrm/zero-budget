package router

import (
	"backend/core/budget/currency/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterCurrencyRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/v1/currencies")

	g.GET("", h.FindAll)
	g.GET("/:code", h.FindOne)
}
