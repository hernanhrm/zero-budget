package domain

import (
	"context"

	"backend/infra/dafi"
)

type UseCaseTx[T any] interface {
	WithTx(tx Transaction) T
}

type UseCaseCommand[C, U any] interface {
	UseCaseCreate[C]
	UseCaseUpdate[U]
	UseCaseDelete
}

type UseCaseCreate[T any] interface {
	Create(ctx context.Context, entity T) error
	CreateBulk(ctx context.Context, entities List[T]) error
}

type UseCaseUpdate[T any] interface {
	Update(ctx context.Context, entity T, filters ...dafi.Filter) error
}

type UseCaseDelete interface {
	Delete(ctx context.Context, filters ...dafi.Filter) error
}

type UseCaseQuery[M any] interface {
	UseCaseFindOne[M]
	UseCaseFindAll[List[M]]
}

type UseCaseQueryRelation[M any] interface {
	UseCaseFindOneRelation[M]
	UseCaseFindAllRelation[[]M]
}

type UseCaseFindOne[T any] interface {
	FindOne(ctx context.Context, criteria dafi.Criteria) (T, error)
}

type UseCaseFindAll[T any] interface {
	FindAll(ctx context.Context, criteria dafi.Criteria) (T, error)
}

type UseCaseFindOneRelation[T any] interface {
	FindOneRelation(ctx context.Context, criteria dafi.Criteria) (T, error)
}

type UseCaseFindAllRelation[T any] interface {
	FindAllRelation(ctx context.Context, criteria dafi.Criteria) (T, error)
}
