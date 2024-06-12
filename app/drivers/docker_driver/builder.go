package dockerDriver

import (
	"citadel/app/models"
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
)

func (driver *DockerDriver) IgniteBuilder(app models.Application, depl models.Deployment) error {
	ct, err := driver.Client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: os.Getenv("BUILDER_IMAGE"),
			Env:   prepareBuilderEnv(app.ID, depl.ID),
			Labels: map[string]string{
				"traefik.enable": "false",
			},
		},
		&container.HostConfig{AutoRemove: false, NetworkMode: "host", Privileged: true},
		nil,
		nil,
		depl.ID,
	)
	if err != nil {
		return err
	}

	if err := driver.Client.ContainerStart(
		context.Background(),
		ct.ID,
		container.StartOptions{},
	); err != nil {
		return err
	}

	return nil
}

func prepareBuilderEnv(appID string, deplID string) []string {
	envList := []string{}
	envList = append(envList, "IMAGE_NAME="+appID)
	envList = append(envList, "FILE_NAME="+deplID)
	envList = append(envList, "REGISTRY_HOST="+os.Getenv("REGISTRY_HOST"))
	envList = append(envList, "REGISTRY_TOKEN="+os.Getenv("REGISTRY_TOKEN"))
	envList = append(envList, "S3_ENDPOINT="+os.Getenv("S3_ENDPOINT"))
	envList = append(envList, "S3_ACCESS_KEY_ID="+os.Getenv("S3_KEY"))
	envList = append(envList, "S3_SECRET_ACCESS_KEY="+os.Getenv("S3_SECRET"))
	envList = append(envList, "S3_BUCKET_NAME="+os.Getenv("S3_BUCKET"))
	envList = append(envList, "S3_REGION="+os.Getenv("S3_REGION"))
	return envList
}
