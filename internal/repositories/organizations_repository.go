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
	return &OrganizationsRepository{
		Repository: &orm.Repository[models.Organization]{Database: db},
	}
}

func (r *OrganizationsRepository) Create(ctx context.Context, org *models.Organization) error {
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

// FindFirstOwnedByUser finds the first organization owned by the user.
// This finds the first organization where it has a member with the role "owner".
func (r *OrganizationsRepository) FindFirstOwnedByUser(ctx context.Context, userId string) (*models.Organization, error) {
	var org models.Organization
	err := r.NewSelect().
		Model(&org).
		Join("JOIN organization_members om ON om.organization_id = organization.id").
		Where("om.user_id = ? AND om.role = ?", userId, models.OrganizationMemberRoleOwner).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAllWhereUserIsMember finds all organizations where the user is a member.
func (r *OrganizationsRepository) FindAllWhereUserIsMember(ctx context.Context, userId string) ([]models.Organization, error) {
	var orgs []models.Organization

	err := r.NewSelect().
		Model(&orgs).
		Relation("OrganizationMembers").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
