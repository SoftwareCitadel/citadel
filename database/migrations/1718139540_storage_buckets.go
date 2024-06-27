package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func storageBucketsMigrationUp_1718139540(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.StorageBucket)(nil)).Exec(ctx)
	return err
}

func storageBucketsMigrationDown_1718139540(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.StorageBucket)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(storageBucketsMigrationUp_1718139540, storageBucketsMigrationDown_1718139540)
}
