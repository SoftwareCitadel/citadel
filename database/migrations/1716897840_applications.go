package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func applicationsMigrationUp_1716897840(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.Application)(nil)).Exec(ctx)
	return err
}

func applicationsMigrationDown_1716897840(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.Application)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(applicationsMigrationUp_1716897840, applicationsMigrationDown_1716897840)
}
