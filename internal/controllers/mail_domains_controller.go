package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	mailsPages "citadel/views/concerns/mails/pages"
	"fmt"

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
	domains, err := c.mailDomainsRepo.FindAllFromOrg(ctx.Context(), ctx.PathValue("orgId"))
	if err != nil {
		return caesar.NewError(500)
	}

	return ctx.Render(mailsPages.ListDomainsPage(domains))
}

type StoreMailDomainValidator struct {
	Domain string `form:"domain" validate:"required"`
}

func (c *MailDomainsController) Store(ctx *caesar.Context) error {
	data, _, ok := caesar.Validate[StoreMailDomainValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	domain := &models.MailDomain{Domain: data.Domain, OrganizationID: ctx.PathValue("orgId")}

	if err := c.mailDomainsRepo.Create(ctx.Context(), domain); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Mail domain created successfully.")

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/domains/" + domain.ID)
}

func (c *MailDomainsController) Show(ctx *caesar.Context) error {
	domain, err := c.mailDomainsRepo.FindOneBy(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	)
	if err != nil {
		return caesar.NewError(500)
	}

	return ctx.Render(mailsPages.ShowDomainPage(*domain))
}

func (c *MailDomainsController) Delete(ctx *caesar.Context) error {
	if err := c.mailDomainsRepo.DeleteOneWhere(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Mail domain deleted successfully.")

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/domains")
}

func (c *MailDomainsController) CheckDNS(ctx *caesar.Context) error {
	// Retrieve the domain from the bun database, where the domain matches the input
	domain, err := c.mailDomainsRepo.FindOneBy(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	)
	if err != nil {
		return caesar.NewError(500)
	}

	// Check the DNS records
	if err := domain.CheckDNS(); err != nil {
		return caesar.NewError(500)
	}

	// Save the domain
	fmt.Println("domain.DNSVerified", domain.DNSVerified)
	if err := c.mailDomainsRepo.UpdateOneWhere(ctx.Context(), domain, "id", ctx.PathValue("id")); err != nil {
		return caesar.NewError(500)
	}

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/domains/" + domain.ID)
}
