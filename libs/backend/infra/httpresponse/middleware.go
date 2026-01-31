package httpresponse

import (
	"errors"
	"net/http"

	apperrors "backend/domain/errors"

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

// Logger interface for error logging.
type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

// MiddlewareConfig holds configuration for the error middleware.
type MiddlewareConfig struct {
	Logger       Logger
	Debug        bool
	CodeMappings map[string]int
}

// ErrorMiddleware returns an Echo middleware that converts errors to RFC 9457 responses.
func ErrorMiddleware(config MiddlewareConfig) echo.MiddlewareFunc {
	mappings := ErrorCodeMapping
	if config.CodeMappings != nil {
		for k, v := range config.CodeMappings {
			mappings[k] = v
		}
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
