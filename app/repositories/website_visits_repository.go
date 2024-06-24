package repositories

import (
	"citadel/app/models"

	"github.com/caesar-rocks/orm"
)

type WebsiteVisitsRepository struct {
	*orm.Repository[models.WebsiteVisit]
}

func NewWebsiteVisitsRepository(db *orm.Database) *WebsiteVisitsRepository {
	return &WebsiteVisitsRepository{Repository: &orm.Repository[models.WebsiteVisit]{
		Database: db,
	}}
}
