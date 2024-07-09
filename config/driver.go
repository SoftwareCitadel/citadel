package config

import (
	"citadel/internal/drivers"
	dockerDriver "citadel/internal/drivers/docker_driver"
	ravelDriver "citadel/internal/drivers/ravel_driver"
	"citadel/internal/repositories"
)

func ProvideDriver(
	env *EnvironmentVariables,
	appsRepo *repositories.ApplicationsRepository,
	deplsRepo *repositories.DeploymentsRepository,
) drivers.Driver {
	switch env.DRIVER {
	case DockerDriver:
		return dockerDriver.New(appsRepo, deplsRepo)
	case RavelDriver:
		return ravelDriver.New()
	default:
		return nil
	}
}
