---
name: scaffold-backend-module
description: |
  Generate a new backend module in `backend/internal/core/<name>`, register it in `cmd/api`, and create the necessary database migration with RLS.
  Use this skill when the user asks to "create a new backend module" or "scaffold a module".
user-invocable: true
allowed-tools: [Write, Edit, Bash, Glob, Read]
---

# Scaffold Backend Module

This skill orchestrates the creation of a new backend feature module. It enforces Hexagonal Architecture, wires the API, and secures the database with RLS.

## Important: Client-Provided IDs

All modules in this project use **client-provided UUIDs** for entity IDs. The client is responsible for generating and sending the ID when creating entities. The repository uses `input.ID` directly instead of generating IDs.

## Usage

**Trigger:** "Create a new backend module named [name] in schema [schema]"

**Replacements:**
- `{{module}}`: Module name in lowercase (e.g., "order").
- `{{Module}}`: Module name in TitleCase (e.g., "Order").
- `{{module}}s`: Plural module name (e.g., "orders").
- `{{schema}}`: Database schema name (e.g., "public", "store").

## Workflow

### 1. Context Analysis (CRITICAL FIRST STEP)

Before generating any code or migrations, you **MUST** check the current state of the project to avoid overwriting existing work or creating duplicates.

1.  **Check for existing module code:**
    ```bash
    ls -d backend/internal/core/{{module}}
    ```
2.  **Check for existing migrations:**
    ```bash
    ls cmd/api/migrations/*{{module}}*
    ```
3.  **Check for existing route registration:**
    ```bash
    grep -i "{{module}}" cmd/api/router/router.go
    ```

**Decision Logic:**
- **If migrations exist:** List them to the user. **DO NOT** generate a "Create Table" migration unless explicitly asked. Ask the user if they want to generate a *new* migration for permissions/policies only, or if the existing migration is sufficient.
- **If code exists:** Warn the user and ask for confirmation before overwriting any files.

---

### 2. Create Backend Module Structure
**Location:** `backend/internal/core/{{module}}/`

Generate the following directory structure and files:

```
backend/internal/core/{{module}}/
├── port/
│   ├── port.go
│   ├── command.go
│   └── query.go
├── adapter/
│   ├── handler/
│   │   └── http.go
│   └── postgres/
│       └── postgres.go
├── core/
│   └── core.go
├── module.go
└── go.mod
```

**File: `go.mod`**
```go
module backend/core/{{module}}

go 1.24.0

toolchain go1.24.12

require (
	github.com/google/uuid v1.6.0
	github.com/guregu/null/v6 v6.0.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/labstack/echo/v4 v4.15.0
	github.com/samber/do/v2 v2.0.0
	github.com/samber/oops v1.21.0
)
```

**File: `module.go`**
```go
package {{module}}

import (
	baseport "backend/port"
	"backend/adapter/database"
	"backend/adapter/di"
	"backend/core/{{module}}/port"
	"backend/core/{{module}}/adapter/handler"
	"backend/core/{{module}}/adapter/postgres"
	"backend/core/{{module}}/core"
	"github.com/samber/do/v2"
)

func Module(i do.Injector) {
	di.Provide(i, func(i do.Injector) (port.Repository, error) {
		db := di.MustInvoke[database.PoolInterface](i)
		logger := di.MustInvoke[baseport.Logger](i)
		return postgres.NewPostgres(db, logger), nil
	})

	di.Provide(i, func(i do.Injector) (port.Service, error) {
		repo := di.MustInvoke[port.Repository](i)
		logger := di.MustInvoke[baseport.Logger](i)
		return core.New(repo, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		logger := di.MustInvoke[baseport.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
```

**File: `port/command.go`**
```go
package port

import (
	"context"
	"backend/adapter/validation"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Create{{Module}} struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (c Create{{Module}}) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &c,
		validation.Field(&c.ID, validation.Required, validation.IsUUID),
		validation.Field(&c.Name, validation.Required, validation.Length(2, 100)),
	)
}

type Update{{Module}} struct {
	Name null.String `json:"name"`
}

func (u Update{{Module}}) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 100)),
	)
}
```

**File: `port/query.go`**
```go
package port

import (
	"time"
	"github.com/google/uuid"
)

type {{Module}} struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type {{Module}}Relation struct {
	{{Module}}
}
```

**File: `port/port.go`**
```go
package port

import baseport "backend/port"

type Repository interface {
	baseport.RepositoryCommand[Create{{Module}}, Update{{Module}}]
	baseport.RepositoryQuery[{{Module}}]
	baseport.RepositoryTx[Repository]
}

type Service interface {
	baseport.UseCaseCommand[Create{{Module}}, Update{{Module}}]
	baseport.UseCaseQuery[{{Module}}]
	baseport.UseCaseQueryRelation[{{Module}}Relation]
	baseport.UseCaseTx[Service]
}
```

