package models

import (
	"context"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type StorageBucket struct {
	ID        string `bun:"id,pk"`
	Name      string `bun:"name"`
	Slug      string `bun:"slug,unique"`
	Host      string `bun:"host"`
	KeyId     string `bun:"key_id"`
	SecretKey string `bun:"secret_key"`

	UserID string `bun:"user_id"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`
}

var _ bun.BeforeAppendModelHook = (*StorageBucket)(nil)

func (bucket *StorageBucket) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		bucket.ID = xid.New().String()
	}
	return nil
}
