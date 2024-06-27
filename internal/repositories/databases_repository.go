package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/Squwid/go-randomizer"
	"github.com/caesar-rocks/orm"
	"github.com/gosimple/slug"
)

type DatabasesRepository struct {
	*orm.Repository[models.Database]
}

func NewDatabasesRepository(db *orm.Database) *DatabasesRepository {
	return &DatabasesRepository{Repository: &orm.Repository[models.Database]{Database: db}}
}

func (r DatabasesRepository) Create(ctx context.Context, db *models.Database) error {
	slug := slug.Make(db.Name)

	for {
		_, err := r.FindOneBy(ctx, "slug", slug)
		if err != nil {
			break
		}

		slug = slug + "-" + randomizer.Noun()
	}
	db.Slug = slug

	_, err := r.NewInsert().Model(db).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r DatabasesRepository) FindAllFromUser(ctx context.Context, userId string) ([]models.Database, error) {
	var items []models.Database = make([]models.Database, 0)

	err := r.NewSelect().Model((*models.Database)(nil)).Where("user_id = ?", userId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
