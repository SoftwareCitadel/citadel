package apiControllers

import (
	caesar "github.com/caesar-rocks/core"
)

type EmailsController struct{}

func NewEmailsController() *EmailsController {
	return &EmailsController{}
}

func (c *EmailsController) Send(ctx *caesar.CaesarCtx) error {
	return nil
}
