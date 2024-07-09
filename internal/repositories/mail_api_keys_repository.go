package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type MailApiKeysRepository struct {
	*orm.Repository[models.MailApiKey]
}

func NewMailApiKeysRepository(db *orm.Database) *MailApiKeysRepository {
	return &MailApiKeysRepository{Repository: &orm.Repository[models.MailApiKey]{
		Database: db,
	}}
}

func (r MailApiKeysRepository) FindAllFromOrgWithRelatedDomain(ctx context.Context, orgId string) ([]models.MailApiKey, error) {
	var items []models.MailApiKey = make([]models.MailApiKey, 0)

	err := r.NewSelect().Model((*models.MailApiKey)(nil)).
		Relation("MailDomain").
		Where("mail_api_key.organization_id = ?", orgId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
