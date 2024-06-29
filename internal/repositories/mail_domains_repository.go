package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type MailDomainsRepository struct {
	*orm.Repository[models.MailDomain]
}

func NewMailDomainsRepository(db *orm.Database) *MailDomainsRepository {
	return &MailDomainsRepository{Repository: &orm.Repository[models.MailDomain]{Database: db}}
}

func (r *MailDomainsRepository) FindAllFromUser(ctx context.Context, userId string) ([]models.MailDomain, error) {
	var items []models.MailDomain = make([]models.MailDomain, 0)

	err := r.NewSelect().Model((*models.MailDomain)(nil)).Where("user_id = ?", userId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindVerifiedDomainWithUser finds a verified domain by its name and user ID.
func (r *MailDomainsRepository) FindVerifiedDomainWithUser(ctx context.Context, domainName, userID string) (*models.MailDomain, error) {
	var domain models.MailDomain
	err := r.NewSelect().
		Model(&domain).
		Where("domain = ? AND user_id = ? AND dns_verified = ?", domainName, userID, true).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}
