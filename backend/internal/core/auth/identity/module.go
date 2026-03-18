package identity

import (
	"backend/core/auth/identity/adapter/handler"
	"backend/core/auth/identity/core"
	"backend/core/auth/identity/port"
	organizationPort "backend/core/auth/organization/port"
	rolePort "backend/core/auth/role/port"
	userPort "backend/core/auth/user/port"
	workspacePort "backend/core/auth/workspace/port"
	workspaceMemberPort "backend/core/auth/workspace_member/port"
	eventbusPort "backend/core/notifications/eventbus/port"
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
		eventBus := di.MustInvoke[eventbusPort.EventBus](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(userSvc, orgSvc, workspaceSvc, roleSvc, workspaceMemberSvc, eventBus, logger, jwtSecret), nil
	})

	di.Provide(i, func(i do.Injector) (handler.HTTP, error) {
		svc := di.MustInvoke[port.Service](i)
		return handler.NewHTTP(svc), nil
	})
}
