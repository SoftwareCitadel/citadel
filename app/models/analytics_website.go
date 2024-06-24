package models

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type AnalyticsWebsite struct {
	ID     string `bun:"id,pk"`
	Name   string `bun:"name,notnull"`
	Domain string `bun:"domain,notnull"`

	UserID string `bun:"user_id,notnull"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*AnalyticsWebsite)(nil)

func (m *AnalyticsWebsite) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = xid.New().String()
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
