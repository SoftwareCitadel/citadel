package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/caesar-rocks/orm"
)

type OrganizationMembersRepository struct {
	*orm.Repository[models.OrganizationMember]
}

func NewOrganizationMembersRepository(db *orm.Database) *OrganizationMembersRepository {
	return &OrganizationMembersRepository{Repository: &orm.Repository[models.OrganizationMember]{
		Database: db,
	}}
}

func (r *OrganizationMembersRepository) FindAllFromOrganizationWithUser(ctx context.Context, orgID string) ([]*models.OrganizationMember, error) {
	var members []*models.OrganizationMember
	if err := r.
		NewSelect().
		Model((*models.OrganizationMember)(nil)).
		Relation("User").
		Where("organization_id = ?", orgID).
		Scan(ctx, &members); err != nil {
		return nil, err
	}

	return members, nil
}
