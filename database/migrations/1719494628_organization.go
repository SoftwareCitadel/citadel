package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func organizationMigrationUp_1719494655(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.Organization)(nil)).Exec(ctx)
	return err
}

func organizationMigrationDown_1719494655(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.Organization)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(organizationMigrationUp_1719494655, organizationMigrationDown_1719494655)
}
