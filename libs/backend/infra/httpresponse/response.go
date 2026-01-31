package httpresponse

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a successful API response.
type Response struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

// OK sends a 200 OK response with the given data.
func OK(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, Response{Data: data})
}

// Created sends a 201 Created response with the given data.
func Created(c echo.Context, data any) error {
	return c.JSON(http.StatusCreated, Response{Data: data})
}

// NoContent sends a 204 No Content response.
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// WithMeta adds metadata to a response and sends it.
func WithMeta(c echo.Context, status int, data, meta any) error {
	return c.JSON(status, Response{Data: data, Meta: meta})
}

// Problem sends an RFC 9457 problem detail response.
func Problem(c echo.Context, problem ProblemDetail) error {
	c.Response().Header().Set(echo.HeaderContentType, ContentTypeProblemJSON)
	return c.JSON(problem.Status, problem)
}

// BadRequest sends a 400 Bad Request problem response.
func BadRequest(c echo.Context, detail string) error {
	return Problem(c, ErrBadRequest.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// NotFound sends a 404 Not Found problem response.
func NotFound(c echo.Context, detail string) error {
	return Problem(c, ErrNotFound.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// Unauthorized sends a 401 Unauthorized problem response.
func Unauthorized(c echo.Context, detail string) error {
	return Problem(c, ErrUnauthorized.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// Forbidden sends a 403 Forbidden problem response.
func Forbidden(c echo.Context, detail string) error {
	return Problem(c, ErrForbidden.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// Conflict sends a 409 Conflict problem response.
func Conflict(c echo.Context, detail string) error {
	return Problem(c, ErrConflict.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// UnprocessableEntity sends a 422 Unprocessable Entity problem response.
func UnprocessableEntity(c echo.Context, detail string) error {
	return Problem(c, ErrUnprocessableEntity.WithDetail(detail).WithInstance(c.Request().URL.Path))
}

// InternalServerError sends a 500 Internal Server Error problem response.
func InternalServerError(c echo.Context, detail string) error {
	return Problem(c, ErrInternalServerError.WithDetail(detail).WithInstance(c.Request().URL.Path))
}
