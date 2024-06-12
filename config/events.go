package config

import (
	"citadel/app/listeners"

	"github.com/caesar-rocks/events"
)

func RegisterEventsEmitter(
	usersListener *listeners.UsersListener,
	deploymentsListener *listeners.DeploymentsListener,
) *events.EventsEmitter {
	emitter := events.NewEventsEmitter()
	emitter.On("users.created", usersListener.OnCreated)
	emitter.On("deployments.created", deploymentsListener.OnCreated)

	return emitter
}
