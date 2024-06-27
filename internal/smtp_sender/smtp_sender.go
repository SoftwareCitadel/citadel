package smtpSender

import (
	"citadel/util"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/net/idna"
)

const (
	SMTP_DIAL_TIMEOUT  = 15 * time.Second
	SMTP_TOTAL_TIMEOUT = 120 * time.Second

	DEFAULT_SMTP_PORT = "25"
)

// New construct a new sender for a given hostname
func New(hostname string) *Sender {
	return &Sender{
		Hostname: hostname,
	}
}

// SenderError represent an smtp error
type SenderError interface {
	Error() string
	IsPermanent() bool
	Code() int
}

type smtpError struct {
	err         error
	isPermanent bool
	code        int
}

func (e smtpError) Error() string {
	return e.err.Error()
}

func (e smtpError) IsPermanent() bool {
	return e.isPermanent
}

func (e smtpError) Code() int {
	return e.code
}

func newSMTPError(err error, isPermanent bool, code int) *smtpError {
	return &smtpError{
		err:         err,
		isPermanent: isPermanent,
		code:        code,
	}
}

type Sender struct {
	Hostname string
}

// Send email
func (s *Sender) Send(from, to string, msg []byte) SenderError {
	toDomain, err := util.GetEmailDomain(to)
	log.Printf("domain %v\n", toDomain)
	if err != nil {
		return newSMTPError(err, true, 510)
	}

	mxs, lerr := lookupMXs(toDomain)
	if lerr != nil {
		return lerr
	}

	var lastErr *smtpError
	for _, mx := range mxs {
		err := deliver(from, to, msg, mx, false, s.Hostname)
		if err == nil {
			return nil
		}
		if err.code > 200 {
			log.Error("error sending email to, cannot retry other MXs", "to", to, "error", err)
			return err
		}

		lastErr = err
	}
	err = fmt.Errorf("all MXs failed, last error: %v", lastErr)
	return newSMTPError(err, false, lastErr.Code())
}

func deliver(from, to string, msg []byte, mx string, insecure bool, domain string) *smtpError {
	smtpURL := fmt.Sprintf("%v:%v", mx, DEFAULT_SMTP_PORT)

	conn, err := net.DialTimeout("tcp", smtpURL, SMTP_DIAL_TIMEOUT)
	if err != nil {
		log.Error("could not dial", "error", err)
		return newSMTPError(err, false, 111)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(SMTP_TOTAL_TIMEOUT)); err != nil {
		log.Error("cannot not set deadline", "error", err)
		return newSMTPError(err, false, 111)
	}

	c, err := smtp.NewClient(conn, mx)
	if err != nil {
		log.Debugf("Error creating client: %v", err)
		return newSMTPError(err, false, 111)
	}

	if err = c.Hello(domain); err != nil {
		log.Debugf("Error saying hello: %v", err)
		return newSMTPError(err, false, 111)
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName:         mx,
			InsecureSkipVerify: insecure,
		}
		err = c.StartTLS(config)
		if err != nil {
			// Unfortunately, many servers use self-signed certs, so if we
			// fail verification we just try again without validating.
			if insecure {
				log.Error("could not start tls", "error", err)
				return newSMTPError(err, false, 111)
			}
			log.Debug("tls error, retrying insecurely")
			return deliver(from, to, msg, mx, true, domain)
		}
	}

	if err := c.Mail(from); err != nil {
		return newSMTPErrorFromSTMP(err)
	}

	if err := c.Rcpt(to); err != nil {
		return newSMTPErrorFromSTMP(err)
	}

	w, err := c.Data()
	if err != nil {
		return newSMTPErrorFromSTMP(err)
	}

	if _, err := w.Write(msg); err != nil {
		return newSMTPErrorFromSTMP(err)
	}

	if err := w.Close(); err != nil {
		log.Debugf("err: %v\n", err)
		return newSMTPErrorFromSTMP(err)
	}

	if err := c.Quit(); err != nil {
		return newSMTPErrorFromSTMP(err)
	}

	return nil
}

func lookupMXs(domain string) ([]string, *smtpError) {
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return nil, newSMTPError(err, true, 512)
	}

	mxs := []string{}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		// TODO: Better handle Temporary errors.
		dnsErr, ok := err.(*net.DNSError)
		if !ok || !dnsErr.IsNotFound {
			return nil, newSMTPError(dnsErr, !dnsErr.Temporary(), 512)
		}
		// Permanent error, we assume MX does not exist and fall back to A.
		log.Debugf("failed to resolve MX for %s, falling back to A", domain)
		mxs = []string{domain}
	} else {
		// Convert the DNS records to a plain string slice. They're already
		// sorted by priority.
		for _, r := range mxRecords {
			mxs = append(mxs, r.Host)
		}
	}

	// Note that mxs could be empty; in that case we do NOT fall back to A.
	// This case is explicitly covered by the SMTP RFC.
	// https://tools.ietf.org/html/rfc5321#section-5.1

	// Cap the list of MXs to 5 hosts, to keep delivery attempt times
	// sane and prevent abuse.
	if len(mxs) > 5 {
		mxs = mxs[:5]
	}

	return mxs, nil
}

func newSMTPErrorFromSTMP(err error) *smtpError {
	terr, ok := err.(*textproto.Error)
	if !ok {
		return newSMTPError(err, false, 0)
	}

	isPermanent := terr.Code >= 500 && terr.Code < 600

	return newSMTPError(err, isPermanent, terr.Code)
}
