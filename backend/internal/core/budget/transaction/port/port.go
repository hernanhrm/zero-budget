package port

import (
	"context"

	basedomain "backend/port"

	"github.com/google/uuid"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateTransaction, UpdateTransaction]
	basedomain.RepositoryQuery[Transaction]
	basedomain.RepositoryTx[Repository]
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
	ExistsForOrganization(ctx context.Context, organizationID string) (bool, error)
}

type Service interface {
	basedomain.UseCaseCommand[CreateTransaction, UpdateTransaction]
	basedomain.UseCaseQuery[Transaction]
	basedomain.UseCaseTx[Service]
}
