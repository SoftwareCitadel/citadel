package smtpBackend

import (
	mailBuilder "citadel/internal/mail_builder"
	"citadel/internal/repositories"
	smtpSender "citadel/internal/smtp_sender"
	"citadel/util"
	"context"
	"errors"
	"io"
	"net/mail"
	"os"

	"citadel/internal/models"

	"github.com/charmbracelet/log"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct {
	apiKeysRepo *repositories.MailApiKeysRepository
	domainsRepo *repositories.MailDomainsRepository
}

// New creates a new Backend.
func New(
	apiKeysRepo *repositories.MailApiKeysRepository,
	domainsRepo *repositories.MailDomainsRepository,
) *Backend {
	return &Backend{apiKeysRepo, domainsRepo}
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		apiKeysRepo: bkd.apiKeysRepo,
		domainsRepo: bkd.domainsRepo,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	apiKeysRepo *repositories.MailApiKeysRepository
	domainsRepo *repositories.MailDomainsRepository

	from   string
	to     string
	apiKey *models.MailApiKey
	domain *models.MailDomain
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		// Check if the user exists
		if username != "citadel" {
			return errors.New("invalid username")
		}

		// Check if the API key is valid
		apiKey, err := s.apiKeysRepo.FindOneBy(context.Background(), "value", password)
		if err != nil {
			return errors.New("invalid api key")
		}

		s.apiKey = apiKey

		return nil
	}), nil
}

// Mail is called after MAIL FROM.
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	// Check if the domain is verified, and if the user is allowed to send from it
	d, err := util.GetEmailDomain(from)
	if err != nil {
		log.Error("failed to get email domain", "error", err)
		return err
	}
	domain, err := s.domainsRepo.FindVerifiedDomainWithOrg(context.Background(), d, s.apiKey.OrganizationID)
	if err != nil {
		log.Error("failed to find domain", "domain", d, "s.apiKey.OrganizationID", s.apiKey.OrganizationID, "error", err)
		return err
	}

	// Save the domain and the sender
	s.from = from
	s.domain = domain

	return nil
}

// Rcpt is called after RCPT TO.
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// Save the recipient
	s.to = to

	return nil
}

// Data is called after DATA.
func (s *Session) Data(r io.Reader) error {
	// Read the message
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}

	// Build the email to send
	builder := mailBuilder.New(msg)
	outputMsg, err := builder.Build()
	if err != nil {
		log.Error("failed to build the email:", "error", err)
		return err
	}

	// Sign the email with DKIM
	outputMsg, err = builder.SignWithDKIM(outputMsg, s.domain.Domain, s.domain.DKIMPrivateKey)
	if err != nil {
		log.Error("failed to sign the email with DKIM:", "error", err)
		return err
	}

	// Send the email
	sender := smtpSender.New(os.Getenv("SMTP_DOMAIN"))
	sender.Send(s.from, s.to, outputMsg)

	return nil
}

// Reset is called after RSET.
func (s *Session) Reset() {}

// Logout is called after QUIT.
func (s *Session) Logout() error {
	return nil
}
