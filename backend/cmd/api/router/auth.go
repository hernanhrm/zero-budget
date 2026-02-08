package router

import (
	"backend/core/auth/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterAuthRoutes(injector do.Injector, e *echo.Echo) {
	authHandler := di.MustInvoke[handler.HTTP](injector)

	authGroup := e.Group("/v1/auth")
	authGroup.POST("/signup", authHandler.Signup)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)
}
