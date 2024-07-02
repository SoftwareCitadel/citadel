package controllers

import (
	"citadel/internal/repositories"
	"citadel/internal/services"

	appsPages "citadel/views/concerns/apps/pages"

	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
)

type EnvController struct {
	appsService *services.AppsService
	repo        *repositories.ApplicationsRepository
}

func NewEnvController(appsService *services.AppsService, repo *repositories.ApplicationsRepository) *EnvController {
	return &EnvController{appsService, repo}
}

func (c *EnvController) Edit(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	return ctx.Render(appsPages.EnvPage(*app))
}

func (c *EnvController) Update(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	ctx.Request.ParseForm()
	env := ctx.Request.FormValue("env")
	app.Env = []byte(env)

	if err := c.repo.UpdateOneWhere(ctx.Context(), app, "id", app.ID); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Environment variables updated successfully.")

	return ctx.Render(appsPages.EnvForm(*app))
}
