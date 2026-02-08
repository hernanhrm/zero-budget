package auth

import (
	"backend/core/auth/adapter/handler"
	"backend/core/auth/core"
	"backend/core/auth/port"
	organizationPort "backend/core/organization/port"
	rolePort "backend/core/role/port"
	userPort "backend/core/user/port"
	workspacePort "backend/core/workspace/port"
	workspaceMemberPort "backend/core/workspace_member/port"
	basedomain "backend/port"
	"backend/adapter/di"
	"github.com/samber/do/v2"
)

func Module(i do.Injector, jwtSecret string) {
	di.Provide(i, func(i do.Injector) (port.Service, error) {
		userSvc := di.MustInvoke[userPort.Service](i)
		orgSvc := di.MustInvoke[organizationPort.Service](i)
		workspaceSvc := di.MustInvoke[workspacePort.Service](i)
		roleSvc := di.MustInvoke[rolePort.Service](i)
		workspaceMemberSvc := di.MustInvoke[workspaceMemberPort.Service](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(userSvc, orgSvc, workspaceSvc, roleSvc, workspaceMemberSvc, logger, jwtSecret), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		return handler.NewHTTP(svc), nil
	})
}
