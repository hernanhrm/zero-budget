package httpresponse

import (
	"context"
	"errors"
	"maps"
	"net/http"

	"backend/domain"
	apperrors "backend/domain/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

// ErrorCodeMapping maps oops error codes to HTTP status codes.
var ErrorCodeMapping = map[string]int{
	apperrors.CodeBadRequest:    http.StatusBadRequest,
	apperrors.CodeUnauthorized:  http.StatusUnauthorized,
	apperrors.CodeForbidden:     http.StatusForbidden,
	apperrors.CodeNotFound:      http.StatusNotFound,
	apperrors.CodeConflict:      http.StatusConflict,
	apperrors.CodeAlreadyExists: http.StatusConflict,
	apperrors.CodeValidation:    http.StatusUnprocessableEntity,
}

// HTTPErrorHandlerConfig holds configuration for the HTTP error handler.
type HTTPErrorHandlerConfig struct {
	Logger       domain.Logger
	Debug        bool
	CodeMappings map[string]int
}

// HTTPErrorHandler returns an Echo HTTPErrorHandler that converts errors to RFC 9457 responses.
func HTTPErrorHandler(config HTTPErrorHandlerConfig) echo.HTTPErrorHandler {
	mappings := ErrorCodeMapping
	if config.CodeMappings != nil {
		maps.Copy(mappings, config.CodeMappings)
	}

	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		ctx := c.Request().Context()
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)

		problem, status := errorToProblem(err, c, mappings, config)

		logError(ctx, config.Logger, problem, status, requestID)

		if err := Problem(c, problem); err != nil && config.Logger != nil {
			config.Logger.Error("failed to send error response", "error", err)
		}
	}
}

func errorToProblem(err error, c echo.Context, mappings map[string]int, config HTTPErrorHandlerConfig) (ProblemDetail, int) {
	instance := c.Request().URL.Path

	var problem ProblemDetail
	if errors.As(err, &problem) {
		return problem.WithInstance(instance), problem.Status
	}

	var validationErrs validation.Errors
	if errors.As(err, &validationErrs) {
		return validationErrorToProblem(validationErrs, instance), http.StatusUnprocessableEntity
	}

	oopsErr, ok := oops.AsOops(err)
	if ok {
		return oopsErrorToProblem(oopsErr, instance, mappings, config)
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return echoErrorToProblem(echoErr, instance), echoErr.Code
	}

	problem = ErrInternalServerError.WithInstance(instance)
	if config.Debug {
		problem = problem.WithDetail(err.Error())
	}

	return problem, http.StatusInternalServerError
}

func validationErrorToProblem(validationErrs validation.Errors, instance string) ProblemDetail {
	fieldErrors := make(map[string][]string)

	for field, err := range validationErrs {
		if innerErrs, ok := err.(validation.Errors); ok {
			for innerField, innerErr := range innerErrs {
				key := field + "." + innerField
				fieldErrors[key] = append(fieldErrors[key], innerErr.Error())
			}
		} else {
			fieldErrors[field] = append(fieldErrors[field], err.Error())
		}
	}

	return NewProblemDetail(
		TypeUnprocessableEntity,
		"Validation Error",
		http.StatusUnprocessableEntity,
	).
		WithDetail(validationErrs.Error()).
		WithInstance(instance).
		WithExtension("errors", fieldErrors)
}

func oopsErrorToProblem(oopsErr oops.OopsError, instance string, mappings map[string]int, config HTTPErrorHandlerConfig) (ProblemDetail, int) {
	codeAny := oopsErr.Code()
	code := ""
	if codeAny != nil {
		code, _ = codeAny.(string)
	}

	status := http.StatusInternalServerError
	if code != "" {
		if mappedStatus, ok := mappings[code]; ok {
			status = mappedStatus
		}
	}

	var validationErrs validation.Errors
	if errors.As(oopsErr, &validationErrs) {
		problem := validationErrorToProblem(validationErrs, instance)
		problem = problem.WithExtension("code", code)
		return problem, http.StatusUnprocessableEntity
	}

	problem := NewProblemDetail(
		problemTypeFromStatus(status),
		http.StatusText(status),
		status,
	).
		WithInstance(instance).
		WithExtension("code", code)

	if config.Debug || status < 500 {
		problem = problem.WithDetail(oopsErr.Error())
	} else {
		problem = problem.WithDetail("An unexpected error occurred")
	}

	return problem, status
}

