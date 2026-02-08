package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreatePermission, UpdatePermission]
	basedomain.RepositoryQuery[Permission]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreatePermission, UpdatePermission]
	basedomain.UseCaseQuery[Permission]
	basedomain.UseCaseQueryRelation[PermissionRelation]
	basedomain.UseCaseTx[Service]
}
