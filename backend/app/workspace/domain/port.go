package domain

import basedomain "backend/domain"

type Repository interface {
	basedomain.RepositoryCommand[CreateWorkspace, UpdateWorkspace]
	basedomain.RepositoryQuery[Workspace]
}

type Service interface {
	basedomain.UseCaseCommand[CreateWorkspace, UpdateWorkspace]
	basedomain.UseCaseQuery[Workspace]
	basedomain.UseCaseQueryRelation[WorkspaceRelation]
}
