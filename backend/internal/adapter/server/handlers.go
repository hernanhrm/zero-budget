package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

const (
	healthStatusHealthy   = "healthy"
	healthStatusUnhealthy = "unhealthy"
	healthStatusDegraded  = "degraded"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Time     string            `json:"time"`
	Services map[string]string `json:"services"`
}

// PingResponse represents the ping response.
type PingResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

// handleHealth returns detailed health status of all registered services.
func (s Server) handleHealth(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	for name, checker := range s.checkers {
		if err := checker.HealthCheck(ctx); err != nil {
			services[name] = healthStatusUnhealthy + ": " + err.Error()
		} else {
			services[name] = healthStatusHealthy
		}
	}

	status := healthStatusHealthy
	for _, svcStatus := range services {
		if svcStatus != healthStatusHealthy {
			status = healthStatusDegraded
			break
		}
	}

	response := HealthResponse{
		Status:   status,
		Version:  "0.1.0",
		Time:     time.Now().UTC().Format(time.RFC3339),
		Services: services,
	}

	httpStatus := http.StatusOK
	if status != healthStatusHealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	if err := c.JSON(httpStatus, response); err != nil {
		return oops.
			Code("health_response_failed").
			Wrapf(err, "failed to send health response")
	}
	return nil
}

// handlePing is a simple ping endpoint.
func (s Server) handlePing(c echo.Context) error {
	response := PingResponse{
		Message: "pong",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := c.JSON(http.StatusOK, response); err != nil {
		return oops.
			Code("ping_response_failed").
			Wrapf(err, "failed to send ping response")
	}
	return nil
}
