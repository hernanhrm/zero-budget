package handler

import (
	"backend/core/budget/currency/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/infra/httpresponse"
	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

type HTTP struct {
	svc    port.Service
	logger basedomain.Logger
}

func NewHTTP(svc port.Service, logger basedomain.Logger) HTTP {
	return HTTP{
		svc:    svc,
		logger: logger.With("component", "currency.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	code := c.Param("code")

	criteria := dafi.Where("code", dafi.Equal, code)
	currency, err := h.svc.FindOne(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, currency)
}

func (h HTTP) FindAll(c echo.Context) error {
	ctx := c.Request().Context()

	parser := dafi.NewQueryParser()
	criteria, err := parser.Parse(c.QueryParams())
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	currencies, err := h.svc.FindAll(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, currencies)
}
