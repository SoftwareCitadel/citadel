package database

import (
	"citadel/database/migrations"
	"embed"

	"github.com/uptrace/bun/migrate"
)

//go:embed migrations
var FS embed.FS

func GetMigrations() *migrate.Migrations {
	if err := migrations.Migrations.Discover(FS); err != nil {
		panic(err)
	}

	return migrations.Migrations
}
