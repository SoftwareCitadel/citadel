package controllers

import (
	mailsPages "citadel/views/pages/mails"

	caesar "github.com/caesar-rocks/core"
)

type MailApiKeysController struct{}

func NewMailApiKeysController() *MailApiKeysController {
	return &MailApiKeysController{}
}

func (c *MailApiKeysController) Index(ctx *caesar.CaesarCtx) error {
	return ctx.Render(mailsPages.APIKeysPage())
}
