package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"backend/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/oops"
)

// Config holds server configuration.
type Config struct {
	Port        int
	ServiceName string
}

// HealthChecker is an interface for services that can report their health.
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// Server represents the HTTP server with all its dependencies.
type Server struct {
	config   Config
	logger   domain.Logger
	checkers map[string]HealthChecker

	Echo *echo.Echo
}

// NewServer creates a new HTTP server instance.
func NewServer(config Config, log domain.Logger, routeSetupFunc func(*echo.Echo)) *Server {
	s := &Server{
		config:   config,
		logger:   log.With("component", "server"),
		checkers: make(map[string]HealthChecker),
		Echo:     echo.New(),
	}

	s.Echo.HideBanner = true
	s.Echo.HidePort = true

	s.setupMiddleware()
	s.setupBaseRoutes()

	if routeSetupFunc != nil {
		routeSetupFunc(s.Echo)
	}

	s.logger.Info("HTTP server initialized",
		"port", config.Port,
		"service", config.ServiceName,
	)

	return s
}

// RegisterHealthChecker adds a health checker for a named service.
func (s Server) RegisterHealthChecker(name string, checker HealthChecker) {
	s.checkers[name] = checker
}

// setupMiddleware configures all HTTP middleware.
func (s Server) setupMiddleware() {
	s.Echo.Use(middleware.Recover())

	s.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		LogLatency:  true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				s.logger.Info("request completed",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
				)
			} else {
				s.logger.Error("request failed",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"error", v.Error,
				)
			}
			return nil
		},
	}))

	s.Echo.Use(middleware.CORS())
}

// setupBaseRoutes configures health and ping endpoints.
func (s Server) setupBaseRoutes() {
	s.Echo.GET("/health", s.handleHealth)
	s.Echo.GET("/ping", s.handlePing)
}

// Start starts the HTTP server and blocks until context is cancelled.
func (s Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.config.Port)

	s.logger.Info("starting HTTP server",
		"address", addr,
		"service", s.config.ServiceName,
	)

	go func() {
		if err := s.Echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()

	s.logger.Info("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(shutdownCtx); err != nil {
		return oops.
			Code("server_shutdown_failed").
			Wrapf(err, "failed to shutdown HTTP server")
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")

	if err := s.Echo.Shutdown(ctx); err != nil {
		return oops.
			Code("server_shutdown_failed").
			Wrapf(err, "failed to shutdown HTTP server")
	}

	return nil
}