**File: `core/core.go`**
```go
package core

import (
	"context"
	"backend/core/{{module}}/port"
	baseport "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"github.com/samber/oops"
)

type service struct {
	repo   port.Repository
	logger baseport.Logger
}

func New(repo port.Repository, logger baseport.Logger) port.Service {
	return service{
		repo:   repo,
		logger: logger.With("component", "{{module}}.core"),
	}
}

func (s service) WithTx(tx baseport.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.{{Module}}, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.{{Module}}{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return item, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (baseport.List[port.{{Module}}], error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return items, nil
}

func (s service) Create(ctx context.Context, input port.Create{{Module}}) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("{{module}} created")
	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs baseport.List[port.Create{{Module}}]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("{{module}}s created", "count", len(inputs))
	return nil
}

func (s service) Update(ctx context.Context, input port.Update{{Module}}, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("{{module}} updated")
	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("{{module}} deleted")
	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (port.{{Module}}Relation, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.{{Module}}Relation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return port.{{Module}}Relation{ {{Module}}: item }, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]port.{{Module}}Relation, error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	result := make([]port.{{Module}}Relation, len(items))
	for i, item := range items {
		result[i] = port.{{Module}}Relation{ {{Module}}: item }
	}
	return result, nil
}
```

**File: `adapter/postgres/postgres.go`**
```go
package postgres

import (
	"context"
	"errors"
	"time"

	"backend/core/{{module}}/port"
	baseport "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/adapter/database"
	"backend/infra/sqlcraft"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
)

type dbConn interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

const tableName = "{{schema}}.{{module}}s"

var columns = []string{
	"id",
	"name",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":        "id",
	"name":      "name",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

type postgres struct {
	db     dbConn
	logger baseport.Logger
}

func NewPostgres(db database.PoolInterface, logger baseport.Logger) port.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "{{module}}.postgres"),
	}
}

func (r postgres) WithTx(tx baseport.Transaction) port.Repository {
	return postgres{
		db:     tx.GetTx(),
		logger: r.logger,
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.{{Module}}, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.{{Module}}{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var item port.{{Module}}
	err = row.Scan(&item.ID, &item.Name, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.{{Module}}{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.{{Module}}{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return item, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (baseport.List[port.{{Module}}], error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		OrderBy(criteria.Sorts...).
		Limit(criteria.Pagination.PageSize).
		Page(criteria.Pagination.PageNumber).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	rows, err := r.db.Query(ctx, result.SQL, result.Args...)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	defer rows.Close()

	var items baseport.List[port.{{Module}}]
	for rows.Next() {
		var item port.{{Module}}
		err = rows.Scan(&item.ID, &item.Name, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r postgres) Create(ctx context.Context, input port.Create{{Module}}) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(input.ID, input.Name, now, now)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) CreateBulk(ctx context.Context, inputs baseport.List[port.Create{{Module}}]) error {
	if inputs.IsEmpty() { return nil }
	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		query = query.WithValues(input.ID, input.Name, now, now)
	}

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("bulk insert", "sql", result.SQL, "count", len(inputs))
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) Update(ctx context.Context, input port.Update{{Module}}, filters ...dafi.Filter) error {
	cols := []string{}
	vals := []any{}
	if input.Name.Valid {
		cols = append(cols, "name")
		vals = append(vals, input.Name.String)
	}
	cols = append(cols, "updated_at")
	vals = append(vals, time.Now())

	query := sqlcraft.Update(tableName).
		WithColumns(cols...).
		WithValues(vals...).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}

func (r postgres) Delete(ctx context.Context, filters ...dafi.Filter) error {
	query := sqlcraft.DeleteFrom(tableName).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)
	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	return nil
}
```

**File: `adapter/handler/http.go`**
```go
package handler

import (
	"backend/core/{{module}}/port"
	baseport "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/infra/httpresponse"
	"github.com/labstack/echo/v4"
	"github.com/samber/oops"
)

type HTTP struct {
	svc    port.Service
	logger baseport.Logger
}

func NewHTTP(svc port.Service, logger baseport.Logger) HTTP {
	return HTTP{
		svc:    svc,
		logger: logger.With("component", "{{module}}.handler"),
	}
}

func (h HTTP) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	criteria := dafi.Where("id", dafi.Equal, id)
	item, err := h.svc.FindOne(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.OK(c, item)
}

func (h HTTP) FindAll(c echo.Context) error {
	ctx := c.Request().Context()
	parser := dafi.NewQueryParser()
	criteria, err := parser.Parse(c.QueryParams())
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	items, err := h.svc.FindAll(ctx, criteria)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.OK(c, items)
}

func (h HTTP) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var input port.Create{{Module}}
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	if err := h.svc.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.Created(c, nil)
}

func (h HTTP) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	var input port.Update{{Module}}
	if err := c.Bind(&input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Code(apperrors.CodeBadRequest).Wrap(err)
	}
	filters := dafi.FilterBy("id", dafi.Equal, id)
	if err := h.svc.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.NoContent(c)
}

func (h HTTP) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	filters := dafi.FilterBy("id", dafi.Equal, id)
	if err := h.svc.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerHandler).Wrap(err)
	}
	return httpresponse.NoContent(c)
}
```

