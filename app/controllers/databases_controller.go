package controllers

import (
	"citadel/app/drivers"
	"citadel/app/models"
	"citadel/app/repositories"
	"citadel/views/pages"
	"os"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type DatabasesController struct {
	dbRepo *repositories.DatabasesRepository
	driver drivers.Driver
}

func NewDatabasesController(dbRepo *repositories.DatabasesRepository, driver drivers.Driver) *DatabasesController {
	return &DatabasesController{dbRepo, driver}
}

func (c *DatabasesController) Index(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	dbs, err := c.dbRepo.FindAllFromUser(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	return ctx.Render(pages.DatabasesPage(dbs))
}

type StoreDatabaseValidator struct {
	Name     string      `form:"name" validate:"required"`
	DBMS     models.DBMS `form:"dbms" validate:"required,oneof=mysql postgres redis"`
	Username string      `form:"username"`
	Password string      `form:"password"`
}

func (c *DatabasesController) Store(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreDatabaseValidator](ctx)
	if !ok {
		return ctx.Redirect("/databases")
	}

	db := &models.Database{
		Name:     data.Name,
		DBMS:     data.DBMS,
		Username: data.Username,
		Password: data.Password,
		UserID:   user.ID,
		Host:     os.Getenv("DB_HOST"),
	}
	if err := c.dbRepo.Create(ctx.Context(), db); err != nil {
		return err
	}

	if err := c.driver.CreateDatabase(*db); err != nil {
		return err
	}

	return ctx.Redirect("/databases")
}
