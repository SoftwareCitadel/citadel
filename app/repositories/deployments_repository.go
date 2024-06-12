package repositories

import (
	"citadel/app/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type DeploymentsRepository struct {
	*orm.Repository[models.Deployment]
}

func NewDeploymentsRepository(db *orm.Database) *DeploymentsRepository {
	return &DeploymentsRepository{Repository: &orm.Repository[models.Deployment]{
		Database: db,
	}}
}

func (r *DeploymentsRepository) FindAllFromApplication(ctx context.Context, appId string) ([]models.Deployment, error) {
	var items []models.Deployment = make([]models.Deployment, 0)

	err := r.NewSelect().Model((*models.Deployment)(nil)).Where("application_id = ?", appId).Order("created_at DESC").Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *DeploymentsRepository) FindOneByIdWithRelatedAppAndCerts(ctx context.Context, id string) (*models.Deployment, error) {
	var item *models.Deployment = new(models.Deployment)

	// Use table aliasing to avoid ambiguity
	err := r.NewSelect().
		Model(item).
		Where("deployment.id = ?", id).
		Relation("Application").
		Relation("Application.Certificates").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return item, nil
}
