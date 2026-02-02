package router

import (
	"net/http"

	"backend/infra/localconfig"
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func SetupRoutes(injector do.Injector) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		RegisterAuthRoutes(injector, e)
		RegisterUserRoutes(injector, e)
		RegisterOrganizationRoutes(injector, e)
		RegisterWorkspaceRoutes(injector, e)
		RegisterPermissionRoutes(injector, e)
		RegisterApiRouteRoutes(injector, e)
		RegisterEmailTemplateRoutes(injector, e)

		e.GET("/v1/docs", func(c echo.Context) error {
			configService := do.MustInvoke[*localconfig.ConfigService](injector)
			docsPath := configService.GetDocsPath()
			if docsPath == "" {
				docsPath = "docs"
			}
			html, err := scalargo.NewV2(
				scalargo.WithSpecDir(docsPath),
				scalargo.WithBaseFileName("openapi.yaml"),
			)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.HTML(http.StatusOK, html)
		})
	}
}
