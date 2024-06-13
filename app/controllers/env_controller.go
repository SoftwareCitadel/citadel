package controllers

import (
	"citadel/app/repositories"
	"citadel/app/services"

	appsPages "citadel/views/pages/apps"

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

func (c *EnvController) Edit(ctx *caesar.CaesarCtx) error {
	app, err := c.appsService.GetAppOwnedByCurrentUser(ctx)
	if err != nil {
		return err
	}

	return ctx.Render(appsPages.EnvPage(*app))
}

func (c *EnvController) Update(ctx *caesar.CaesarCtx) error {
	app, err := c.appsService.GetAppOwnedByCurrentUser(ctx)
	if err != nil {
		return err
	}

	ctx.Request.ParseForm()
	env := ctx.Request.FormValue("env")
	app.Env = []byte(env)

	if err := c.repo.UpdateOneWhere(ctx.Context(), "id", app.ID, app); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Environment variables updated successfully.")

	return ctx.Render(appsPages.EnvForm(*app))
}
