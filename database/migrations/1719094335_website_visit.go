package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func websiteVisitMigrationUp_1719094335(ctx context.Context, db *bun.DB) error {
	return nil
}

func websiteVisitMigrationDown_1719094335(ctx context.Context, db *bun.DB) error {
	return nil
}

func init() {
	Migrations.MustRegister(websiteVisitMigrationUp_1719094335, websiteVisitMigrationDown_1719094335)
}
