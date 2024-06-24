package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Refer to Bun's documentation: https://bun.uptrace.dev/
type WebsiteVisit struct {
	ID int64 `bun:"id,pk,autoincrement"`

	// Add your own fields here...

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

var _ bun.BeforeAppendModelHook = (*WebsiteVisit)(nil)

func (m *WebsiteVisit) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
