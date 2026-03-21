package middleware

import (
	"backend/infra/httpresponse"

	"github.com/labstack/echo/v4"
)

func APIKeyAuth(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-Key")
			if key == "" || key != expectedKey {
				return httpresponse.Unauthorized(c, "invalid or missing API key")
			}
			return next(c)
		}
	}
}
