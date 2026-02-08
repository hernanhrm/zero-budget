package domain

import "context"

type Transaction interface {
	GetTx() Tx
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
	Commit(ctx context.Context, tx Transaction) error
	Rollback(ctx context.Context, tx Transaction) error
}
