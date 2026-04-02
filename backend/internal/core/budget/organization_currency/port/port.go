package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateOrganizationCurrency, UpdateOrganizationCurrency]
	basedomain.RepositoryQuery[OrganizationCurrency]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateOrganizationCurrency, UpdateOrganizationCurrency]
	basedomain.UseCaseQuery[OrganizationCurrency]
	basedomain.UseCaseTx[Service]
}
