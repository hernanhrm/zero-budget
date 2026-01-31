package domain

import (
	basedomain "backend/domain"
)

type Repository interface {
	basedomain.RepositoryCommand[CreateOrganization, UpdateOrganization]
	basedomain.RepositoryQuery[Organization]
}

type Service interface {
	basedomain.UseCaseCommand[CreateOrganization, UpdateOrganization]
	basedomain.UseCaseQuery[Organization]
	basedomain.UseCaseQueryRelation[OrganizationRelation]
}
