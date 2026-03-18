package eventbus

import (
	"backend/core/notifications/eventbus/core"
	"backend/core/notifications/eventbus/port"
	basedomain "backend/port"
	"backend/adapter/di"

	"github.com/samber/do/v2"
)

func Module(i do.Injector) {
	di.Provide(i, func(i do.Injector) (port.EventBus, error) {
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(logger), nil
	})
}
