package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	mailsPages "citadel/views/concerns/mails/pages"

	caesar "github.com/caesar-rocks/core"
)

type MailApiKeysController struct {
	mailDomainsRepo *repositories.MailDomainsRepository
	mailApiKeysRepo *repositories.MailApiKeysRepository
}

func NewMailApiKeysController(
	mailDomainsRepo *repositories.MailDomainsRepository,
	mailApiKeysRepo *repositories.MailApiKeysRepository,
) *MailApiKeysController {
	return &MailApiKeysController{mailDomainsRepo: mailDomainsRepo, mailApiKeysRepo: mailApiKeysRepo}
}

func (c *MailApiKeysController) Index(ctx *caesar.Context) error {
	mailDomains, err := c.mailDomainsRepo.FindAllFromOrg(ctx.Context(), ctx.PathValue("orgId"))
	if err != nil {
		return err
	}

	apiKeys, err := c.mailApiKeysRepo.FindAllFromOrgWithRelatedDomain(ctx.Context(), ctx.PathValue("orgId"))
	if err != nil {
		return err
	}

	return ctx.Render(mailsPages.APIKeysPage(mailDomains, apiKeys))
}

type StoreMailApiKeyValidator struct {
	Name     string `form:"name"`
	DomainID string `form:"domain_id"`
}

func (c *MailApiKeysController) Store(ctx *caesar.Context) error {
	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	var onboarding bool

	if data.Name == "" {
		data.Name = "Onboarding"
		onboarding = true
	}

	apiKey := &models.MailApiKey{Name: data.Name, OrganizationID: ctx.PathValue("orgId"), MailDomainID: data.DomainID}
	if err := c.mailApiKeysRepo.Create(ctx.Context(), apiKey); err != nil {
		return err
	}

	if onboarding {
		return ctx.Render(mailsPages.AddApiKeyOnboarding(apiKey.Value))
	}

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/api_keys")
}

func (c *MailApiKeysController) Update(ctx *caesar.Context) error {
	apiKey, err := c.mailApiKeysRepo.FindOneBy(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Verify the domain ownership, if domain specified.
	if data.DomainID != "" {
		if _, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", data.DomainID, "organization_id", ctx.PathValue("orgId")); err != nil {
			return err
		}
	}

	apiKey.Name = data.Name
	apiKey.MailDomainID = data.DomainID

	if err := c.mailApiKeysRepo.UpdateOneWhere(
		ctx.Context(), apiKey,
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	); err != nil {
		return err
	}

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/api_keys")
}

func (c *MailApiKeysController) Delete(ctx *caesar.Context) error {
	if err := c.mailApiKeysRepo.DeleteOneWhere(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"organization_id", ctx.PathValue("orgId"),
	); err != nil {
		return err
	}

	return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/mails/api_keys")
}
