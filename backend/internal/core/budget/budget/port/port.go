package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateBudget, UpdateBudget]
	basedomain.RepositoryQuery[Budget]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateBudget, UpdateBudget]
	basedomain.UseCaseQuery[Budget]
	basedomain.UseCaseTx[Service]
}
