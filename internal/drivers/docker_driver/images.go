package dockerDriver

import (
	"context"

	"github.com/docker/docker/api/types/image"
)

func (driver *DockerDriver) ImageExists(imageName string) (bool, error) {
	smr, err := driver.Client.ImageList(
		context.Background(),
		image.ListOptions{
			All: true,
		},
	)
	if err != nil {
		return false, err
	}

	for _, image := range smr {
		for _, tag := range image.RepoTags {
			if tag == imageName+":latest" {
				return true, nil
			}
		}
	}

	return false, nil
}
