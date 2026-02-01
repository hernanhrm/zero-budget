// Package localconfig provides configuration loading from environment variables and .env files.
package localconfig

import (
	"os"
	"strconv"
)

// ConfigOptions defines options for loading configuration.
type ConfigOptions struct {
	EnvPath     string // Relative path to directory containing the env file (default: ".")
	EnvFileName string // Name of the env file (default: ".env")
}

// LocalConfig holds the complete application configuration.
type LocalConfig struct {
	Service   Service
	Database  Database
	JWTSecret string
}

// Service holds service-specific configuration.
type Service struct {
	Port           func() int
	Name           string
	DocsPath       string
	MigrationsPath string
	SkipMigrations bool
}

// Database holds database connection configuration.
type Database struct {
	URL string
}

func getPort() int {
	p := os.Getenv("PORT")
	if p == "" {
		return 8080
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 8080
	}
	return port
}
