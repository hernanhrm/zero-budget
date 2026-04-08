package account

import (
	"backend/core/budget/account/adapter/handler"
	"backend/core/budget/account/adapter/postgres"
	"backend/core/budget/account/core"
	"backend/core/budget/account/port"
	transactionport "backend/core/budget/transaction/port"
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
		txnRepo := di.MustInvoke[transactionport.Repository](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(repo, txnRepo, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
