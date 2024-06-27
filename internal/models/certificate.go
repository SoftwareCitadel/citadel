package models

import (
	"context"
	"time"

	"github.com/rs/xid"
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

func (cert *Certificate) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		cert.ID = xid.New().String()
		cert.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		cert.UpdatedAt = time.Now()
	}
	return nil
}
