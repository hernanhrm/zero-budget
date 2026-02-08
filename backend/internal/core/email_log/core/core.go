package core

import (
	"context"

	"backend/core/email_log/port"
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
		logger: logger.With("component", "email_log.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.EmailLog, error) {
	log, err := s.repo.FindOne(ctx, criteria)
	if err != nil {
		return port.EmailLog{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return log, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.EmailLog], error) {
	logs, err := s.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	return logs, nil
}

func (s service) Create(ctx context.Context, input port.CreateEmailLog) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email log created", "templateId", input.TemplateID, "status", input.Status)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateEmailLog]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email logs created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateEmailLog, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email log updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("email log deleted")

	return nil
}
