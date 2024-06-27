package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	analyticsWebsitesPages "citadel/views/concerns/analytics_websites/pages"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
	"github.com/charmbracelet/log"
)

type AnalyticsWebsitesController struct {
	*caesar.BaseResourceController

	analRepo *repositories.AnalyticsWebsitesRepository
}

func NewAnalyticsWebsitesController(
	analRepo *repositories.AnalyticsWebsitesRepository,
) *AnalyticsWebsitesController {
	return &AnalyticsWebsitesController{analRepo: analRepo}
}

func (c *AnalyticsWebsitesController) Index(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	websites, err := c.analRepo.FindAllFromUser(ctx.Context(), user.ID)
	if err != nil {
		log.Error("err", err)
		return caesar.NewError(400)
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(websites)
	}

	return ctx.Render(analyticsWebsitesPages.IndexPage(websites))
}

type CreateAnalyticsWebsiteValidator struct {
	Name   string `form:"name" validate:"required,min=3"`
	Domain string `form:"domain" validate:"required,min=3"`
}

func (c *AnalyticsWebsitesController) Store(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[CreateAnalyticsWebsiteValidator](ctx)
	if !ok {
		return nil
	}

	website := &models.AnalyticsWebsite{
		UserID: user.ID,
		Name:   data.Name,
		Domain: data.Domain,
	}
	if err := c.analRepo.Create(ctx.Context(), website); err != nil {
		return caesar.NewError(400)
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(website)
	}

	return ctx.Redirect("/analytics_websites/" + website.ID)
}

func (c *AnalyticsWebsitesController) Show(ctx *caesar.Context) error {
	website, err := c.analRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(404)
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(website)
	}

	return ctx.Render(analyticsWebsitesPages.ShowPage(*website))
}

func (c *AnalyticsWebsitesController) Edit(ctx *caesar.Context) error {
	website, err := c.analRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(404)
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(website)
	}

	return ctx.Render(analyticsWebsitesPages.EditPage(*website))
}

type UpdateAnalyticsWebsiteValidator struct {
	Name   string `form:"name" validate:"required,min=3"`
	Domain string `form:"domain" validate:"required"`
}

func (c *AnalyticsWebsitesController) Update(ctx *caesar.Context) error {
	website, err := c.analRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(404)
	}

	data, errors, ok := caesar.Validate[UpdateAnalyticsWebsiteValidator](ctx)
	if !ok {
		return ctx.Render(analyticsWebsitesPages.AnalyticsWebsiteSettingsForm(*website, errors))
	}

	website.Name = data.Name
	website.Domain = data.Domain

	if err := c.analRepo.UpdateOneWhere(ctx.Context(), "id", website.ID, website); err != nil {
		return caesar.NewError(400)
	}

	toast.Success(ctx, "Website updated successfully")

	return ctx.Render(analyticsWebsitesPages.AnalyticsWebsiteSettingsForm(*website, nil))
}

func (c *AnalyticsWebsitesController) Delete(ctx *caesar.Context) error {
	website, err := c.analRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(404)
	}

	if err := c.analRepo.DeleteOneWhere(ctx.Context(), "id", website.ID); err != nil {
		return caesar.NewError(400)
	}

	return ctx.Redirect("/analytics_websites")
}

func (c *AnalyticsWebsitesController) Track(ctx *caesar.Context) error {
	return nil
}
