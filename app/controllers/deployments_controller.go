package controllers

import (
	"bytes"
	"citadel/app/models"
	"citadel/app/repositories"
	"citadel/app/services"
	"citadel/util"
	appsPages "citadel/views/pages/apps"
	"io"
	"log/slog"

	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/drive"
	"github.com/caesar-rocks/events"
	"github.com/rs/xid"
)

type DeploymentsController struct {
	appsService *services.AppsService
	appsRepo    *repositories.ApplicationsRepository
	deplRepo    *repositories.DeploymentsRepository
	drive       *drive.Drive
	emitter     *events.EventsEmitter
}

func NewDeploymentsController(appsService *services.AppsService, appsRepo *repositories.ApplicationsRepository, deplRepo *repositories.DeploymentsRepository, drive *drive.Drive, emitter *events.EventsEmitter) *DeploymentsController {
	return &DeploymentsController{appsService: appsService, appsRepo: appsRepo, deplRepo: deplRepo, drive: drive, emitter: emitter}
}

func (c *DeploymentsController) Index(ctx *caesar.CaesarCtx) error {
	app, err := c.appsService.GetAppOwnedByCurrentUser(ctx)
	if err != nil {
		return err
	}

	return ctx.Render(appsPages.DeploymentsPage(*app))
}

func (c *DeploymentsController) Store(ctx *caesar.CaesarCtx) error {
	app, err := c.appsService.GetAppOwnedByCurrentUser(ctx)
	if err != nil {
		return err
	}

	ctx.Request.ParseMultipartForm(10 << 20) // 10 MB
	tarball, _, err := ctx.Request.FormFile("tarball")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, tarball); err != nil {
		return err
	}

	depl := &models.Deployment{
		ID:            xid.New().String(),
		Application:   app,
		ApplicationID: app.ID,
		Status:        models.DeploymentStatusBuilding,
		Origin:        models.DeploymentOriginCli,
	}

	if err := c.drive.Use("s3").Put(depl.ID, buf.Bytes()); err != nil {
		return err
	}

	if err := c.deplRepo.Create(ctx.Context(), depl); err != nil {
		return err
	}

	bytes, err := util.EncodeJSON(depl)
	if err != nil {
		return err
	}
	c.emitter.Emit("deployments.created", bytes)

	return ctx.SendText("Deployment created")
}

func (c *DeploymentsController) List(ctx *caesar.CaesarCtx) error {
	app, err := c.appsService.GetAppOwnedByCurrentUser(ctx)
	if err != nil {
		slog.Warn("Failed to get app", "err", err, "slug", ctx.PathValue("slug"))
		return err
	}

	depls, err := c.deplRepo.FindAllFromApplication(ctx.Context(), app.ID)
	if err != nil {
		return err
	}

	return ctx.Render(appsPages.DeploymentsList(*app, depls))
}
