package router

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func SetupRoutes(injector do.Injector) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		RegisterUserRoutes(injector, e)

		e.GET("/docs", func(c echo.Context) error {
			html, err := scalargo.NewV2(
				scalargo.WithSpecDir("docs"),
				scalargo.WithBaseFileName("openapi.yaml"),
			)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.HTML(http.StatusOK, html)
		})
	}
}
