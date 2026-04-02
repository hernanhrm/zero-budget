package transaction

import (
	"backend/core/budget/transaction/adapter/handler"
	"backend/core/budget/transaction/adapter/postgres"
	"backend/core/budget/transaction/core"
	"backend/core/budget/transaction/port"
	basedomain "backend/port"
	"backend/adapter/database"
	"backend/adapter/di"
	"github.com/samber/do/v2"
)

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
