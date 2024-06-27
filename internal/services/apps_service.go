package services

import (
	"citadel/internal/models"
	"citadel/internal/repositories"

	"github.com/caesar-rocks/auth"
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

func (s *AppsService) GetAppOwnedByCurrentUser(ctx *caesar.Context) (*models.Application, error) {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return nil, err
	}

	app, err := s.appsRepo.FindOneBy(ctx.Context(), "slug", ctx.PathValue("slug"))
	if err != nil {
		return nil, err
	}

	if app.UserID != user.ID {
		return nil, caesar.NewError(403)
	}

	return app, nil
}
