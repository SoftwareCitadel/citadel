package models

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type OrganizationMember struct {
	ID   string                 `bun:"id,pk"`
	Role OrganizationMemberRole `bun:"role,notnull"`

	OrganizationID string        `bun:"organization_id"`
	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id"`

	User   *User  `bun:"rel:belongs-to,join:user_id=id"`
	UserID string `bun:"user_id"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type OrganizationMemberRole string

const (
	OrganizationMemberRoleOwner OrganizationMemberRole = "owner"
	OrganizationMemberRoleAdmin OrganizationMemberRole = "member"
)

var _ bun.BeforeAppendModelHook = (*OrganizationMember)(nil)

func (m *OrganizationMember) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = xid.New().String()
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
