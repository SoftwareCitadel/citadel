package config

import (
	"citadel/views/pages/errors"

	"github.com/caesar-rocks/core"
)

func ProvideErrorHandler() *core.ErrorHandler {
	return &core.ErrorHandler{Handle: func(ctx *core.CaesarCtx, err error) {
		code := core.RetrieveErrorCode(err)
		if code == 404 {
			ctx.WithStatus(404).Render(errors.NotFoundPage())
		} else {
			ctx.WithStatus(code).Render(errors.ServerErrorPage(code))
		}
	}}
}
