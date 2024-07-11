package smtpBackend

import (
	mailBuilder "citadel/internal/mail_builder"
	"citadel/internal/repositories"
	smtpSender "citadel/internal/smtp_sender"
	"citadel/util"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
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
	emailsRepo  *repositories.EmailRepository
}

// New creates a new Backend.
func New(
	apiKeysRepo *repositories.MailApiKeysRepository,
	domainsRepo *repositories.MailDomainsRepository,
	emailsRepo *repositories.EmailRepository,
) *Backend {
	return &Backend{apiKeysRepo, domainsRepo, emailsRepo}
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		apiKeysRepo: bkd.apiKeysRepo,
		domainsRepo: bkd.domainsRepo,
		emailsRepo:  bkd.emailsRepo,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	apiKeysRepo *repositories.MailApiKeysRepository
	domainsRepo *repositories.MailDomainsRepository
	emailsRepo  *repositories.EmailRepository

	from      string
	to        string
	apiKey    *models.MailApiKey
	domain    *models.MailDomain
	newAPIKey string
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
		if password == "" {
			//Generate a random password
			password, err = PasswordGenerator(32)
			if err != nil {
				return errors.New("failed to create random password")
			}
			// Hash the password for storage in the database
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
			s.newAPIKey = string(newAPIKey)
			s.apiKey = newMailApiKey
			// Set flag to indicate that this API key is newly issued
			s.apiKey.IsNewlyCreated = true
			return nil
		}
		return errors.New("invalid api key")
	}), nil
}

// New function to determine if API key is new or not
func (s *Session) GetNewAPIKeyInfo() (string, bool) {
	if s.apiKey.IsNewlyCreated {
		s.apiKey.IsNewlyCreated = false // reset the flag for future use
		return s.newAPIKey, true
	}
	return "", false
}

// Handle API Key Auth and update client side information
func HandleAuth(w http.ResponseWriter, r *http.Request) {
	s := &Session{}
	server, err := s.Auth("PLAIN")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Print(server)
	// Check if a new API Key was created
	newAPIKey, isNew := s.GetNewAPIKeyInfo()
	if isNew {
		// Return the new API key information as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "new_key_created",
			"apiKey": newAPIKey,
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "authentication_successful",
		})
	}
}

// Generate a random Password
func PasswordGenerator(length int) (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		all       = lowercase + uppercase + digits + symbols
	)

	var password strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		password.WriteByte(all[n.Int64()])
	}
	result := []byte(password.String())

	// Ensure at least one character from each category
	categories := []string{lowercase, uppercase, digits, symbols}
	for _, category := range categories {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(category))))
		if err != nil {
			return "", err
		}
		pos, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {
			return "", err
		}
		result[pos.Int64()] = category[n.Int64()]
	}

	return string(result), nil
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

	// Set the email params for storing to the database
	e := &models.Email{
		Sender:         s.from,
		Recipient:      s.to,
		OrganizationID: s.domain.OrganizationID,
		Subject:        msg.Header.Get("Subject"),
		Body:           string(outputMsg),
		CreatedAt:      (time.Now()),
		UpdatedAt:      (time.Now()),
	}
	// Store the email in the database
	err = s.emailsRepo.Create(context.Background(), e)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) HandleListSentEmails(w http.ResponseWriter, r *http.Request) {
	// Capture OrgID
	orgID := s.domain.OrganizationID
	if orgID == "" {
		http.Error(w, "Organization ID is required", http.StatusBadRequest)
		return
	}

	// Capture additionaly query inputs from user

	filter := make(map[string]string)
	if sender := r.URL.Query().Get("sender"); sender != "" {
		filter["sender"] = sender
	}
	if recipient := r.URL.Query().Get("recipient"); recipient != "" {
		filter["recipient"] = recipient
	}

	// Set return limit
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20 // Default limit
	}

	// Set return offset
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	// Retrieve emails from the database
	emails, err := s.ListSentEmails(r.Context(), orgID, filter, limit, offset)
	if err != nil {
		log.Error("Error retrieving emails", "error", err)
		http.Error(w, "Failed to retrieve emails", http.StatusInternalServerError)
		return
	}
	// Craft the response
	response := struct {
		Status string         `json:"status"`
		Emails []models.Email `json:"emails"`
	}{
		Status: "success",
		Emails: emails,
	}
	// Serve the response to the rendering page
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error encoding response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Listener for retrieving emails
func (s *Session) ListSentEmails(ctx context.Context, orgID string, filter map[string]string, limit, offset int) ([]models.Email, error) {
	return s.emailsRepo.FindEmails(ctx, orgID, filter, limit, offset)
}

// Reset is called after RSET.
func (s *Session) Reset() {}

// Logout is called after QUIT.
func (s *Session) Logout() error {
	return nil
}
