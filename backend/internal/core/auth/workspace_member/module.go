package workspace_member

import (
	roleport "backend/core/auth/role/port"
	userport "backend/core/auth/user/port"
	"backend/core/auth/workspace_member/adapter/handler"
	"backend/core/auth/workspace_member/adapter/postgres"
	"backend/core/auth/workspace_member/core"
	"backend/core/auth/workspace_member/port"
	"backend/adapter/database"
	"backend/adapter/di"
	basedomain "backend/port"

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
		userSvc := di.MustInvoke[userport.Service](i)
		roleSvc := di.MustInvoke[roleport.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(repo, userSvc, roleSvc, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
