package auth

import (
	"backend/app/auth/domain"
	"backend/app/auth/handler"
	"backend/app/auth/service"
	organizationDomain "backend/app/organization/domain"
	roleDomain "backend/app/role/domain"
	userDomain "backend/app/user/domain"
	workspaceDomain "backend/app/workspace/domain"
	workspaceMemberDomain "backend/app/workspace_member/domain"
	basedomain "backend/domain"
	"backend/infra/di"
	"github.com/samber/do/v2"
)

func Module(i do.Injector, jwtSecret string) {
	di.Provide(i, func(i do.Injector) (domain.Service, error) {
		userSvc := di.MustInvoke[userDomain.Service](i)
		orgSvc := di.MustInvoke[organizationDomain.Service](i)
		workspaceSvc := di.MustInvoke[workspaceDomain.Service](i)
		roleSvc := di.MustInvoke[roleDomain.Service](i)
		workspaceMemberSvc := di.MustInvoke[workspaceMemberDomain.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return service.New(userSvc, orgSvc, workspaceSvc, roleSvc, workspaceMemberSvc, logger, jwtSecret), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[domain.Service](i)
		return handler.NewHTTP(svc), nil
	})
}
