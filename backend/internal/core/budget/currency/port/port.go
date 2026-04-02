package port

import basedomain "backend/port"

type Repository interface {
	basedomain.RepositoryQuery[Currency]
}

type Service interface {
	basedomain.UseCaseQuery[Currency]
}
