package migrations

import (
	"citadel/app/models"
	"context"

	"github.com/uptrace/bun"
)

func databasesMigrationUp_1718442825(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.Database)(nil)).Exec(ctx)
	return err
}

func databasesMigrationDown_1718442825(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.Database)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(databasesMigrationUp_1718442825, databasesMigrationDown_1718442825)
}
