package dockerDriver

import "context"

func (d *DockerDriver) ContainerExists(containerName string) bool {
	info, err := d.Client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return false
	}
	return info.ID != ""
}
