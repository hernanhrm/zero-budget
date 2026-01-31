package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	basedomain "backend/domain"
	"backend/infra/di"
	"backend/infra/logger"
	"backend/infra/server"
	"feature/user"
	"feature/user/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	injector := di.New()

	di.ProvideValue[basedomain.Logger](injector, logger.NewProduction())

	// Register feature modules
	user.Module(injector)

	// Build server config
	config := server.Config{
		Port:        8080,
		ServiceName: "zero-budget-api",
	}

	log := di.MustInvoke[basedomain.Logger](injector)

	// Create server with route setup
	srv := server.NewServer(config, log, func(e *echo.Echo) {
		// Register feature routes
		// userHandler := di.MustInvoke[handler.HTTP](injector)
		// userHandler.RegisterRoutes(e.Group("/users"))
		_ = handler.HTTP{} // placeholder until database is registered
	})

	// Register health checkers
	// db := di.MustInvoke[*database.Database](injector)
	// srv.RegisterHealthChecker("database", db)

	if err := srv.Start(ctx); err != nil {
		log.Error("server error", "error", err)
	}

	if err := di.Shutdown(injector); err != nil {
		log.Error("shutdown error", "error", err)
	}
}
