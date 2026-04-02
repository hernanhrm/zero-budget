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
			"/v1/currencies":               {Resource: "currency", Actions: middleware.ReadOnlyActions},
			"/v1/currencies/:code":         {Resource: "currency", Actions: middleware.ReadOnlyActions},
			"/v1/organization-currencies":      {Resource: "organizationCurrency"},
			"/v1/organization-currencies/:id":  {Resource: "organizationCurrency"},
			"/v1/accounts":                 {Resource: "account"},
			"/v1/accounts/:id":             {Resource: "account"},
			"/v1/categories":               {Resource: "category"},
			"/v1/categories/:id":           {Resource: "category"},
			"/v1/budgets":                  {Resource: "budget"},
			"/v1/budgets/:id":              {Resource: "budget"},
			"/v1/transactions":             {Resource: "transaction"},
			"/v1/transactions/:id":         {Resource: "transaction"},
		}))

		RegisterEmailTemplateRoutes(injector, e)
		RegisterEventRoutes(injector, e)
		RegisterCurrencyRoutes(injector, e)
		RegisterOrganizationCurrencyRoutes(injector, e)
		RegisterAccountRoutes(injector, e)
		RegisterCategoryRoutes(injector, e)
		RegisterBudgetRoutes(injector, e)
		RegisterTransactionRoutes(injector, e)

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
