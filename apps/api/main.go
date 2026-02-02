package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"api/router"
	"backend/app/api_route"
	"backend/app/auth"
	"backend/app/email_log"
	"backend/app/email_template"
	"backend/app/organization"
	"backend/app/permission"
	"backend/app/role"
	"backend/app/user"
	"backend/app/workspace"
	"backend/app/workspace_member"
	"backend/domain"
	"backend/infra/database"
	"backend/infra/di"
	"backend/infra/localconfig"
	"backend/infra/logger"
	"backend/infra/server"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dbURL, migrationsPath string, log domain.Logger) error {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

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

	if !cfg.Service.SkipMigrations {
		if err := runMigrations(cfg.Database.URL, cfg.Service.MigrationsPath, log); err != nil {
			log.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
		log.Info("migrations completed successfully")
	} else {
		log.Info("migrations skipped")
	}

	db, err := database.NewConnection(ctx, cfg.Database.URL, log)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	di.ProvideValue(injector, db)
	di.ProvideValue(injector, db.Pool)

	// Register feature modules
	user.Module(injector)
	email_log.Module(injector)
	email_template.Module(injector)
	organization.Module(injector)
	workspace.Module(injector)
	permission.Module(injector)
	role.Module(injector)
	api_route.Module(injector)
	workspace_member.Module(injector)
	auth.Module(injector, cfg.JWTSecret)

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
