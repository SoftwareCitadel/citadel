package models

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type MailDomain struct {
	ID         string     `bun:"id,pk"`
	Domain     string     `bun:"domain"`
	IsVerified bool       `bun:"is_verified"`
	DnsEntries []DnsEntry `bun:"dns_entries,type:jsonb"`

	User   *User  `bun:"rel:belongs-to,join:user_id=id"`
	UserID string `bun:"user_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var _ bun.BeforeAppendModelHook = (*MailDomain)(nil)

func (domain *MailDomain) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		domain.ID = xid.New().String()
		domain.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		domain.UpdatedAt = time.Now()
	}
	return nil
}
