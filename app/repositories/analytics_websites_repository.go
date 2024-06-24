package repositories

import (
	"citadel/app/models"

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
