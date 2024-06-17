package models

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type Database struct {
	ID       string `bun:"id,pk"`
	Name     string `bun:"name"`
	Slug     string `bun:"slug,unique"`
	DBMS     DBMS   `bun:"dbms"`
	Host     string `bun:"host"`
	Username string `bun:"username"`
	Password string `bun:"password"`

	UserID string `bun:"user_id"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DBMS string

const (
	Postgres DBMS = "postgres"
	MySQL    DBMS = "mysql"
	Redis    DBMS = "redis"
)

var _ bun.BeforeAppendModelHook = (*Database)(nil)

func (db *Database) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		db.ID = xid.New().String()
		db.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		db.UpdatedAt = time.Now()
	}
	return nil
}

func (db *Database) GetURI() string {
	switch db.DBMS {
	case MySQL:
		return "mysql://" + db.Username + ":" + db.Password + "@" + db.Host + "/" + db.Name
	case Postgres:
		return "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + "/" + db.Name
	case Redis:
		return "redis://default:" + db.Password + "@" + db.Host
	default:
		return ""
	}
}
