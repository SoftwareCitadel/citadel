package config

import (
	"citadel/internal/listeners"

	"github.com/caesar-rocks/events"
)

func RegisterEventsEmitter(
	usersListener *listeners.UsersListener,
	deploymentsListener *listeners.DeploymentsListener,
	smtpListener *listeners.SMTPListener,
) *events.EventsEmitter {
	emitter := events.NewEventsEmitter()
	emitter.On("users.created", usersListener.OnCreated)
	emitter.On("deployments.created", deploymentsListener.OnCreated)
	emitter.On("smtp.outbound_email", smtpListener.OnOutboundEmail)

	return emitter
}
