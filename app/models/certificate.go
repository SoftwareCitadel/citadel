package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Certificate struct {
	ID            string            `bun:"id,pk"`
	Domain        string            `bun:"domain"`
	Status        CertificateStatus `bun:"status"`
	ValidDns      bool              `bun:"valid_dns"`
	DnsEntries    []DnsEntry        `bun:"dns_entries,type:jsonb"`
	ApplicationID string            `bun:"application_id"`
	Application   *Application      `bun:"rel:belongs-to,join:application_id=id"`
	CreatedAt     time.Time         `bun:"created_at"`
	UpdatedAt     time.Time         `bun:"updated_at"`
}

type CertificateStatus string

const (
	CertificateStatusPending  CertificateStatus = "pending"
	CertificateStatusVerified CertificateStatus = "verified"
)

type DnsEntry struct {
	Hostname string `bun:"hostname"`
	Type     string `bun:"type"`
	Value    string `bun:"value"`
}

var _ bun.BeforeAppendModelHook = (*Deployment)(nil)

func (certificate *Certificate) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		certificate.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		certificate.UpdatedAt = time.Now()
	}
	return nil
}
