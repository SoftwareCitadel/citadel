package services

import (
	"citadel/app/models"
	"citadel/app/repositories"
	"citadel/util"
	"context"

	"github.com/caesar-rocks/events"
)

type UsersService struct {
	repo    *repositories.UsersRepository
	emitter *events.EventsEmitter
}

func NewUsersService(repo *repositories.UsersRepository, emitter *events.EventsEmitter) *UsersService {
	return &UsersService{repo, emitter}
}

func (s *UsersService) CreateAndEmitEvent(ctx context.Context, user *models.User) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	bytes, err := util.EncodeJSON(user)
	if err != nil {
		return err
	}
	s.emitter.Emit("users.created", bytes)

	return nil
}
