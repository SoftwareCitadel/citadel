package repositories

import (
	"citadel/app/models"
	"context"

	"github.com/Squwid/go-randomizer"
	"github.com/caesar-rocks/orm"
	"github.com/gosimple/slug"
)

type ApplicationsRepository struct {
	*orm.Repository[models.Application]
}

func NewApplicationsRepository(db *orm.Database) *ApplicationsRepository {
	return &ApplicationsRepository{Repository: &orm.Repository[models.Application]{
		Database: db,
	}}
}

func (r ApplicationsRepository) Create(ctx context.Context, app *models.Application) error {
	slug := slug.Make(app.Name)

	for {
		_, err := r.FindOneBy(ctx, "slug", slug)
		if err != nil {
			break
		}

		slug = slug + "-" + randomizer.Noun()
	}
	app.Slug = slug

	_, err := r.NewInsert().Model(app).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r ApplicationsRepository) FindAllFromUser(ctx context.Context, userId string) ([]models.Application, error) {
	var items []models.Application = make([]models.Application, 0)

	err := r.NewSelect().Model((*models.Application)(nil)).Where("user_id = ?", userId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
