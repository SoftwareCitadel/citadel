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

func (r *MailDomainsRepository) FindAllFromOrg(ctx context.Context, orgId string) ([]models.MailDomain, error) {
	var items []models.MailDomain = make([]models.MailDomain, 0)

	err := r.NewSelect().Model((*models.MailDomain)(nil)).Where("organization_id = ?", orgId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindVerifiedDomainWithOrg finds a verified domain by its name and org ID.
func (r *MailDomainsRepository) FindVerifiedDomainWithOrg(ctx context.Context, domainName, orgID string) (*models.MailDomain, error) {
	var domain models.MailDomain
	err := r.NewSelect().
		Model(&domain).
		Where("domain = ?", domainName).
		Where("organization_id = ?", orgID).
		Where("dns_verified = true").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}
