package dockerDriver

import (
	"citadel/internal/models"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
)

func (driver *DockerDriver) watchEvents() {
	eventsChan, errChan := driver.Client.Events(context.Background(), types.EventsOptions{})

	go func() {
		for {
			select {
			case err := <-errChan:
				slog.Error("Error while watching events", "error", err)
			case event := <-eventsChan:
				if err := driver.handleEvent(event); err != nil {
					slog.Error("Error while handling event", "error", err)
				}
			}
		}
	}()
}

func (driver *DockerDriver) handleEvent(event events.Message) error {
	if event.Type != "container" {
		return nil
	}

	image := event.Actor.Attributes["image"]

	if image == os.Getenv("BUILDER_IMAGE") && event.Action == "die" {
		depl, err := driver.DeplsRepo.FindOneByIdWithRelatedAppAndCerts(context.Background(), event.Actor.Attributes["name"])
		if err != nil {
			slog.Error("Error while finding deployment", "error", err, "deployment_id", event.Actor.Attributes["name"])
			return err
		}

		if event.Actor.Attributes["exitCode"] != "0" {
			return driver.handleBuildFailed(depl)
		}

		return driver.handleBuildSuccess(depl)
	} else {
		// Get deployment
		depl, err := driver.DeplsRepo.FindOneBy(context.Background(), "id", event.Actor.Attributes["deployment_id"])
		if err != nil {
			return err
		}

		if event.Action == "die" {
			depl.Status = models.DeploymentStatusDeployFailed
			if err := driver.DeplsRepo.UpdateOneWhere(context.Background(), depl, "id", depl.ID); err != nil {
				return err
			}
		} else if event.Action == "start" {
			depl.Status = models.DeploymentStatusSuccess
			if err := driver.DeplsRepo.UpdateOneWhere(context.Background(), depl, "id", depl.ID); err != nil {
				return err
			}
		}

		slog.Info("Unknown event", "event", fmt.Sprintf("%+v", event))
	}

	return nil
}

func (driver *DockerDriver) handleBuildFailed(depl *models.Deployment) error {
	depl.Status = models.DeploymentStatusBuildFailed
	if err := driver.DeplsRepo.UpdateOneWhere(context.Background(), depl, "id", depl.ID); err != nil {
		return err
	}

	return nil
}

func (driver *DockerDriver) handleBuildSuccess(depl *models.Deployment) error {
	depl.Status = models.DeploymentStatusDeploying
	if err := driver.DeplsRepo.UpdateOneWhere(context.Background(), depl, "id", depl.ID); err != nil {
		return err
	}

	if err := driver.IgniteApplication(*depl.Application, *depl); err != nil {
		return err
	}

	return nil
}
