package handler

import (
	"backend/core/auth/port"
	"backend/infra/httpresponse"
	"github.com/labstack/echo/v4"
	"strings"
)

type HTTP struct {
	svc port.Service
}

func NewHTTP(svc port.Service) HTTP {
	return HTTP{svc: svc}
}

func (h HTTP) Signup(c echo.Context) error {
	ctx := c.Request().Context()

	var input port.SignupWithEmail
	if err := c.Bind(&input); err != nil {
		return err
	}

	result, err := h.svc.SignupWithEmail(ctx, input)
	if err != nil {
		return err
	}

	return httpresponse.Created(c, result)
}

func (h HTTP) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var input port.LoginWithEmail
	if err := c.Bind(&input); err != nil {
		return err
	}

	result, err := h.svc.LoginWithEmail(ctx, input)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, result)
}

func (h HTTP) Refresh(c echo.Context) error {
	ctx := c.Request().Context()

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return httpresponse.Unauthorized(c, "Missing or invalid Authorization header")
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	var input port.RefreshTokenRequest
	if err := c.Bind(&input); err != nil {
		return httpresponse.Unauthorized(c, "Invalid refresh token")
	}

	result, err := h.svc.RefreshToken(ctx, accessToken, input.RefreshToken)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, result)
}
