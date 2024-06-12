package migrations

import (
	"citadel/app/models"
	"context"

	"github.com/uptrace/bun"
)

func storage_bucketsMigrationUp_1718139540(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.StorageBucket)(nil)).Exec(ctx)
	return err
}

func storage_bucketsMigrationDown_1718139540(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.StorageBucket)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(storage_bucketsMigrationUp_1718139540, storage_bucketsMigrationDown_1718139540)
}
