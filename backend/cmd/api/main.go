package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"api/router"
	"backend/core/notifications/email_log"
	"backend/core/notifications/email_template"
	"backend/core/notifications/email_dispatcher"
	"backend/core/notifications/eventbus"
	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/adapter/database"
	"backend/adapter/di"
	"backend/adapter/localconfig"
	"backend/adapter/logger"
	"backend/adapter/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	injector := di.New()

	log := logger.NewProduction()
	di.ProvideValue(injector, log)

	configSvc, err := localconfig.NewConfigService(log)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	di.ProvideValue(injector, configSvc)

	cfg := configSvc.Get()
	di.ProvideValue(injector, cfg)

	db, err := database.NewConnection(ctx, cfg.Database.URL, log)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	di.ProvideValue(injector, db)
	di.ProvideValue(injector, db.Pool)

	// Register feature modules
	email_log.Module(injector)
	email_template.Module(injector)
	eventbus.Module(injector)
	email_dispatcher.Module(injector, cfg.Resend.APIKey, cfg.Resend.FromAddress)

	// Start event bus
	bus := di.MustInvoke[eventbusPort.EventBus](injector)
	bus.Start(ctx)

	// Build server config
	config := server.Config{
		Port:        cfg.Service.Port(),
		ServiceName: cfg.Service.Name,
	}

	// Create server with route setup
	srv := server.NewServer(config, log, router.SetupRoutes(injector))

	// Register health checkers
	srv.RegisterHealthChecker("database", db)

	if err := srv.Start(ctx); err != nil {
		log.Error("server error", "error", err)
	}

	if err := di.Shutdown(injector); err != nil {
		log.Error("shutdown error", "error", err)
	}
}
