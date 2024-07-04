package listeners

import (
	mailBuilder "citadel/internal/mail_builder"
	"citadel/internal/models"
	smtpSender "citadel/internal/smtp_sender"
	"citadel/util"
	"net/mail"
	"os"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/charmbracelet/log"
)

type SMTPListener struct{}

func NewSMTPListener() *SMTPListener {
	return &SMTPListener{}
}

type OutboundEmailEvent struct {
	Msg    *mail.Message
	Domain *models.MailDomain
	From   string
	To     string
}

// OnOutboundEmail is called when an email is sent from the system.
func (smtpListener *SMTPListener) OnOutboundEmail(msg *message.Message) ([]*message.Message, error) {
	// Parse the event
	var event OutboundEmailEvent
	if err := util.DecodeJSON(msg.Payload, &event); err != nil {
		return nil, err
	}

	// Build the email to send
	builder := mailBuilder.New(event.Msg)
	outputMsg, err := builder.Build()
	if err != nil {
		log.Error("failed to build the email:", "error", err)
		return nil, err
	}

	// Sign the email with DKIM
	outputMsg, err = builder.SignWithDKIM(outputMsg, event.Domain.Domain, event.Domain.DKIMPrivateKey)
	if err != nil {
		log.Error("failed to sign the email with DKIM:", "error", err)
		return nil, err
	}

	// Send the email
	sender := smtpSender.New(os.Getenv("SMTP_DOMAIN"))
	sender.Send(event.From, event.To, outputMsg)

	return nil, nil
}
