package handler

import (
	"backend/core/budget/account/accounttype"
	"backend/core/budget/account/port"
	"backend/infra/dafi"
	"backend/infra/httpresponse"
	basedomain "backend/port"
	apperrors "backend/port/errors"

	"github.com/guregu/null/v6"
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
		logger: logger.With("component", "account.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	criteria := dafi.Where("id", dafi.Equal, id)
	acct, err := h.svc.FindOne(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, acct)
}

func (h HTTP) FindAll(c echo.Context) error {
	ctx := c.Request().Context()

	parser := dafi.NewQueryParser()
	criteria, err := parser.Parse(c.QueryParams())
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	accts, err := h.svc.FindAll(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, accts)
}

func (h HTTP) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var input port.CreateAccount
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	input.Type = accounttype.Normalize(input.Type)

	if err := h.svc.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.Created(c, nil)
}

func (h HTTP) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var input port.UpdateAccount
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	if input.Type.Valid {
		input.Type = null.StringFrom(accounttype.Normalize(input.Type.String))
	}

	filters := dafi.FilterBy("id", dafi.Equal, id)
	if err := h.svc.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.NoContent(c)
}

func (h HTTP) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	filters := dafi.FilterBy("id", dafi.Equal, id)
	if err := h.svc.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.NoContent(c)
}
