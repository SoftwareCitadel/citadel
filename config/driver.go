package config

import (
	"citadel/internal/drivers"
	dockerDriver "citadel/internal/drivers/docker_driver"
	"citadel/internal/repositories"
)

func ProvideDriver(appsRepo *repositories.ApplicationsRepository, deplsRepo *repositories.DeploymentsRepository) drivers.Driver {
	return dockerDriver.NewDockerDriver(appsRepo, deplsRepo)
}
