package port

import (
	basedomain "backend/port"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateOrganization, UpdateOrganization]
	basedomain.RepositoryQuery[Organization]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateOrganization, UpdateOrganization]
	basedomain.UseCaseQuery[Organization]
	basedomain.UseCaseQueryRelation[OrganizationRelation]
	basedomain.UseCaseTx[Service]
}
