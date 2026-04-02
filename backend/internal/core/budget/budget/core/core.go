package core

import (
	"context"

	"backend/core/budget/budget/port"
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
		logger: logger.With("component", "budget.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Budget, error) {
	b, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.Budget{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return b, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Budget], error) {
	budgets, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return budgets, nil
}

func (s service) Create(ctx context.Context, input port.CreateBudget) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("budget created", "name", input.Name)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateBudget]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("budgets created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateBudget, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("budget updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("budget deleted")

	return nil
}
