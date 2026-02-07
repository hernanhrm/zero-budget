package user

import (
	basedomain "backend/domain"
	"backend/infra/database"
	"backend/infra/di"
	"backend/app/user/domain"
	"backend/app/user/handler"
	"backend/app/user/repository"
	"backend/app/user/service"
	"github.com/samber/do/v2"
)

// Module registers all user feature services into the DI container
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
