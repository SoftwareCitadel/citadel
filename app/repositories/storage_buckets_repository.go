package repositories

import (
	"citadel/app/models"
	"context"

	"github.com/Squwid/go-randomizer"
	"github.com/caesar-rocks/orm"
	"github.com/gosimple/slug"
)

type StorageBucketsRepository struct {
	*orm.Repository[models.StorageBucket]
}

func NewStorageBucketsRepository(db *orm.Database) *StorageBucketsRepository {
	return &StorageBucketsRepository{Repository: &orm.Repository[models.StorageBucket]{
		Database: db,
	}}
}

func (r StorageBucketsRepository) FindAllFromUser(ctx context.Context, userId string) ([]models.StorageBucket, error) {
	var items []models.StorageBucket = make([]models.StorageBucket, 0)

	err := r.NewSelect().Model((*models.StorageBucket)(nil)).Where("user_id = ?", userId).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r StorageBucketsRepository) Create(ctx context.Context, storageBucket *models.StorageBucket) error {
	slug := slug.Make(storageBucket.Name)

	for {
		_, err := r.FindOneBy(ctx, "slug", slug)
		if err != nil {
			break
		}

		slug = slug + "-" + randomizer.Noun()
	}
	storageBucket.Slug = slug

	_, err := r.NewInsert().Model(storageBucket).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
