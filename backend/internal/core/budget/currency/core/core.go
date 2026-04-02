package core

import (
	"context"

	"backend/core/budget/currency/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"github.com/samber/oops"
)

type service struct {
	repo   port.Repository
	logger basedomain.Logger
}

func New(repo port.Repository, logger basedomain.Logger) port.Service {
	return service{
		repo:   repo,
		logger: logger.With("component", "currency.service"),
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Currency, error) {
	currency, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.Currency{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return currency, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Currency], error) {
	currencies, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return currencies, nil
}
