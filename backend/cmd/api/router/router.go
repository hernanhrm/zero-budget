package router

import (
	"net/http"

	"api/middleware"

	"backend/adapter/localconfig"
	"backend/adapter/di"
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func SetupRoutes(injector do.Injector) func(e *echo.Echo) {
	return func(e *echo.Echo) {
		cfg := di.MustInvoke[localconfig.LocalConfig](injector)
		permClient := middleware.NewPermissionClient(cfg.Identity.URL)

		e.Use(middleware.RequirePermission(permClient, middleware.PathResources{
			"/v1/email-templates":          {Resource: "emailTemplate"},
			"/v1/email-templates/:id":      {Resource: "emailTemplate"},
			"/v1/email-templates/:id/logs": {Resource: "emailLog", Actions: middleware.ReadOnlyActions},
		}))

		RegisterEmailTemplateRoutes(injector, e)
		RegisterEventRoutes(injector, e)

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
