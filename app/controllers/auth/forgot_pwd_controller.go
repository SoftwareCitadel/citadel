package authControllers

import (
	"bytes"
	"citadel/app/repositories"
	"citadel/views/mails"
	authPages "citadel/views/pages/auth"
	"log/slog"

	caesar "github.com/caesar-rocks/core"
	mailer "github.com/caesar-rocks/mail"
)

type ForgotPwdController struct {
	mailer *mailer.Mailer
	repo   *repositories.UsersRepository
}

func NewForgotPwdController(mailer *mailer.Mailer, repo *repositories.UsersRepository) *ForgotPwdController {
	return &ForgotPwdController{mailer, repo}
}

func (c *ForgotPwdController) Show(ctx *caesar.CaesarCtx) error {
	return ctx.Render(authPages.ForgotPasswordPage())
}

type ForgotPwdValidator struct {
	Email string `form:"email" validate:"required,email"`
}

func (c *ForgotPwdController) Handle(ctx *caesar.CaesarCtx) error {
	data, _, ok := caesar.Validate[ForgotPwdValidator](ctx)
	if !ok {
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	user, _ := c.repo.FindOneBy(ctx.Context(), "email", data.Email)
	if user == nil {
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	var buf bytes.Buffer
	res := mails.ForgotPasswordMail(data.Email, "http://localhost:3000/reset-password")
	if err := res.Render(ctx.Context(), &buf); err != nil {
		slog.Error("Failed to render email", "err", err)
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	if err := c.mailer.Send(mailer.Mail{
		From:    "Software Citadel <contact@softwarecitadel.com>",
		To:      data.Email,
		Subject: "Reset your password",
		Html:    buf.String(),
	}); err != nil {
		slog.Error("Failed to send email", "err", err)
	}

	return ctx.Render(authPages.ForgotPasswordPage())
}
