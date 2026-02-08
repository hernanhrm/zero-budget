package workspace_member

import (
	"backend/core/workspace_member/adapter/handler"
	"backend/core/workspace_member/adapter/postgres"
	"backend/core/workspace_member/core"
	"backend/core/workspace_member/port"
	basedomain "backend/port"
	"backend/adapter/database"
	"backend/adapter/di"
	"github.com/samber/do/v2"
)

// Module registers all workspace_member feature services into the DI container
func Module(i do.Injector) {
	di.Provide(i, func(i do.Injector) (port.Repository, error) {
		db := di.MustInvoke[database.PoolInterface](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return postgres.NewPostgres(db, logger), nil
	})

	di.Provide(i, func(i do.Injector) (port.Service, error) {
		repo := di.MustInvoke[port.Repository](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(repo, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
