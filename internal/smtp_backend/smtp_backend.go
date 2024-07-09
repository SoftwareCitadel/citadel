package smtpBackend

import (
	mailBuilder "citadel/internal/mail_builder"
	"citadel/internal/repositories"
	smtpSender "citadel/internal/smtp_sender"
	"citadel/util"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

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
		// new check implementation
		apiKey, err := s.apiKeysRepo.FindOneBy(context.Background(), "value", password)
		if err != nil {
			// Verify the API key if found
			decodedHash, err := hex.DecodeString(apiKey.Value)
			if err != nil {
				return errors.New("invalid api key format")
			}
			if err := bcrypt.CompareHashAndPassword(decodedHash, []byte(password)); err == nil {
				s.apiKey = apiKey
				return nil
			}
		}
		//If the api wasn't found or didn't match, check if a new creation is needed
		if password == "CREATE_NEW_API_KEY" {
			//Generate a new API key
			newAPIKey, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return errors.New("failed to create new API key")
			}
			// Hash the key
			hashedKey, err := bcrypt.GenerateFromPassword([]byte(newAPIKey), bcrypt.DefaultCost)
			if err != nil {
				return errors.New("failed to hash API key")
			}
			// Encode the hash
			encodedHash := hex.EncodeToString(hashedKey)
			//Create a new MailApiKey object
			newMailApiKey := &models.MailApiKey{
				Name:           "New API Key",
				Value:          encodedHash,
				OrganizationID: s.apiKey.OrganizationID,
				CreatedAt:      (time.Now()),
				UpdatedAt:      (time.Now()),
			}
			// Store the hashed key in the database
			err = s.apiKeysRepo.Create(context.Background(), newMailApiKey)
			if err != nil {
				return errors.New("failed to store API key")
			}
			// Give the user the option to save the API Key
			fmt.Printf("Your new API key is: %s\n", newAPIKey)
			fmt.Println("Please store this key safely. It will not be shown again.")

			s.apiKey = newMailApiKey
			return nil
		}

		return errors.New("invalid api key")
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
		return err
	}

	// Sign the email with DKIM
	outputMsg, err = builder.SignWithDKIM(outputMsg, s.domain.Domain, s.domain.DKIMPrivateKey)
	if err != nil {
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
