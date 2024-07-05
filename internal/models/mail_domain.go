package models

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

const (
	DKIM_BITS_SIZE = 2048
)

type MailDomain struct {
	ID     string `bun:"id,pk"`
	Domain string `bun:"domain"`

	DKIMPrivateKey string `bun:"dkim_private_key"`
	DKIMPublicKey  string `bun:"dkim_public_key"`

	DNSVerified        bool            `bun:"dns_verified"`
	ExpectedDNSRecords json.RawMessage `bun:"expected_dns_records,type:jsonb,default:'[]'"`

	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	OrganizationID string        `bun:"organization_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ExpectedDNSRecordType represents the type of a DNS record, either MX or TXT
type ExpectedDNSRecordType string

const (
	ExpectedDNSRecordTypeMX  ExpectedDNSRecordType = "MX"
	ExpectedDNSRecordTypeTXT ExpectedDNSRecordType = "TXT"
)

// ExpectedDNSRecord represents a DNS record that is expected to be set for a domain
type ExpectedDNSRecord struct {
	Verified bool `json:"verified" example:"true" doc:"Record verification status"`

	Type  ExpectedDNSRecordType `json:"type" example:"MX" doc:"Record type"`
	Host  string                `json:"host" example:"example.com" doc:"Record host"`
	Value string                `json:"value" example:"mail.example.com" doc:"Record value"`
}

var _ bun.BeforeAppendModelHook = (*MailDomain)(nil)

func (d *MailDomain) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		// Generate a new ID
		d.ID = xid.New().String()

		// Set the creation time
		d.CreatedAt = time.Now()

		// Assign a new DKIM key pair
		if err := d.assignDKIMKeysPair(); err != nil {
			return err
		}

		// Assign expected DNS records
		if err := d.assignExpectedDNSRecords(); err != nil {
			return err
		}
	case *bun.UpdateQuery:
		// Update the modification time
		d.UpdatedAt = time.Now()
	}

	return nil
}

// assignDKIMKeysPair assigns a new DKIM key pair to the domain
func (d *MailDomain) assignDKIMKeysPair() error {
	// We generate a new RSA key pair
	key, err := rsa.GenerateKey(rand.Reader, DKIM_BITS_SIZE)
	if err != nil {
		return err
	}

	// Export the keys as base64 encoded strings
	d.DKIMPrivateKey = exportRsaPrivateKeyAsStr(key)
	d.DKIMPublicKey, err = exportRsaPublicKeyAsStr(&key.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

// exportRsaPrivateKeyAsStr exports an RSA private key as a base64 encoded string
func exportRsaPrivateKeyAsStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	return base64.StdEncoding.EncodeToString(privkeyBytes)
}

// exportRsaPublicKeyAsStr exports an RSA public key as a base64 encoded string
func exportRsaPublicKeyAsStr(key *rsa.PublicKey) (string, error) {
	privkeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(privkeyBytes), nil
}

// assignExpectedDNSRecords assigns the expected DNS records for the domain
func (d *MailDomain) assignExpectedDNSRecords() error {
	records := []ExpectedDNSRecord{
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  "mail._domainkey." + d.Domain,
			Value: fmt.Sprintf("v=DKIM1; k=rsa; p=%s", d.DKIMPublicKey),
		},
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  "_dmarc." + d.Domain,
			Value: "v=DMARC1; p=none;",
		},
		{
			Type:  ExpectedDNSRecordTypeTXT,
			Host:  d.Domain,
			Value: fmt.Sprintf("v=spf1 mx a:%s -all", os.Getenv("SMTP_DOMAIN")),
		},
		{
			Type:  ExpectedDNSRecordTypeMX,
			Host:  d.Domain,
			Value: os.Getenv("SMTP_DOMAIN"),
		},
	}

	data, err := json.Marshal(records)
	if err != nil {
		return err
	}

	d.ExpectedDNSRecords = data

	return nil
}

// setExpectedDNSRecords sets the expected DNS records for a domain
func (d *MailDomain) setExpectedDNSRecords(records []ExpectedDNSRecord) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}

	d.ExpectedDNSRecords = data

	return nil
}

// GetExpectedDNSRecords returns the expected DNS records for a domain
func (d *MailDomain) GetExpectedDNSRecords() []ExpectedDNSRecord {
	var records []ExpectedDNSRecord
	if err := json.Unmarshal(d.ExpectedDNSRecords, &records); err != nil {
		return []ExpectedDNSRecord{}
	}
	return records
}

// CheckDNS checks the DNS records for a domain
func (d *MailDomain) CheckDNS() error {
	// Retrieve the expected DNS records
	expectedRecords := d.GetExpectedDNSRecords()

	// Check the DNS records
	for i, record := range expectedRecords {
		log.Info("Checking DNS record", "type", record.Type, "host", record.Host, "value", record.Value)

		switch record.Type {
		case ExpectedDNSRecordTypeMX:
			log.Info("Looking up MX records", "domain", d.Domain)

			mxs, err := net.LookupMX(d.Domain)
			if err != nil {
				log.Error("Failed to lookup MX records", "err", err)
				expectedRecords[i].Verified = false
				continue
			}

			for _, mx := range mxs {
				log.Info("MX record", "host", mx.Host, "pref", mx.Pref, "value", record.Value)
				if mx.Host == record.Value {
					expectedRecords[i].Verified = true
					break
				}
			}
		case ExpectedDNSRecordTypeTXT:
			log.Info("Looking up TXT records", "domain", record.Host+"."+d.Domain)

			txt, err := net.LookupTXT(record.Host + "." + d.Domain)
			if err != nil {
				log.Error("Failed to lookup TXT records", "err", err)
				expectedRecords[i].Verified = false
				continue
			}

			for _, t := range txt {
				log.Info("TXT record", "t", t, "value", record.Value)
				if t == record.Value {
					expectedRecords[i].Verified = true
					break
				}
			}
		}
	}

	// Update the DNS verification status
	d.DNSVerified = true
	for _, record := range expectedRecords {
		if !record.Verified {
			d.DNSVerified = false
			break
		}
	}

	// Save the updated records
	if err := d.setExpectedDNSRecords(expectedRecords); err != nil {
		return err
	}

	return nil
}
