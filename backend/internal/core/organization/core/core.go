package core

import (
	"context"

	"backend/core/organization/port"
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
		logger: logger.With("component", "organization.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Organization, error) {
	organization, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.Organization{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return organization, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Organization], error) {
	organizations, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return organizations, nil
}

func (s service) Create(ctx context.Context, input port.CreateOrganization) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization created", "name", input.Name)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateOrganization]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organizations created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateOrganization, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization deleted")

	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (port.OrganizationRelation, error) {
	organization, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.OrganizationRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return port.OrganizationRelation{Organization: organization}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]port.OrganizationRelation, error) {
	organizations, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := make([]port.OrganizationRelation, len(organizations))
	for i, organization := range organizations {
		result[i] = port.OrganizationRelation{Organization: organization}
	}

	return result, nil
}
