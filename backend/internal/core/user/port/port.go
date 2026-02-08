package port

import (
	basedomain "backend/port"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateUser, UpdateUser]
	basedomain.RepositoryQuery[User]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateUser, UpdateUser]
	basedomain.UseCaseQuery[User]
	basedomain.UseCaseQueryRelation[UserRelation]
	basedomain.UseCaseTx[Service]
}
