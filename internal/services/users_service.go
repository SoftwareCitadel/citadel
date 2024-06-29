package services

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	"citadel/util"
	"context"

	"github.com/caesar-rocks/events"
)

type UsersService struct {
	usersRepo      *repositories.UsersRepository
	orgsRepo       *repositories.OrganizationsRepository
	orgMembersRepo *repositories.OrganizationMembersRepository
	emitter        *events.EventsEmitter
}

func NewUsersService(usersRepo *repositories.UsersRepository, orgsRepo *repositories.OrganizationsRepository, orgMembersRepo *repositories.OrganizationMembersRepository, emitter *events.EventsEmitter) *UsersService {
	return &UsersService{usersRepo, orgsRepo, orgMembersRepo, emitter}
}

func (s *UsersService) CreateAndEmitEvent(ctx context.Context, user *models.User) error {
	if err := s.usersRepo.Create(ctx, user); err != nil {
		return err
	}

	org := &models.Organization{Name: user.FullName}
	if err := s.orgsRepo.Create(ctx, org); err != nil {
		return err
	}

	member := &models.OrganizationMember{UserID: user.ID, OrganizationID: org.ID, Role: models.OrganizationMemberRoleOwner}
	if err := s.orgMembersRepo.Create(ctx, member); err != nil {
		return err
	}

	bytes, err := util.EncodeJSON(user)
	if err != nil {
		return err
	}
	s.emitter.Emit("users.created", bytes)

	return nil
}
