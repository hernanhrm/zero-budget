package router

import (
	"backend/core/user/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterUserRoutes(injector do.Injector, e *echo.Echo) {
	userHandler := di.MustInvoke[handler.HTTP](injector)

	usersGroup := e.Group("/v1/users")

	usersGroup.POST("", userHandler.Create)
	usersGroup.PUT("/:id", userHandler.Update)
	usersGroup.DELETE("/:id", userHandler.Delete)
	usersGroup.GET("", userHandler.FindAll)
	usersGroup.GET("/one", userHandler.FindOne)
}
