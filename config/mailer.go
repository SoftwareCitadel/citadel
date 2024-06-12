package config

import mailer "github.com/caesar-rocks/mail"

func ProvideMailer(env *EnvironmentVariables) *mailer.Mailer {
	return mailer.NewMailer(mailer.MailCfg{
		APIService: mailer.RESEND,
		APIKey:     env.RESEND_KEY,
	})
}
