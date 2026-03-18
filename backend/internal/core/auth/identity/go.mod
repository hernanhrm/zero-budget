module backend/core/auth/identity

go 1.24.12

require (
	backend/adapter/di v0.0.0
	backend/core/auth/organization v0.0.0
	backend/core/auth/role v0.0.0
	backend/core/auth/user v0.0.0
	backend/core/auth/workspace v0.0.0
	backend/core/auth/workspace_member v0.0.0
	backend/core/notifications/eventbus v0.0.0
	backend/core/notifications/events v0.0.0
	backend/infra/dafi v0.0.0
	backend/infra/httpresponse v0.0.0
	backend/port v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/labstack/echo/v4 v4.15.0
	github.com/samber/do/v2 v2.0.0
	github.com/samber/oops v1.21.0
	golang.org/x/crypto v0.46.0
)

require (
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/guregu/null/v6 v6.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	github.com/samber/lo v1.52.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

replace (
	backend/adapter/di => ../../../adapter/di
	backend/core/auth/organization => ../organization
	backend/core/auth/role => ../role
	backend/core/auth/user => ../user
	backend/core/auth/workspace => ../workspace
	backend/core/auth/workspace_member => ../workspace_member
	backend/core/notifications/eventbus => ../../notifications/eventbus
	backend/core/notifications/events => ../../notifications/events
	backend/infra/dafi => ../../../../pkg/dafi
	backend/infra/httpresponse => ../../../../pkg/httpresponse
	backend/port => ../../../port
)
