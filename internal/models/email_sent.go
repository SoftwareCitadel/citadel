package models

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type Email struct {
	ID             string        `bun:"id,pk"`
	Subject        string        `bun:"subject"`
	Slug           string        `bun:"slug,unique"`
	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	OrganizationID string        `bun:"rel:organization_id"`
	Sender         string        `bun:"sender"`
	Recipient      string        `bun:"recipient"`
	Body           string        `bun:"body"`
	Inc_attach     bool          `bun:"incl_attach"`
	Send_time      time.Time     `bun:"send_time"`
	CreatedAt      time.Time     `bun:"created_at"`
	UpdatedAt      time.Time     `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*Deployment)(nil)

func (e *Email) BeforeAppend(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		e.ID = xid.New().String()
		e.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		e.UpdatedAt = time.Now()
	}
	return nil
}
