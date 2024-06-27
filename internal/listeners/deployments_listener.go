package listeners

import (
	"citadel/internal/drivers"
	"citadel/internal/models"
	"citadel/util"

	"github.com/ThreeDotsLabs/watermill/message"
)

type DeploymentsListener struct {
	driver drivers.Driver
}

func NewDeploymentsListener(driver drivers.Driver) *DeploymentsListener {
	return &DeploymentsListener{driver}
}

func (deplListener *DeploymentsListener) OnCreated(msg *message.Message) ([]*message.Message, error) {
	var deployment models.Deployment
	if err := util.DecodeJSON(msg.Payload, &deployment); err != nil {
		return nil, err
	}

	if err := deplListener.driver.IgniteBuilder(*deployment.Application, deployment); err != nil {
		return nil, err
	}

	return nil, nil
}
