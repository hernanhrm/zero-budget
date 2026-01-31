package workspace

import (
	"backend/app/workspace/domain"
	"backend/app/workspace/handler"
	"backend/app/workspace/repository"
	"backend/app/workspace/service"
	basedomain "backend/domain"
	"backend/infra/database"
	"backend/infra/di"
	"github.com/samber/do/v2"
)

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
