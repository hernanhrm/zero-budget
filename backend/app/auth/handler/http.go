package handler

import (
	"backend/app/auth/domain"
	"backend/infra/httpresponse"
	"github.com/labstack/echo/v4"
)

type HTTP struct {
	svc domain.Service
}

func NewHTTP(svc domain.Service) HTTP {
	return HTTP{svc: svc}
}

func (h HTTP) Signup(c echo.Context) error {
	ctx := c.Request().Context()

	var input domain.SignupWithEmail
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

	var input domain.LoginWithEmail
	if err := c.Bind(&input); err != nil {
		return err
	}

	result, err := h.svc.LoginWithEmail(ctx, input)
	if err != nil {
		return err
	}

	return httpresponse.OK(c, result)
}
