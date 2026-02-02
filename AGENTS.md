# AGENTS.md - Zero Budget Development Guidelines

This document provides guidelines for agentic coding agents working in this repository.

## Project Overview

This is a Go-based monorepo using Nx for task orchestration. The project follows a layered architecture:
- **apps/api**: Main API application entry point
- **backend/app/***: Business logic modules (api_route, auth, organization, permission, role, user, workspace, workspace_member)
- **backend/domain**: Domain layer with interfaces and shared types
- **backend/infra**: Infrastructure layer (database, DI, logging, HTTP, validation)

## Build, Lint, and Test Commands

### Running Tasks via Nx

All tasks should be run through Nx for consistency and caching benefits:

```bash
# Run all tasks (lint, test, build)
bun nx run-many -t lint test build

# Run a specific target for all projects
bun nx run-many -t test

# Run a target for a specific project
bun nx run api:lint
bun nx run api:test

# Run with affected projects only
bun nx affected -t test build
```

### Go-Specific Commands

```bash
# Run tests for a specific package
cd apps/api && go test ./...
go test ./backend/infra/...

# Run a single test file
go test -run TestSlogAdapter_JSONFormat ./backend/infra/logger

# Run tests with verbose output
go test -v ./backend/infra/logger

# Run tests with coverage
go test -cover ./...

# Run gofmt linter
gofmt -l .

# Run go vet
go vet ./...

# Tidy go modules (per package)
cd apps/api && go mod tidy
cd backend/infra/database && go mod tidy
```

### Database Migrations

```bash
# Run migrations (via Nx)
bun nx run api:migrate:up

# Rollback one migration
bun nx run api:migrate:down

# Create new migration
bun nx run api:migrate:create --args.name=add_users_table
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

**Error codes** in `backend/domain/errors/codes.go`:
- `not_found` (404), `bad_request` (400), `validation` (422), `conflict` (409)
- `unauthorized` (401), `forbidden` (403), `already_exists` (409)

**Layer constants** for `oops.In()`: `handler`, `service`, `repository`, `middleware`, `infrastructure`

### Architecture Layers

1. **Handler Layer** (`app/*/handler`): HTTP request/response handling, validation
2. **Service Layer** (`app/*/service`): Business logic, use cases
3. **Repository Layer** (`app/*/repository`): Data access, persistence
4. **Domain Layer** (`backend/domain`): Interfaces, domain types, shared errors
5. **Infrastructure Layer** (`backend/infra`): Database, DI, logging, HTTP server

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

- Use `localconfig` package for configuration management (see `backend/infra/localconfig`)
- Use `domain.Logger` interface from `backend/domain`
- Create loggers via `backend/infra/logger` package
- Use `NewProduction()` for production, `NewDevelopment()` for local dev
- Use `backend/infra/httpresponse` for consistent API responses

<!-- nx configuration start-->
<!-- Leave the start & end comments to automatically receive updates. -->

# General Guidelines for working with Nx

- When running tasks (for example build, lint, test, e2e, etc.), always prefer running the task through `nx` (i.e. `nx run`, `nx run-many`, `nx affected`) instead of using the underlying tooling directly
- You have access to the Nx MCP server and its tools, use them to help the user
- When answering questions about the repository, use the `nx_workspace` tool first to gain an understanding of the workspace architecture where applicable.
- When working in individual projects, use the `nx_project_details` mcp tool to analyze and understand the specific project structure and dependencies
- For questions around nx configuration, best practices or if you're unsure, use the `nx_docs` tool to get relevant, up-to-date docs. Always use this instead of assuming things about nx configuration
- If the user needs help with an Nx configuration or project graph error, use the `nx_workspace` tool to get any errors
- For Nx plugin best practices, check `node_modules/@nx/<plugin>/PLUGIN.md`. Not all plugins have this file - proceed without it if unavailable.

<!-- nx configuration end-->
