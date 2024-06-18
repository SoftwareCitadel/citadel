package controllers

import (
	"citadel/app/models"
	"citadel/app/repositories"
	mailsPages "citadel/views/pages/mails"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
)

type MailDomainsController struct {
	mailDomainsRepo *repositories.MailDomainsRepository
}

func NewMailDomainsController(mailDomainsRepo *repositories.MailDomainsRepository) *MailDomainsController {
	return &MailDomainsController{mailDomainsRepo}
}

func (c *MailDomainsController) Index(ctx *caesar.CaesarCtx) error {
	return ctx.Render(mailsPages.DomainsPage())
}

type StoreMailDomainValidator struct {
	Domain string `json:"domain" validate:"required"`
}

func (c *MailDomainsController) Store(ctx *caesar.CaesarCtx) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	data, _, ok := caesar.Validate[StoreMailDomainValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	domain := &models.MailDomain{Domain: data.Domain, UserID: user.ID}

	if err := c.mailDomainsRepo.Create(ctx.Context(), domain); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Mail domain created successfully.")

	return ctx.Redirect("/mail_domains")
}

func (c *MailDomainsController) Show(ctx *caesar.CaesarCtx) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	domain, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(500)
	}

	if domain.UserID != user.ID {
		return caesar.NewError(403)
	}

	return ctx.Render(mailsPages.ShowDomainPage())
}

func (c *MailDomainsController) Delete(ctx *caesar.CaesarCtx) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	domain, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(500)
	}

	if domain.UserID != user.ID {
		return caesar.NewError(403)
	}

	if err := c.mailDomainsRepo.DeleteOneWhere(ctx.Context(), "id", ctx.PathValue("id")); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Mail domain deleted successfully.")

	return ctx.Redirect("/mail_domains")
}
