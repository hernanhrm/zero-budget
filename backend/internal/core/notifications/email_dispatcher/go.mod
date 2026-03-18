module backend/core/notifications/email_dispatcher

go 1.24.12

require (
	backend/adapter/di v0.0.0
	backend/core/notifications/email_log v0.0.0
	backend/core/notifications/email_template v0.0.0
	backend/core/notifications/eventbus v0.0.0
	backend/core/notifications/events v0.0.0
	backend/infra/dafi v0.0.0
	backend/port v0.0.0
	backend/port/errors v0.0.0
	github.com/google/uuid v1.6.0
	github.com/guregu/null/v6 v6.0.0
	github.com/samber/do/v2 v2.0.0
	github.com/samber/oops v1.21.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	github.com/samber/lo v1.52.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

replace (
	backend/adapter/di => ../../../adapter/di
	backend/core/notifications/email_log => ../email_log
	backend/core/notifications/email_template => ../email_template
	backend/core/notifications/eventbus => ../eventbus
	backend/core/notifications/events => ../events
	backend/infra/dafi => ../../../../pkg/dafi
	backend/port => ../../../port
	backend/port/errors => ../../../port/errors
)
