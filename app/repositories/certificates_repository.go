package repositories

import (
	"citadel/app/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type CertificatesRepository struct {
	*orm.Repository[models.Certificate]
}

func NewCertificatesRepository(db *orm.Database) *CertificatesRepository {
	return &CertificatesRepository{Repository: &orm.Repository[models.Certificate]{
		Database: db,
	}}
}

func (r *CertificatesRepository) FindAllFromApp(ctx context.Context, appId string) ([]models.Certificate, error) {
	var items []models.Certificate

	err := r.NewSelect().Model((*models.Certificate)(nil)).Where(
		"application_id = ?", appId,
	).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
