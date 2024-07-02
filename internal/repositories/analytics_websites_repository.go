package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type AnalyticsWebsitesRepository struct {
	*orm.Repository[models.AnalyticsWebsite]
}

func NewAnalyticsWebsitesRepository(db *orm.Database) *AnalyticsWebsitesRepository {
	return &AnalyticsWebsitesRepository{Repository: &orm.Repository[models.AnalyticsWebsite]{
		Database: db,
	}}
}

func (r AnalyticsWebsitesRepository) FindAllFromOrg(ctx context.Context, orgId string) ([]models.AnalyticsWebsite, error) {
	var items []models.AnalyticsWebsite = make([]models.AnalyticsWebsite, 0)

	err := r.NewSelect().Model((*models.AnalyticsWebsite)(nil)).Where("organization_id = ?", orgId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
