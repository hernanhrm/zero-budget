package domain

import basedomain "backend/domain"

type Repository interface {
	basedomain.RepositoryCommand[CreatePermission, UpdatePermission]
	basedomain.RepositoryQuery[Permission]
}

type Service interface {
	basedomain.UseCaseCommand[CreatePermission, UpdatePermission]
	basedomain.UseCaseQuery[Permission]
	basedomain.UseCaseQueryRelation[PermissionRelation]
}
