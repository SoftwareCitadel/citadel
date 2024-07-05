package mailBuilder

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/mail"
	"time"

	"github.com/emersion/go-msgauth/dkim"
	"github.com/rs/xid"
)

// MailBuilder is responsible for preparing an email message for sending.
type MailBuilder struct {
	msg *mail.Message
}

// New creates a new MailBuilder.
func New(msg *mail.Message) *MailBuilder {
	return &MailBuilder{
		msg: msg,
	}
}

// Build creates a new email message, ready to be sent.
func (mb *MailBuilder) Build() ([]byte, error) {
	// Create a new message
	var buf bytes.Buffer

	// Set Message-ID header
	emailBase64 := base64.URLEncoding.EncodeToString([]byte(mb.msg.Header.Get("From")))
	mb.msg.Header["Message-Id"] = []string{fmt.Sprintf("<%s@%s>", xid.New().String(), emailBase64)}

	// Set date header
	mb.msg.Header["Date"] = []string{time.Now().Format(time.RFC1123Z)}

	// Write the initial headers
	for key, values := range mb.msg.Header {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
		}
	}

	// Write the body
	buf.WriteString("\r\n")
	if _, err := io.Copy(&buf, mb.msg.Body); err != nil {
		return nil, fmt.Errorf("failed to copy body: %w", err)
	}

	return buf.Bytes(), nil
}

// signWithDKIM signs the email message with DKIM.
func (mb *MailBuilder) SignWithDKIM(msg []byte, domain, privateKey string) ([]byte, error) {
	// Decode the private key
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Parse the private key
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get error keys
	var headerKeys []string
	for key := range mb.msg.Header {
		headerKeys = append(headerKeys, key)
	}

	// Set DKIM options
	dkimOpts := &dkim.SignOptions{
		Domain:                 domain,
		Selector:               "default",
		Signer:                 privKey,
		HeaderCanonicalization: dkim.CanonicalizationRelaxed,
		BodyCanonicalization:   dkim.CanonicalizationRelaxed,
		Hash:                   crypto.SHA256,
		HeaderKeys:             headerKeys,
	}

	// Sign the message
	var signedMsg bytes.Buffer
	if err := dkim.Sign(&signedMsg, bytes.NewReader(msg), dkimOpts); err != nil {
		return nil, fmt.Errorf("failed to sign DKIM: %w", err)
	}

	return signedMsg.Bytes(), nil
}
