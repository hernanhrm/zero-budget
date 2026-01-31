package router

import (
	"backend/app/organization/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterOrganizationRoutes(injector do.Injector, e *echo.Echo) {
	organizationHandler := di.MustInvoke[handler.HTTP](injector)

	organizationsGroup := e.Group("/organizations")

	organizationsGroup.POST("", organizationHandler.Create)
	organizationsGroup.PUT("/:id", organizationHandler.Update)
	organizationsGroup.DELETE("/:id", organizationHandler.Delete)
	organizationsGroup.GET("", organizationHandler.FindAll)
	organizationsGroup.GET("/one", organizationHandler.FindOne)
}
