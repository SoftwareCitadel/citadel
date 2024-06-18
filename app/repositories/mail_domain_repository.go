package repositories

import (
	"citadel/app/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type MailDomainsRepository struct {
	*orm.Repository[models.MailDomain]
}

func NewMailDomainsRepository(db *orm.Database) *MailDomainsRepository {
	return &MailDomainsRepository{Repository: &orm.Repository[models.MailDomain]{Database: db}}
}

func (r MailDomainsRepository) FindAllFromUser(ctx context.Context, userId string) ([]models.MailDomain, error) {
	var items []models.MailDomain = make([]models.MailDomain, 0)

	err := r.NewSelect().Model((*models.MailDomain)(nil)).Where("user_id = ?", userId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
