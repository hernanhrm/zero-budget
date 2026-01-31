package domain

import (
	"context"

	basedomain "backend/domain"
	"backend/infra/dafi"
)

type Repository interface {
	FindOne(ctx context.Context, criteria dafi.Criteria) (WorkspaceMember, error)
	FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[WorkspaceMember], error)
	Create(ctx context.Context, input CreateWorkspaceMember) error
	CreateBulk(ctx context.Context, inputs basedomain.List[CreateWorkspaceMember]) error
	Update(ctx context.Context, input UpdateWorkspaceMember, filters ...dafi.Filter) error
	Delete(ctx context.Context, filters ...dafi.Filter) error
}

type Service interface {
	FindOne(ctx context.Context, criteria dafi.Criteria) (WorkspaceMember, error)
	FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[WorkspaceMember], error)
	Create(ctx context.Context, input CreateWorkspaceMember) error
	CreateBulk(ctx context.Context, inputs basedomain.List[CreateWorkspaceMember]) error
	Update(ctx context.Context, input UpdateWorkspaceMember, filters ...dafi.Filter) error
	Delete(ctx context.Context, filters ...dafi.Filter) error
}
