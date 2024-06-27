package repositories

import (
	"citadel/internal/models"

	"github.com/caesar-rocks/orm"
)

type UsersRepository struct {
	*orm.Repository[models.User]
}

func NewUsersRepository(db *orm.Database) *UsersRepository {
	return &UsersRepository{Repository: &orm.Repository[models.User]{Database: db}}
}
