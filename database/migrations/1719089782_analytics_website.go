package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func analyticsWebsiteMigrationUp_1719089782(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.AnalyticsWebsite)(nil)).Exec(ctx)
	return err
}

func analyticsWebsiteMigrationDown_1719089782(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.AnalyticsWebsite)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(analyticsWebsiteMigrationUp_1719089782, analyticsWebsiteMigrationDown_1719089782)
}
