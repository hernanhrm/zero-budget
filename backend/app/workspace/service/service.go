package service

import (
	"context"

	"backend/app/workspace/domain"
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
		logger: logger.With("component", "workspace.service"),
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.Workspace, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.Workspace{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return item, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.Workspace], error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return items, nil
}

func (s service) Create(ctx context.Context, input domain.CreateWorkspace) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("workspace created")
	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateWorkspace]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("workspaces created", "count", len(inputs))
	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdateWorkspace, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}
	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("workspace updated")
	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	s.logger.WithContext(ctx).Info("workspace deleted")
	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (domain.WorkspaceRelation, error) {
	item, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.WorkspaceRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	return domain.WorkspaceRelation{Workspace: item}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]domain.WorkspaceRelation, error) {
	items, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}
	result := make([]domain.WorkspaceRelation, len(items))
	for i, item := range items {
		result[i] = domain.WorkspaceRelation{Workspace: item}
	}
	return result, nil
}
