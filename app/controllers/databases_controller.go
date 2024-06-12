package controllers

import (
	"citadel/views/pages"

	caesar "github.com/caesar-rocks/core"
)

type DatabasesController struct{}

func NewDatabasesController() *DatabasesController {
	return &DatabasesController{}
}

func (c *DatabasesController) Index(ctx *caesar.CaesarCtx) error {
	return ctx.Render(pages.DatabasesPage())
}
