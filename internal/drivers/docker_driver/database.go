package dockerDriver

import (
	"citadel/internal/models"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func (driver *DockerDriver) CreateDatabase(db models.Database) error {
	cfg := buildConfig(db)
	resp, err := driver.Client.ContainerCreate(context.Background(), cfg, nil, nil, nil, "citadel-"+db.Slug)
	if err != nil {
		return err
	}
	return driver.Client.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
}

func (driver *DockerDriver) DeleteDatabase(db models.Database) error {
	containerID := fmt.Sprintf("citadel-%s", db.Slug)
	err := driver.Client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}

func buildConfig(db models.Database) *container.Config {
	image, envs, exposedPorts := prepareImageAndEnvsAndPorts(db)
	labels := prepareLabels(db)
	return &container.Config{
		Image:        image,
		Env:          envs,
		ExposedPorts: exposedPorts,
		Labels:       labels,
	}
}

func prepareImageAndEnvsAndPorts(db models.Database) (string, []string, nat.PortSet) {
	var image string
	var envs []string
	ports := nat.PortSet{}

	switch db.DBMS {
	case models.Postgres:
		image = "postgres:13-alpine"
		envs = []string{
			fmt.Sprintf("POSTGRES_DB=%s", db.Name),
			fmt.Sprintf("POSTGRES_USER=%s", db.Username),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", db.Password),
			fmt.Sprintf("PG_PASSWORD=%s", db.Password),
		}
		ports[nat.Port("5432/tcp")] = struct{}{}
	case models.MySQL:
		image = "mysql:8.3.0"
		envs = []string{
			fmt.Sprintf("MYSQL_DATABASE=%s", db.Name),
			fmt.Sprintf("MYSQL_USER=%s", db.Username),
			fmt.Sprintf("MYSQL_PASSWORD=%s", db.Password),
			"MYSQL_RANDOM_ROOT_PASSWORD=yes",
		}
		ports[nat.Port("3306/tcp")] = struct{}{}
	case models.Redis:
		image = "redis:6-alpine"
		envs = []string{
			fmt.Sprintf("REDIS_PASSWORD=%s", db.Password),
		}
		ports[nat.Port("6379/tcp")] = struct{}{}
	}

	return image, envs, ports
}

func prepareLabels(db models.Database) map[string]string {
	hostname := fmt.Sprintf("HostSNI(`%s`)", db.Host)
	routerName := fmt.Sprintf("citadel-builder-%s", db.Slug)
	port := 0
	switch db.DBMS {
	case models.Postgres:
		port = 5432
	case models.MySQL:
		port = 3306
	case models.Redis:
		port = 6379
	}

	return map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.tcp.routers.%s.rule", routerName):                      hostname,
		fmt.Sprintf("traefik.tcp.routers.%s.entrypoints", routerName):               string(db.DBMS),
		fmt.Sprintf("traefik.tcp.services.%s.loadbalancer.server.port", routerName): fmt.Sprintf("%d", port),
		fmt.Sprintf("traefik.tcp.routers.%s.tls", routerName):                       "true",
		fmt.Sprintf("traefik.tcp.routers.%s.tls.certresolver", routerName):          "letsencrypt",
	}
}
