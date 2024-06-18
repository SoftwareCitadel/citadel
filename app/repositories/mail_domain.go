package repositories

import (
	"citadel/app/models"

	"github.com/caesar-rocks/orm"
)

type MailDomainsRepository struct {
	*orm.Repository[models.MailDomain]
}

func NewMailDomainsRepository(db *orm.Database) *MailDomainsRepository {
	return &MailDomainsRepository{Repository: &orm.Repository[models.MailDomain]{Database: db}}
}