---

### 3. API Registration

**File: `backend/cmd/api/router/{{module}}.go`**
Create a new file to map the routes.

```go
package router

import (
	"backend/core/{{module}}/adapter/handler"
	"backend/adapter/di"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func Register{{Module}}Routes(injector do.Injector, e *echo.Echo) {
	h := di.MustInvoke[handler.HTTP](injector)

	g := e.Group("/{{module}}s")

	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.FindAll)
	g.GET("/:id", h.FindOne)
}
```

**Task: Update `backend/cmd/api/router/router.go`**
Add `Register{{Module}}Routes(injector, e)` to the `SetupRoutes` function.

**Task: Update `backend/cmd/api/main.go`**
Import the module (`backend/core/{{module}}`) and call `{{module}}.Module(injector)` in `main()`.

---

### 4. Database Migration (Conditional)

**CONDITION:** **Only generate this file if NO existing migration for this module was found in Step 1.**

**Task: Create Migration File**
Use `pg-aiguide` or manual file creation in `backend/cmd/api/migrations/YYYYMMDDHHMMSS_create_{{module}}s_table.up.sql`.

**Migration Template:**
```sql
CREATE TABLE {{schema}}.{{module}}s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Note: While DEFAULT gen_random_uuid() is kept for backward compatibility,
-- the client should always provide the ID when creating entities.

-- Indexes
CREATE INDEX idx_{{module}}s_name ON {{schema}}.{{module}}s(name);

-- RLS Enablement
ALTER TABLE {{schema}}.{{module}}s ENABLE ROW LEVEL SECURITY;

-- Seed Permissions
INSERT INTO auth.permissions (slug, description) VALUES
    ('{{module}}.read', 'View {{module}}s'),
    ('{{module}}.create', 'Create {{module}}s'),
    ('{{module}}.update', 'Update {{module}}s'),
    ('{{module}}.delete', 'Delete {{module}}s');

-- Seed API Routes Mapping
DO $$
DECLARE
    p_read UUID;
    p_create UUID;
    p_update UUID;
    p_delete UUID;
BEGIN
    SELECT id INTO p_read FROM auth.permissions WHERE slug = '{{module}}.read';
    SELECT id INTO p_create FROM auth.permissions WHERE slug = '{{module}}.create';
    SELECT id INTO p_update FROM auth.permissions WHERE slug = '{{module}}.update';
    SELECT id INTO p_delete FROM auth.permissions WHERE slug = '{{module}}.delete';

    INSERT INTO auth.api_routes (method, path, permission_id) VALUES
    ('GET', '/{{module}}s', p_read),
    ('GET', '/{{module}}s/:id', p_read),
    ('POST', '/{{module}}s', p_create),
    ('PUT', '/{{module}}s/:id', p_update),
    ('DELETE', '/{{module}}s/:id', p_delete);
END $$;

-- Auto-grant to Admin Role
DO $$
DECLARE
    r_admin UUID;
    p_read UUID;
    p_create UUID;
    p_update UUID;
    p_delete UUID;
BEGIN
    SELECT id INTO r_admin FROM auth.roles WHERE name = 'Admin' AND workspace_id IS NULL;
    
    SELECT id INTO p_read FROM auth.permissions WHERE slug = '{{module}}.read';
    SELECT id INTO p_create FROM auth.permissions WHERE slug = '{{module}}.create';
    SELECT id INTO p_update FROM auth.permissions WHERE slug = '{{module}}.update';
    SELECT id INTO p_delete FROM auth.permissions WHERE slug = '{{module}}.delete';

    INSERT INTO auth.role_permissions (role_id, permission_id) VALUES 
    (r_admin, p_read), (r_admin, p_create), (r_admin, p_update), (r_admin, p_delete);
END $$;

-- RLS Policies
-- Read: Requires '{{module}}.read' permission
CREATE POLICY {{module}}s_read ON {{schema}}.{{module}}s FOR SELECT USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, NULL, '{{module}}.read')
);

-- Write: Requires '{{module}}.create' permission
CREATE POLICY {{module}}s_insert ON {{schema}}.{{module}}s FOR INSERT WITH CHECK (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, NULL, '{{module}}.create')
);

-- Update: Requires '{{module}}.update' permission
CREATE POLICY {{module}}s_update ON {{schema}}.{{module}}s FOR UPDATE USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, NULL, '{{module}}.update')
);

-- Delete: Requires '{{module}}.delete' permission
CREATE POLICY {{module}}s_delete ON {{schema}}.{{module}}s FOR DELETE USING (
    auth.check_permission(current_setting('app.current_user_id', true)::UUID, NULL, '{{module}}.delete')
);
```

---

### 5. Final Steps
1. Run `go mod tidy` in `backend/internal/core/{{module}}`.
2. Run `go mod tidy` in `backend/cmd/api`.
3. Notify the user to run migrations.
