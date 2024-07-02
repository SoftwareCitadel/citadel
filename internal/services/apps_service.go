package services

import (
	"citadel/internal/models"
	"citadel/internal/repositories"

	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/events"
)

type AppsService struct {
	usersRepo *repositories.UsersRepository
	appsRepo  *repositories.ApplicationsRepository
	emitter   *events.EventsEmitter
}

func NewAppsService(usersRepo *repositories.UsersRepository, appsRepo *repositories.ApplicationsRepository, emitter *events.EventsEmitter) *AppsService {
	return &AppsService{usersRepo, appsRepo, emitter}
}

func (s *AppsService) GetAppOwnedByCurrentOrg(ctx *caesar.Context) (*models.Application, error) {
	app, err := s.appsRepo.FindOneBy(
		ctx.Context(),
		"slug", ctx.PathValue("slug"),
		"organization_id", ctx.PathValue("orgId"),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}
