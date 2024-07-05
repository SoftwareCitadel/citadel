package apiControllers

import (
	"os"
	"strings"

	caesar "github.com/caesar-rocks/core"
	gomail "gopkg.in/mail.v2"
)

type EmailsController struct{}

func NewEmailsController() *EmailsController {
	return &EmailsController{}
}

type SendEmailValidator struct {
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

func (c *EmailsController) Send(ctx *caesar.Context) error {
	data, errors, ok := caesar.Validate[SendEmailValidator](ctx)
	if !ok {
		return ctx.SendJSON(struct {
			Errors map[string]string `json:"errors"`
		}{Errors: errors}, 400)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	m.SetBody("text/plain", data.Text)
	if data.HTML != "" {
		m.AddAlternative("text/html", data.HTML)
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return caesar.NewError(401)
	}
	splitted := strings.Split(authHeader, "Bearer ")
	if len(splitted) != 2 {
		return caesar.NewError(401)
	}

	d := gomail.NewDialer(os.Getenv("SMTP_DOMAIN"), 465, "citadel", splitted[1])

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
