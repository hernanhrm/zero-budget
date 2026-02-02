package domain

import basedomain "backend/domain"

type Repository interface {
	basedomain.RepositoryCommand[CreateEmailTemplate, UpdateEmailTemplate]
	basedomain.RepositoryQuery[EmailTemplate]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateEmailTemplate, UpdateEmailTemplate]
	basedomain.UseCaseQuery[EmailTemplate]
	basedomain.UseCaseTx[Service]
}
