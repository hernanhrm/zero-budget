package router

import (
	"backend/app/email_template/handler"
	"backend/infra/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterEmailTemplateRoutes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	emailTemplatesGroup := e.Group("/email-templates")

	emailTemplatesGroup.POST("", h.Create)
	emailTemplatesGroup.PUT("/:id", h.Update)
	emailTemplatesGroup.DELETE("/:id", h.Delete)
	emailTemplatesGroup.GET("", h.FindAll)
	emailTemplatesGroup.GET("/:id", h.FindOne)

	RegisterEmailLogRoutes(injector, emailTemplatesGroup)
}
