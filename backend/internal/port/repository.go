package domain

import (
	"context"

	"backend/infra/dafi"
)

type RepositoryTx[T any] interface {
	WithTx(tx Transaction) T
}

type RepositoryCommand[C, U any] interface {
	RepositoryCreate[C]
	RepositoryUpdate[U]
	RepositoryDelete
}

type RepositoryCreate[T any] interface {
	Create(ctx context.Context, entity T) error
	CreateBulk(ctx context.Context, entities List[T]) error
}

type RepositoryUpdate[T any] interface {
	Update(ctx context.Context, entity T, filters ...dafi.Filter) error
}

type RepositoryDelete interface {
	Delete(ctx context.Context, filters ...dafi.Filter) error
}

type RepositoryQuery[M any] interface {
	FindOne(ctx context.Context, criteria dafi.Criteria) (M, error)
	FindAll(ctx context.Context, criteria dafi.Criteria) (List[M], error)
}

type RepositoryQueryRelation[M any] interface {
	FindOneRelation(ctx context.Context, criteria dafi.Criteria) (M, error)
	FindAllRelation(ctx context.Context, criteria dafi.Criteria) ([]M, error)
}
