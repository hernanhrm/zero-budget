package handler

import (
	"backend/core/api_route/port"
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
		logger: logger.With("component", "api_route.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()

	// Assuming method and path are passed as query params because path can contain slashes
	method := c.QueryParam("method")
	path := c.QueryParam("path")

	if method == "" || path == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("method and path are required")
	}

	criteria := dafi.Criteria{
		Filters: dafi.FilterBy("method", dafi.Equal, method).And("path", dafi.Equal, path),
	}

	item, err := h.svc.FindOne(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.OK(c, item)
}

func (h HTTP) FindAll(c echo.Context) error {
	ctx := c.Request().Context()
	parser := dafi.NewQueryParser()
	criteria, err := parser.Parse(c.QueryParams())
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	items, err := h.svc.FindAll(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.OK(c, items)
}

func (h HTTP) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var input port.CreateApiRoute
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	if err := h.svc.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.Created(c, nil)
}

func (h HTTP) Update(c echo.Context) error {
	ctx := c.Request().Context()
	method := c.QueryParam("method")
	path := c.QueryParam("path")

	if method == "" || path == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("method and path are required")
	}

	var input port.UpdateApiRoute
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	filters := []dafi.Filter{
		{Field: "method", Operator: dafi.Equal, Value: method},
		{Field: "path", Operator: dafi.Equal, Value: path},
	}

	if err := h.svc.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.NoContent(c)
}

func (h HTTP) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	method := c.QueryParam("method")
	path := c.QueryParam("path")

	if method == "" || path == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("method and path are required")
	}

	filters := []dafi.Filter{
		{Field: "method", Operator: dafi.Equal, Value: method},
		{Field: "path", Operator: dafi.Equal, Value: path},
	}
	if err := h.svc.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.NoContent(c)
}
