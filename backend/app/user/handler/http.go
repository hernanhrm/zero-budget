package handler

import (
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"backend/infra/httpresponse"
	"backend/app/user/domain"

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
		logger: logger.With("component", "user.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	criteria := dafi.Where("id", dafi.Equal, id)
	user, err := h.svc.FindOne(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, user)
}

func (h HTTP) FindAll(c echo.Context) error {
	ctx := c.Request().Context()

	parser := dafi.NewQueryParser()
	criteria, err := parser.Parse(c.QueryParams())
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}

	users, err := h.svc.FindAll(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}

	return httpresponse.OK(c, users)
}

func (h HTTP) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var input domain.CreateUser
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
	id := c.Param("id")

	var input domain.UpdateUser
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
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

func (h HTTP) RegisterRoutes(g *echo.Group) {
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindOne)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
