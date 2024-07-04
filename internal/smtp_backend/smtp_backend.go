package smtpBackend

import (
	"citadel/internal/listeners"
	"citadel/internal/repositories"
	"citadel/util"
	"context"
	"errors"
	"io"
	"net/mail"

	"citadel/internal/models"

	"github.com/caesar-rocks/events"
	"github.com/charmbracelet/log"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct {
	apiKeysRepo *repositories.MailApiKeysRepository
	domainsRepo *repositories.MailDomainsRepository
	emitter     *events.EventsEmitter
}

// New creates a new Backend.
func New(
	apiKeysRepo *repositories.MailApiKeysRepository,
	domainsRepo *repositories.MailDomainsRepository,
	emitter *events.EventsEmitter,
) *Backend {
	return &Backend{apiKeysRepo, domainsRepo, emitter}
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		apiKeysRepo: bkd.apiKeysRepo,
		domainsRepo: bkd.domainsRepo,
		emitter:     bkd.emitter,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	apiKeysRepo *repositories.MailApiKeysRepository
	domainsRepo *repositories.MailDomainsRepository
	emitter     *events.EventsEmitter

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

	// Emit event
	event := &listeners.OutboundEmailEvent{
		Msg:    msg,
		Domain: s.domain,
		From:   s.from,
		To:     s.to,
	}

	bytes, err := util.EncodeJSON(event)
	if err != nil {
		return err
	}
	s.emitter.Emit("smtp.outbound_email", bytes)

	return nil
}

// Reset is called after RSET.
func (s *Session) Reset() {}

// Logout is called after QUIT.
func (s *Session) Logout() error {
	return nil
}
