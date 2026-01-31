package domain

import (
	basedomain "backend/domain"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateUser, UpdateUser]
	basedomain.RepositoryQuery[User]
}

type Service interface {
	basedomain.UseCaseCommand[CreateUser, UpdateUser]
	basedomain.UseCaseQuery[User]
	basedomain.UseCaseQueryRelation[UserRelation]
}