func echoErrorToProblem(echoErr *echo.HTTPError, instance string) ProblemDetail {
	status := echoErr.Code
	detail := ""

	if msg, ok := echoErr.Message.(string); ok {
		detail = msg
	}

	return NewProblemDetail(
		problemTypeFromStatus(status),
		http.StatusText(status),
		status,
	).
		WithDetail(detail).
		WithInstance(instance)
}

func logError(ctx context.Context, logger domain.Logger, problem ProblemDetail, status int, requestID string) {
	if logger == nil {
		return
	}

	attrs := []any{
		"request_id", requestID,
		"status", status,
		"type", problem.Type,
		"title", problem.Title,
	}

	if problem.Detail != "" {
		attrs = append(attrs, "detail", problem.Detail)
	}

	if code, ok := problem.Extensions["code"]; ok {
		attrs = append(attrs, "code", code)
	}

	if status >= 500 {
		logger.WithContext(ctx).Error(problem.Title, attrs...)
		return
	}

	logger.WithContext(ctx).Warn(problem.Title, attrs...)
}

// Logger interface for error logging.
// Deprecated: Use HTTPErrorHandler with domain.Logger instead.
type Logger interface {
	Error(msg string, keysAndValues ...any)
}

// MiddlewareConfig holds configuration for the error middleware.
// Deprecated: Use HTTPErrorHandlerConfig instead.
type MiddlewareConfig struct {
	Logger       Logger
	Debug        bool
	CodeMappings map[string]int
}

// ErrorMiddleware returns an Echo middleware that converts errors to RFC 9457 responses.
func ErrorMiddleware(config MiddlewareConfig) echo.MiddlewareFunc {
	mappings := ErrorCodeMapping
	if config.CodeMappings != nil {
		maps.Copy(mappings, config.CodeMappings)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			return handleError(c, err, mappings, config)
		}
	}
}

func handleError(c echo.Context, err error, mappings map[string]int, config MiddlewareConfig) error {
	// Check if it's already a ProblemDetail
	var problem ProblemDetail
	if errors.As(err, &problem) {
		return Problem(c, problem.WithInstance(c.Request().URL.Path))
	}

	// Check if it's an oops error
	oopsErr, ok := oops.AsOops(err)
	if ok {
		return handleOopsError(c, oopsErr, mappings, config)
	}

	// Check if it's an Echo HTTP error
	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return handleEchoError(c, echoErr)
	}

	// Default to internal server error
	if config.Logger != nil {
		config.Logger.Error("unhandled error", "error", err)
	}

	problem = ErrInternalServerError.WithInstance(c.Request().URL.Path)
	if config.Debug {
		problem = problem.WithDetail(err.Error())
	}

	return Problem(c, problem)
}

func handleOopsError(c echo.Context, oopsErr oops.OopsError, mappings map[string]int, config MiddlewareConfig) error {
	codeAny := oopsErr.Code()
	code := ""
	if codeAny != nil {
		code, _ = codeAny.(string)
	}
	status := http.StatusInternalServerError

	if code != "" {
		if mappedStatus, ok := mappings[code]; ok {
			status = mappedStatus
		}
	}

	if config.Logger != nil && status >= 500 {
		config.Logger.Error("server error",
			"code", code,
			"error", oopsErr.Error(),
		)
	}

	problem := NewProblemDetail(
		problemTypeFromStatus(status),
		http.StatusText(status),
		status,
	).
		WithDetail(oopsErr.Error()).
		WithInstance(c.Request().URL.Path).
		WithExtension("code", code)

	if !config.Debug && status >= 500 {
		problem.Detail = "An unexpected error occurred"
	}

	return Problem(c, problem)
}

func handleEchoError(c echo.Context, echoErr *echo.HTTPError) error {
	status := echoErr.Code
	detail := ""

	if msg, ok := echoErr.Message.(string); ok {
		detail = msg
	}

	problem := NewProblemDetail(
		problemTypeFromStatus(status),
		http.StatusText(status),
		status,
	).
		WithDetail(detail).
		WithInstance(c.Request().URL.Path)

	return Problem(c, problem)
}

func problemTypeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return TypeBadRequest
	case http.StatusUnauthorized:
		return TypeUnauthorized
	case http.StatusForbidden:
		return TypeForbidden
	case http.StatusNotFound:
		return TypeNotFound
	case http.StatusMethodNotAllowed:
		return TypeMethodNotAllowed
	case http.StatusConflict:
		return TypeConflict
	case http.StatusUnprocessableEntity:
		return TypeUnprocessableEntity
	case http.StatusServiceUnavailable:
		return TypeServiceUnavailable
	default:
		return TypeInternalServerError
	}
}
