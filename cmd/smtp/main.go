package main

import (
	"citadel/config"
	"citadel/internal/listeners"
	"citadel/internal/repositories"
	smtpBackend "citadel/internal/smtp_backend"
	"context"
	"crypto/tls"
	"log/slog"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/caesar-rocks/events"
	"github.com/charmbracelet/log"
	"github.com/emersion/go-smtp"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.ProvideEnvironmentVariables, config.ProvideDatabase),
		fx.Provide(repositories.NewMailApiKeysRepository, repositories.NewMailDomainsRepository),
		fx.Provide(listeners.NewSMTPListener),
		fx.Provide(newSMTPServer, newEventsEmitter),
		fx.Invoke(func(*smtp.Server) {}),
	).Run()
}

func newSMTPServer(
	lc fx.Lifecycle,
	env *config.EnvironmentVariables,
	apiKeysRepo *repositories.MailApiKeysRepository, domainsRepo *repositories.MailDomainsRepository,
	emitter *events.EventsEmitter,
) *smtp.Server {
	// Create the server
	backend := smtpBackend.New(apiKeysRepo, domainsRepo, emitter)

	srv := smtp.NewServer(backend)
	srv.Addr = env.SMTP_ADDR
	srv.Domain = env.SMTP_DOMAIN
	srv.WriteTimeout = 10 * time.Second
	srv.ReadTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = false

	// Set the TLS configuration
	tlsConfig, err := certmagic.TLS([]string{env.SMTP_DOMAIN})
	if err != nil {
		log.Fatal("failed to get TLS configuration", "error", err)
	}
	tlsConfig.ClientAuth = tls.RequestClientCert
	tlsConfig.NextProtos = []string{"smtp"}

	srv.TLSConfig = tlsConfig

	// Register the server with the lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("SMTP server listening on", "addr", srv.Addr)
			go func() {
				if err := srv.ListenAndServeTLS(); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func newEventsEmitter(smtpListener *listeners.SMTPListener) *events.EventsEmitter {
	emitter := events.NewEventsEmitter()
	emitter.On("smtp.outbound_email", smtpListener.OnOutboundEmail)

	return emitter
}
