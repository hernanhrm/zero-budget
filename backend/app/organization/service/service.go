package service

import (
	"context"

	"backend/app/organization/domain"
	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"github.com/samber/oops"
)

type service struct {
	repo   domain.Repository
	logger basedomain.Logger
}

func New(repo domain.Repository, logger basedomain.Logger) domain.Service {
	return service{
		repo:   repo,
		logger: logger.With("component", "organization.service"),
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Organization, error) {
	organization, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.Organization{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return organization, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Organization], error) {
	organizations, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return organizations, nil
}

func (s service) Create(ctx context.Context, input domain.CreateOrganization) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization created", "name", input.Name)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateOrganization]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organizations created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdateOrganization, filters ...dafi.Filter) error {
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

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (domain.OrganizationRelation, error) {
	organization, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.OrganizationRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return domain.OrganizationRelation{Organization: organization}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]domain.OrganizationRelation, error) {
	organizations, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := make([]domain.OrganizationRelation, len(organizations))
	for i, organization := range organizations {
		result[i] = domain.OrganizationRelation{Organization: organization}
	}

	return result, nil
}
