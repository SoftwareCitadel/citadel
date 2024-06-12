package dockerDriver

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/swarm"
)

func (driver *DockerDriver) initializeSwarm() error {
	log.Println("Initializing swarm if not initialized...")

	if _, err := driver.Client.SwarmInspect(context.Background()); err == nil {
		return err
	}

	_, err := driver.Client.SwarmInit(context.Background(), swarm.InitRequest{
		ListenAddr: "192.168.65.4:2377",
	})
	if err == nil {
		return err
	}

	return nil
}
