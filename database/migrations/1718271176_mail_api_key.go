package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func mailApiKeyMigrationUp_1718271176(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.MailApiKey)(nil)).Exec(ctx)
	return err
}

func mailApiKeyMigrationDown_1718271176(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.MailApiKey)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(mailApiKeyMigrationUp_1718271176, mailApiKeyMigrationDown_1718271176)
}
