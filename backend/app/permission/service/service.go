package service

import (
	"context"

	"backend/app/permission/domain"
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
		logger: logger.With("component", "permission.service"),
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Permission, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.Permission{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return item, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Permission], error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return items, nil
}

func (s service) Create(ctx context.Context, input domain.CreatePermission) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("permission created")
	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreatePermission]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("permissions created", "count", len(inputs))
	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdatePermission, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("permission updated")
	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("permission deleted")
	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (domain.PermissionRelation, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.PermissionRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return domain.PermissionRelation{Permission: item}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]domain.PermissionRelation, error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	result := make([]domain.PermissionRelation, len(items))
	for i, item := range items {
		result[i] = domain.PermissionRelation{Permission: item}
	}
	return result, nil
}
