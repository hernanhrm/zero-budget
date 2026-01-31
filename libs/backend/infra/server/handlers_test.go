package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/infra/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHealth(t *testing.T) {
	s := &Server{
		logger:   logger.NewNoop(),
		checkers: make(map[string]HealthChecker),
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := s.handleHealth(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response HealthResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, healthStatusHealthy, response.Status)
}

func TestHandlePing(t *testing.T) {
	s := &Server{
		logger: logger.NewNoop(),
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := s.handlePing(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response PingResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response.Message)
	assert.NotEmpty(t, response.Time)
}

func TestHealthResponse_Structure(t *testing.T) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    "1h30m45s",
		Services: map[string]string{
			"database": "healthy",
			"logger":   "healthy",
			"config":   "healthy",
		},
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled HealthResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response, unmarshaled)
}

func TestPingResponse_Structure(t *testing.T) {
	response := PingResponse{
		Message: "pong",
		Time:    "2023-12-14T10:30:00Z",
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled PingResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response, unmarshaled)
}

func TestHealthStatus_Constants(t *testing.T) {
	assert.Equal(t, "healthy", healthStatusHealthy)
	assert.Equal(t, "unhealthy", healthStatusUnhealthy)
	assert.Equal(t, "degraded", healthStatusDegraded)
}
