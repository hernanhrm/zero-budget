package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateTransaction, UpdateTransaction]
	basedomain.RepositoryQuery[Transaction]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateTransaction, UpdateTransaction]
	basedomain.UseCaseQuery[Transaction]
	basedomain.UseCaseTx[Service]
}
