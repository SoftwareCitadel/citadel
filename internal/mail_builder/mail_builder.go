package mailBuilder

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/emersion/go-msgauth/dkim"
	"github.com/rs/xid"
)

// MailBuilder is responsible for preparing an email message for sending.
type MailBuilder struct {
	msg         *mail.Message
	embedded    []*file // add embedding of files to MailBuilder
	attachments []*file // add attachments to MailBuilder
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

	// Create a multipart writer
	writer := multipart.NewWriter(&buf)
	mb.msg.Header["Content-Type"] = []string{"multipart/mixed; boundary=" + writer.Boundary()}

	// Write the initial headers
	for key, values := range mb.msg.Header {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
		}
	}

	buf.WriteString("\r\n")

	// Write the body
	bodyPart, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain"}})
	if err != nil {
		return nil, fmt.Errorf("failed to create body part: %w", err)
	}
	if _, err := io.Copy(bodyPart, mb.msg.Body); err != nil {
		return nil, fmt.Errorf("failed to copy body: %w", err)
	}

	// Write the attachments
	for _, attachment := range mb.attachments {
		attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {fmt.Sprintf("application/octet-stream; name=%q", attachment.Name)},
			"Content-Transfer-Encoding": {"base64"},
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=%q", attachment.Name)},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment part: %w", err)
		}

		encoder := base64.NewEncoder(base64.StdEncoding, attachmentPart)
		if err := attachment.CopyFunc(encoder); err != nil {
			return nil, fmt.Errorf("failed to copy attachment: %w", err)
		}
		encoder.Close()
	}

	// Close the multipart writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
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

// Set the file structure
type file struct {
	Name     string
	Header   map[string]string // Check for conflicts with msg header
	CopyFunc func(w io.Writer) error
}

// Set the file settings

type FileSettings func(*file)

// Gather the document(s)
func (m *MailBuilder) AppendFileToEmail(list []*file, name string, settings []FileSettings) []*file {
	f := &file{
		Name:   filepath.Base(name),
		Header: make(map[string]string), // Check for conflicts with msg header
		CopyFunc: func(w io.Writer) error {
			g, err := os.Open(name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, g); err != nil {
				g.Close()
				return err
			}
			return g.Close()
		},
	}
	for _, s := range settings {
		s(f)
	}
	if list == nil {
		return []*file{f}
	}
	return append(list, f)
}

// Attach the document(s)

func (m *MailBuilder) Attach(filename string, settings ...FileSettings) {
	m.attachments = m.AppendFileToEmail(m.attachments, filename, settings)
}

// Embed the document(s)
func (m *MailBuilder) Embed(filename string, settings ...FileSettings) {
	m.embedded = m.AppendFileToEmail(m.embedded, filename, settings)
}

// Append headers from sendmail
func (mb *MailBuilder) SetHeader(key, value string) {
    mb.msg.Header[key] = []string{value}
}
