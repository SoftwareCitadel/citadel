package controllers

import (
	mailsPages "citadel/views/pages/mails"

	caesar "github.com/caesar-rocks/core"
)

type MailDomainsController struct {
}

func NewMailDomainsController() *MailDomainsController {
	return &MailDomainsController{}
}

func (c *MailDomainsController) Index(ctx *caesar.CaesarCtx) error {
	return ctx.Render(mailsPages.DomainsPage())
}

func (c *MailDomainsController) Show(ctx *caesar.CaesarCtx) error {
	return ctx.Render(mailsPages.ShowDomainPage())
}
