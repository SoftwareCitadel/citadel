package config

import (
	"citadel/app/drivers"
	dockerDriver "citadel/app/drivers/docker_driver"
	"citadel/app/repositories"
)

func ProvideDriver(appsRepo *repositories.ApplicationsRepository, deplsRepo *repositories.DeploymentsRepository) drivers.Driver {
	return dockerDriver.NewDockerDriver(appsRepo, deplsRepo)
}
