package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/Squwid/go-randomizer"
	"github.com/caesar-rocks/orm"
	"github.com/gosimple/slug"
)

type OrganizationsRepository struct {
	*orm.Repository[models.Organization]
}

func NewOrganizationsRepository(db *orm.Database) *OrganizationsRepository {
	return &OrganizationsRepository{Repository: &orm.Repository[models.Organization]{
		Database: db,
	}}
}

func (r OrganizationsRepository) Create(ctx context.Context, org *models.Organization) error {
	slug := slug.Make(org.Name)

	for {
		if _, err := r.FindOneBy(ctx, "slug", slug); err != nil {
			break
		}

		slug = slug + "-" + randomizer.Noun()
	}

	org.Slug = slug

	if _, err := r.NewInsert().Model(org).Exec(ctx); err != nil {
		return err
	}

	return nil
}
