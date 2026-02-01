package domain

import (
	basedomain "backend/domain"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateRole, UpdateRole]
	basedomain.RepositoryQuery[Role]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateRole, UpdateRole]
	basedomain.UseCaseQuery[Role]
	basedomain.UseCaseTx[Service]
}
