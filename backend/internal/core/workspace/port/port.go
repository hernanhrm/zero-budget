package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryCommand[CreateWorkspace, UpdateWorkspace]
	basedomain.RepositoryQuery[Workspace]
	basedomain.RepositoryTx[Repository]
}

type Service interface {
	basedomain.UseCaseCommand[CreateWorkspace, UpdateWorkspace]
	basedomain.UseCaseQuery[Workspace]
	basedomain.UseCaseQueryRelation[WorkspaceRelation]
	basedomain.UseCaseTx[Service]
}
