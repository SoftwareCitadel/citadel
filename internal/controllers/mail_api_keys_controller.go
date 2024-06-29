package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	mailsPages "citadel/views/concerns/mails/pages"

	caesarAuth "github.com/caesar-rocks/auth"
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
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	mailDomains, err := c.mailDomainsRepo.FindAllFromUser(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	apiKeys, err := c.mailApiKeysRepo.FindAllFromUserWithRelatedDomain(ctx.Context(), user.ID)
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
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Verify the domain ownership, if domain specified.
	if data.DomainID != "" {
		domain, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", data.DomainID)
		if err != nil {
			return err
		}

		if domain.UserID != user.ID {
			return caesar.NewError(403)
		}
	}

	var onboarding bool

	if data.Name == "" {
		data.Name = "Onboarding"
		onboarding = true
	}

	apiKey := &models.MailApiKey{Name: data.Name, UserID: user.ID, MailDomainID: data.DomainID}
	if err := c.mailApiKeysRepo.Create(ctx.Context(), apiKey); err != nil {
		return err
	}

	if onboarding {
		return ctx.Render(mailsPages.AddApiKeyOnboarding(apiKey.Value))
	}

	return ctx.Redirect("/mails/api_keys")
}

func (c *MailApiKeysController) Update(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	apiKey, err := c.mailApiKeysRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return err
	}

	if apiKey.UserID != user.ID {
		return caesar.NewError(403)
	}

	data, _, ok := caesar.Validate[StoreMailApiKeyValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Verify the domain ownership, if domain specified.
	if data.DomainID != "" {
		domain, err := c.mailDomainsRepo.FindOneBy(ctx.Context(), "id", data.DomainID)
		if err != nil {
			return err
		}

		if domain.UserID != user.ID {
			return caesar.NewError(403)
		}
	}

	apiKey.Name = data.Name
	apiKey.MailDomainID = data.DomainID

	if err := c.mailApiKeysRepo.UpdateOneWhere(ctx.Context(), "id", ctx.PathValue("id"), apiKey); err != nil {
		return err
	}

	return ctx.Redirect("/mails/api_keys")
}

func (c *MailApiKeysController) Delete(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	apiKey, err := c.mailApiKeysRepo.FindOneBy(ctx.Context(), "id", ctx.PathValue("id"))
	if err != nil {
		return err
	}

	if apiKey.UserID != user.ID {
		return caesar.NewError(403)
	}

	if err := c.mailApiKeysRepo.DeleteOneWhere(ctx.Context(), "id", ctx.PathValue("id")); err != nil {
		return err
	}

	return ctx.Redirect("/mails/api_keys")
}
