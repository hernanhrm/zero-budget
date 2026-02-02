package domain

import basedomain "backend/domain"

type Repository interface {
	basedomain.RepositoryCommand[CreateEmailLog, UpdateEmailLog]
	basedomain.RepositoryQuery[EmailLog]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateEmailLog, UpdateEmailLog]
	basedomain.UseCaseQuery[EmailLog]
	basedomain.UseCaseTx[Service]
}
