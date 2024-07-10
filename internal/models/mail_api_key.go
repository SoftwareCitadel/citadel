package models

import (
	"citadel/util"
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type MailApiKey struct {
	ID    string `bun:"id,pk"`
	Name  string `bun:"name"`
	Value string `bun:"value"`

	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	OrganizationID string        `bun:"organization_id"`

	MailDomain   *MailDomain `bun:"rel:belongs-to,join:mail_domain_id=id"`
	MailDomainID string      `bun:"mail_domain_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	IsNewlyCreated bool
}

var _ bun.BeforeAppendModelHook = (*MailApiKey)(nil)

func (domain *MailApiKey) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		domain.ID = xid.New().String()
		domain.CreatedAt = time.Now()
		value, err := util.GenerateSecretKey()
		if err != nil {
			return err
		}
		domain.Value = value
	case *bun.UpdateQuery:
		domain.UpdatedAt = time.Now()
	}
	return nil
}
