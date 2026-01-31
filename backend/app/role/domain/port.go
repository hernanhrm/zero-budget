package domain

import (
	basedomain "backend/domain"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateRole, UpdateRole]
	basedomain.RepositoryQuery[Role]
}

type Service interface {
	basedomain.UseCaseCommand[CreateRole, UpdateRole]
	basedomain.UseCaseQuery[Role]
}
