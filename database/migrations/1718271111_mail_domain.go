package migrations

import (
	"citadel/app/models"
	"context"

	"github.com/uptrace/bun"
)

func mailDomainMigrationUp_1718271111(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.MailDomain)(nil)).Exec(ctx)
	return err
}

func mailDomainMigrationDown_1718271111(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.MailDomain)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(mailDomainMigrationUp_1718271111, mailDomainMigrationDown_1718271111)
}
