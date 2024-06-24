package authControllers

import (
	"bytes"
	"citadel/app/repositories"
	authPages "citadel/views/concerns/auth/pages"
	"citadel/views/mails"
	"log/slog"
	"os"
	"time"

	caesar "github.com/caesar-rocks/core"
	mailer "github.com/caesar-rocks/mail"
	"github.com/golang-jwt/jwt"
)

type ForgotPwdController struct {
	mailer *mailer.Mailer
	repo   *repositories.UsersRepository
}

func NewForgotPwdController(mailer *mailer.Mailer, repo *repositories.UsersRepository) *ForgotPwdController {
	return &ForgotPwdController{mailer, repo}
}

func (c *ForgotPwdController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.ForgotPasswordPage())
}

type ForgotPwdValidator struct {
	Email string `form:"email" validate:"required,email"`
}

func (c *ForgotPwdController) Handle(ctx *caesar.Context) error {
	data, _, ok := caesar.Validate[ForgotPwdValidator](ctx)
	if !ok {
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	user, _ := c.repo.FindOneBy(ctx.Context(), "email", data.Email)
	if user == nil {
		return ctx.Render(authPages.ForgotPasswordSuccessAlert())
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("APP_KEY")))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	res := mails.ForgotPasswordMail(user.FullName, os.Getenv("APP_URL")+"/auth/reset_password/"+tokenString)
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

	return ctx.Render(authPages.ForgotPasswordSuccessAlert())
}
