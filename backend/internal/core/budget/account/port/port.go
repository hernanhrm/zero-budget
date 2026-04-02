package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateAccount, UpdateAccount]
	basedomain.RepositoryQuery[Account]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateAccount, UpdateAccount]
	basedomain.UseCaseQuery[Account]
	basedomain.UseCaseTx[Service]
}
