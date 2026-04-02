package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateCategory, UpdateCategory]
	basedomain.RepositoryQuery[Category]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateCategory, UpdateCategory]
	basedomain.UseCaseQuery[Category]
	basedomain.UseCaseTx[Service]
}
