package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create source: %v\n", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		fmt.Println("Running migrations up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "migration up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid steps: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Rolling back %d migration(s)...\n", steps)
		if err := m.Steps(-steps); err != nil {
			fmt.Fprintf(os.Stderr, "migration down failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Rollback completed successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\nUsage: migrate [up|down [steps]|version]\n", cmd)
		os.Exit(1)
	}
}
