package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	mailsPages "citadel/views/concerns/mails/pages"

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

func (c *MailDomainsController) Index(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	domains, err := c.mailDomainsRepo.FindAllFromUser(ctx.Context(), user.ID)
	if err != nil {
		return caesar.NewError(500)
	}

	return ctx.Render(mailsPages.ListDomainsPage(domains))
}

type StoreMailDomainValidator struct {
	Domain string `form:"domain" validate:"required"`
}

func (c *MailDomainsController) Store(ctx *caesar.Context) error {
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

	return ctx.Redirect("/mails/domains")
}

func (c *MailDomainsController) Show(ctx *caesar.Context) error {
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

	return ctx.Render(mailsPages.ShowDomainPage(*domain))
}

func (c *MailDomainsController) Delete(ctx *caesar.Context) error {
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

	return ctx.Redirect("/mails/domains")
}

func (c *MailDomainsController) CheckDNS(ctx *caesar.Context) error {
	// Retrieve the user from the context
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	// Retrieve the domain from the bun database, where the domain matches the input
	domain, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return caesar.NewError(500)
	}

	// Check domain ownership
	if domain.UserID != user.ID {
		return caesar.NewError(403)
	}

	// Check the DNS records
	if err := domain.CheckDNS(); err != nil {
		return caesar.NewError(500)
	}

	return ctx.Redirect("/mails/domains")
}
