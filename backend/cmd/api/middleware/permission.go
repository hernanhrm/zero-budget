package middleware

import (
	"fmt"

	"backend/infra/httpresponse"

	"github.com/labstack/echo/v4"
)

// DefaultCRUDActions maps HTTP methods to standard CRUD permission actions.
var DefaultCRUDActions = map[string]string{
	"POST":   "create",
	"GET":    "read",
	"PUT":    "update",
	"DELETE": "delete",
}

// ReadOnlyActions maps only GET to the read permission action.
var ReadOnlyActions = map[string]string{
	"GET": "read",
}

// ResourceRule defines which Better Auth resource a route pattern maps to.
type ResourceRule struct {
	Resource string            // Better Auth resource name, e.g. "emailTemplate"
	Actions  map[string]string // HTTP method → action override; nil uses DefaultCRUDActions
}

// PathResources maps Echo route patterns to resource rules.
type PathResources map[string]ResourceRule

// RequirePermission returns a global middleware that checks permissions
// against the identity service based on the matched route pattern.
func RequirePermission(client *PermissionClient, resources PathResources) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rule, ok := resources[c.Path()]
			if !ok {
				return next(c)
			}

			actions := rule.Actions
			if actions == nil {
				actions = DefaultCRUDActions
			}

			action, ok := actions[c.Request().Method]
			if !ok {
				return next(c)
			}

			permission := fmt.Sprintf("%s:%s", rule.Resource, action)

			hasPermission, err := client.HasPermission(c.Request().Context(), permission, c.Request().Header)
			if err != nil {
				return err
			}

			if !hasPermission {
				return httpresponse.Forbidden(c, "insufficient permissions")
			}

			return next(c)
		}
	}
}
