package core

import (
	"context"

	"backend/core/budget/organization_currency/port"
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
		logger: logger.With("component", "organization_currency.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.OrganizationCurrency, error) {
	oc, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return oc, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.OrganizationCurrency], error) {
	ocs, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return ocs, nil
}

func (s service) Create(ctx context.Context, input port.CreateOrganizationCurrency) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency created", "currencyCode", input.CurrencyCode)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateOrganizationCurrency]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currencies created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateOrganizationCurrency, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency deleted")

	return nil
}
