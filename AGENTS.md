# AGENTS.md - Zero Budget Development Guidelines

This document provides guidelines for agentic coding agents working in this repository.

## Project Overview

This is a Go-based monorepo. The backend lives under the `backend/` directory with its own `go.work` workspace. The project follows a layered architecture:
- **backend/cmd/api**: Main API application entry point
- **backend/internal/app/***: Business logic modules (api_route, auth, organization, permission, role, user, workspace, workspace_member)
- **backend/internal/domain**: Domain layer with interfaces and shared types
- **backend/internal/infra**: Infrastructure layer (database, DI, logging, validation)
- **backend/pkg**: Reusable packages (dafi, sqlcraft, httpresponse)
- **frontend/**: Frontend application (TBD)

## Build, Lint, and Test Commands

### Go Commands

All Go commands should be run from the `backend/` directory where `go.work` lives:

```bash
# Build the API
cd backend && go build ./cmd/api/...

# Run tests for all packages
cd backend && go test ./...

# Run tests for a specific package
cd backend && go test ./internal/infra/logger/...

# Run a single test
cd backend && go test -run TestSlogAdapter_JSONFormat ./internal/infra/logger

# Run tests with verbose output
cd backend && go test -v ./internal/infra/logger

# Run tests with coverage
cd backend && go test -cover ./...

# Run gofmt linter
cd backend && gofmt -l .

# Run go vet
cd backend && go vet ./...

# Tidy go modules (per package)
cd backend/cmd/api && go mod tidy
cd backend/internal/infra/database && go mod tidy
```

### Database Migrations

```bash
# Run migrations
cd backend/cmd/api && migrate -path migrations -database "$DATABASE_URL" up

# Rollback one migration
cd backend/cmd/api && migrate -path migrations -database "$DATABASE_URL" down 1

# Create new migration
cd backend/cmd/api && migrate create -ext sql -dir migrations -format 20060102150405 <name>

# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Database (Podman)

```bash
# Start database
cd backend/internal/infra/database && podman compose up -d

# Stop database
cd backend/internal/infra/database && podman compose down

# View logs
cd backend/internal/infra/database && podman compose logs -f
```

## Code Style Guidelines

### General Principles

- Use value receivers by default for method definitions; avoid pointer receivers unless necessary
- Keep functions small and focused on a single responsibility
- Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- Write tests for all exported functions and significant internal logic

### Naming Conventions

- **Variables and Functions**: camelCase (e.g., `newSlogAdapter`, `healthStatusHealthy`)
- **Constants**: camelCase or SCREAMING_SNAKE_CASE for config constants (e.g., `healthStatusHealthy`)
- **Types**: PascalCase (e.g., `SlogAdapter`, `HealthResponse`)
- **Packages**: lowercase, concise, descriptive (e.g., `logger`, `dafi`, `sqlcraft`)
- **URL Parameters**: kebab-case (e.g., `/users/:user-id` not `/users/:userId`)
- **Database Columns**: snake_case

### Import Organization

Imports should be organized in three groups with blank lines between:

```go
import (
    // Standard library
    "context"
    "net/http"
    "time"

    // Third-party packages
    "github.com/labstack/echo/v4"
    "github.com/samber/oops"

    // Internal - path aliases
    "api/router"
    "backend/app/user"
    "backend/domain"
    "backend/infra/database"
)
```

Use path aliases for internal imports (e.g., `backend/`, `api/` as defined in `go.work`).

### Error Handling

This project uses the [`samber/oops`](https://github.com/samber/oops) library for structured error handling:

```go
// Return an error with a code and context
return oops.Code("not_found").Errorf("user %d not found", userID)

// Wrap an existing error with context
return oops.Code("database_failed").Wrapf(err, "failed to fetch user %d", userID)

// Add layer context for debugging
return oops.In("repository").Code("not_found").Errorf("user %d not found", userID)
```

**Error codes** in `backend/internal/domain/errors/codes.go`:
- `not_found` (404), `bad_request` (400), `validation` (422), `conflict` (409)
- `unauthorized` (401), `forbidden` (403), `already_exists` (409)

**Layer constants** for `oops.In()`: `handler`, `service`, `repository`, `middleware`, `infrastructure`

### Architecture Layers

1. **Handler Layer** (`internal/app/*/handler`): HTTP request/response handling, validation
2. **Service Layer** (`internal/app/*/service`): Business logic, use cases
3. **Repository Layer** (`internal/app/*/repository`): Data access, persistence
4. **Domain Layer** (`internal/domain`): Interfaces, domain types, shared errors
5. **Infrastructure Layer** (`internal/infra`): Database, DI, logging, HTTP server

Dependencies flow downward (handler -> service -> repository -> infrastructure).

### Testing

- Use table-driven tests with the `tests` struct pattern
- Use `stretchr/testify/assert` for assertions
- Name test files `*_test.go` in the same package
- Example structure:

```go
func TestSlogAdapter_JSONFormat(t *testing.T) {
    tests := []struct {
        name      string
        level     Level
        logFunc   func(domain.Logger)
        shouldLog bool
    }{
        // test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic...
        })
    }
}
```

### Configuration and Logging

- Use `localconfig` package for configuration management (see `backend/internal/infra/localconfig`)
- Use `domain.Logger` interface from `backend/domain`
- Create loggers via `backend/internal/infra/logger` package
- Use `NewProduction()` for production, `NewDevelopment()` for local dev
- Use `backend/internal/infra/httpresponse` for consistent API responses
