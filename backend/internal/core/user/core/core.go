package core

import (
	"context"

	"backend/core/user/port"
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
		logger: logger.With("component", "user.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.User, error) {
	user, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.User{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return user, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.User], error) {
	users, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return users, nil
}

func (s service) Create(ctx context.Context, input port.CreateUser) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user created", "email", input.Email)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateUser]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("users created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateUser, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("user deleted")

	return nil
}

func (s service) FindOneRelation(ctx context.Context, criteria dafi.Criteria) (port.UserRelation, error) {
	user, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.UserRelation{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return port.UserRelation{User: user}, nil
}

func (s service) FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]port.UserRelation, error) {
	users, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	result := make([]port.UserRelation, len(users))
	for i, user := range users {
		result[i] = port.UserRelation{User: user}
	}

	return result, nil
}
