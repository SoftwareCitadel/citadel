package migrations

import (
	"citadel/internal/models"
	"context"

	"github.com/uptrace/bun"
)

func organizationMemberMigrationUp_1719494629(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.OrganizationMember)(nil)).Exec(ctx)
	return err
}

func organizationMemberMigrationDown_1719494629(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDropTable().Model((*models.OrganizationMember)(nil)).Exec(ctx)
	return err
}

func init() {
	Migrations.MustRegister(organizationMemberMigrationUp_1719494629, organizationMemberMigrationDown_1719494629)
}
