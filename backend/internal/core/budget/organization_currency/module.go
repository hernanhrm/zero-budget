package organization_currency

import (
	currencypkg "backend/core/budget/currency/port"
	"backend/core/budget/organization_currency/adapter/handler"
	"backend/core/budget/organization_currency/adapter/postgres"
	"backend/core/budget/organization_currency/core"
	"backend/core/budget/organization_currency/port"
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
		currencyRepo := di.MustInvoke[currencypkg.Repository](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(repo, currencyRepo, logger), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return handler.NewHTTP(svc, logger), nil
	})
}
