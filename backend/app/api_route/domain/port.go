package domain

import basedomain "backend/domain"

type Repository interface {
	basedomain.RepositoryCommand[CreateApiRoute, UpdateApiRoute]
	basedomain.RepositoryQuery[ApiRoute]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateApiRoute, UpdateApiRoute]
	basedomain.UseCaseQuery[ApiRoute]
	basedomain.UseCaseQueryRelation[ApiRouteRelation]
	basedomain.UseCaseTx[Service]
}
