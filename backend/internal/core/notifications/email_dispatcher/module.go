package email_dispatcher

import (
	emailLogPort "backend/core/notifications/email_log/port"
	emailTemplatePort "backend/core/notifications/email_template/port"
	eventbusPort "backend/core/notifications/eventbus/port"
	"backend/core/notifications/events"
	"backend/core/notifications/email_dispatcher/core"
	"backend/core/notifications/email_dispatcher/port"
	basedomain "backend/port"
	"backend/adapter/di"

	"github.com/samber/do/v2"
)

func Module(i do.Injector, resendAPIKey, resendFromAddress string) {
	di.Provide(i, func(i do.Injector) (port.EmailSender, error) {
		logger := di.MustInvoke[basedomain.Logger](i)
		if resendAPIKey != "" {
			return core.NewResendSender(resendAPIKey, resendFromAddress, logger), nil
		}
		logger.Warn("RESEND_API_KEY not set, using log-only email sender")
		return core.NewLogSender(logger), nil
	})

	di.Provide(i, func(i do.Injector) (port.Service, error) {
		emailTemplateSvc := di.MustInvoke[emailTemplatePort.Service](i)
		emailLogSvc := di.MustInvoke[emailLogPort.Service](i)
		emailSender := di.MustInvoke[port.EmailSender](i)
		logger := di.MustInvoke[basedomain.Logger](i)
		return core.New(emailTemplateSvc, emailLogSvc, emailSender, logger), nil
	})

	// Subscribe to events
	bus := di.MustInvoke[eventbusPort.EventBus](i)
	svc := di.MustInvoke[port.Service](i)

	bus.Subscribe(events.UserSignedUp, svc.HandleUserSignedUp)
	bus.Subscribe(events.UserVerificationEmail, svc.HandleUserVerificationEmail)
	bus.Subscribe(events.UserPasswordReset, svc.HandleUserPasswordReset)
	bus.Subscribe(events.OrganizationInvitationCreated, svc.HandleOrganizationInvitationCreated)
}
