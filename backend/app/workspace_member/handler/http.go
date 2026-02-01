package handler

import (
	"backend/app/workspace_member/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/httpresponse"
	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

type HTTP struct {
	svc    domain.Service
	logger basedomain.Logger
}

func NewHTTP(svc domain.Service, logger basedomain.Logger) HTTP {
	return HTTP{
		svc:    svc,
		logger: logger.With("component", "workspace_member.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.QueryParam("userId")

	if userID == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("userId is required")
	}

	criteria := dafi.Where("userId", dafi.Equal, userID)
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

	var input domain.CreateWorkspaceMember
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
	userID := c.QueryParam("userId")

	if userID == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("userId is required")
	}

	var input domain.UpdateWorkspaceMember
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	filters := dafi.FilterBy("userId", dafi.Equal, userID)
	if err := h.svc.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.NoContent(c)
}

func (h HTTP) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.QueryParam("userId")

	if userID == "" {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Errorf("userId is required")
	}

	filters := dafi.FilterBy("userId", dafi.Equal, userID)
	if err := h.svc.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.NoContent(c)
}
