package service

import (
	"context"

	"backend/app/role/domain"
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
		logger: logger.With("component", "role.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) domain.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Role, error) {
	role, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.Role{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return role, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Role], error) {
	roles, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return roles, nil
}

func (s service) Create(ctx context.Context, input domain.CreateRole) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("role created", "name", input.Name)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateRole]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("roles created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdateRole, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("role updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("role deleted")

	return nil
}
