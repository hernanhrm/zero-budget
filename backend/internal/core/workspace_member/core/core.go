package core

import (
	"context"

	"backend/core/workspace_member/port"
	basedomain "backend/port"
	"backend/infra/dafi"
)

type service struct {
	repo   port.Repository
	logger basedomain.Logger
}

func New(repo port.Repository, logger basedomain.Logger) port.Service {
	return service{
		repo:   repo,
		logger: logger.With("component", "workspace_member.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:   s.repo.WithTx(tx),
		logger: s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.WorkspaceMember, error) {
	return s.repo.FindOne(ctx, criteria)
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.WorkspaceMember], error) {
	return s.repo.FindAll(ctx, criteria)
}

func (s service) Create(ctx context.Context, input port.CreateWorkspaceMember) error {
	if err := input.Validate(ctx); err != nil {
		return err
	}
	return s.repo.Create(ctx, input)
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateWorkspaceMember]) error {
	for _, input := range inputs {
		if err := input.Validate(ctx); err != nil {
			return err
		}
	}
	return s.repo.CreateBulk(ctx, inputs)
}

func (s service) Update(ctx context.Context, input port.UpdateWorkspaceMember, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return err
	}
	return s.repo.Update(ctx, input, filters...)
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	return s.repo.Delete(ctx, filters...)
}
