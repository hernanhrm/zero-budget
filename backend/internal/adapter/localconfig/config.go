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
	Resend    Resend
	Identity  Identity
	JWTSecret string
}

// Identity holds identity service configuration.
type Identity struct {
	URL            string
	InternalAPIKey string
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

// Resend holds Resend email service configuration.
type Resend struct {
	APIKey      string
	FromAddress string
}

func getPort() int {
	p := os.Getenv("SERVICE_PORT")
	if p == "" {
		return 8080
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 8080
	}
	return port
}
