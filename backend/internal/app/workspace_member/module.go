package workspace_member

import (
	"backend/app/workspace_member/domain"
	"backend/app/workspace_member/handler"
	"backend/app/workspace_member/repository"
	"backend/app/workspace_member/service"
	basedomain "backend/domain"
	"backend/infra/database"
	"backend/infra/di"
	"github.com/samber/do/v2"
)

// Module registers all workspace_member feature services into the DI container
func Module(i do.Injector) {
	di.Provide(i, func(i do.Injector) (domain.Repository, error) {
		db := di.MustInvoke[database.PoolInterface](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return repository.NewPostgres(db, logger), nil
	})

	di.Provide(i, func(i do.Injector) (domain.Service, error) {
		repo := di.MustInvoke[domain.Repository](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return service.New(repo, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[domain.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
