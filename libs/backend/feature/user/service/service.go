package service

import (
	"context"

	basedomain "backend/domain"
	apperrors "backend/domain/errors"
	"backend/infra/dafi"
	"feature/user/domain"

	"github.com/samber/oops"
)

type service struct {
	repo   domain.Repository
	logger basedomain.Logger
}

func New(repo domain.Repository, logger basedomain.Logger) domain.Service {
	return service{
		repo:   repo,
		logger: logger.With("component", "user.service"),
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.User, error) {
	user, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.User{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return user, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.User], error) {
	users, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return users, nil
}

func (s service) Create(ctx context.Context, input domain.CreateUser) error {
	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.Info("user created", "email", input.Email)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateUser]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.Info("users created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdateUser, filters ...dafi.Filter) error {
	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.Info("user updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.Info("user deleted")

	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (domain.UserRelation, error) {
	user, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.UserRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return domain.UserRelation{User: user}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]domain.UserRelation, error) {
	users, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := make([]domain.UserRelation, len(users))
	for i, user := range users {
		result[i] = domain.UserRelation{User: user}
	}

	return result, nil
}
