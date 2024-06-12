package migrations

import (
	"citadel/app/models"
	"context"

	"github.com/uptrace/bun"
)

func deploymentMigrationUp_1717507533(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.Deployment)(nil)).Exec(ctx)
	return err
}

func deploymentMigrationDown_1717507533(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropColumn().Model((*models.Deployment)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(deploymentMigrationUp_1717507533, deploymentMigrationDown_1717507533)
}
