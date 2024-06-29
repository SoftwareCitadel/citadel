package repositories

import (
	"citadel/internal/models"

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
