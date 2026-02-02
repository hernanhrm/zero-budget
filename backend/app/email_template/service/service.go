package service

import (
	"context"

	"backend/app/email_template/domain"
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
		logger: logger.With("component", "email_template.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) domain.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (domain.EmailTemplate, error) {
	tmpl, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return domain.EmailTemplate{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return tmpl, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[domain.EmailTemplate], error) {
	tmpls, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return tmpls, nil
}

func (s service) Create(ctx context.Context, input domain.CreateEmailTemplate) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email template created", "event", input.Event)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[domain.CreateEmailTemplate]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email templates created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input domain.UpdateEmailTemplate, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email template updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email template deleted")

	return nil
}
