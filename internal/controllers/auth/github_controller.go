package authControllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	"citadel/internal/services"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type GithubController struct {
	auth    *auth.Auth
	repo    *repositories.UsersRepository
	service *services.UsersService
}

func NewGithubController(auth *auth.Auth, repo *repositories.UsersRepository, service *services.UsersService) *GithubController {
	return &GithubController{auth, repo, service}
}

func (c *GithubController) Redirect(ctx *caesar.Context) error {
	return c.auth.Social.Use("github").Redirect(ctx)
}

func (c *GithubController) Callback(ctx *caesar.Context) error {
	oauthUser, err := c.auth.Social.Use("github").Callback(ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	user, err := c.repo.FindOneBy(ctx.Context(), "github_user_id", oauthUser.UserID)
	if err != nil {
		user = &models.User{Email: oauthUser.Email, FullName: oauthUser.Name, GitHubUserID: oauthUser.UserID}
		if err := c.service.CreateAndEmitEvent(ctx.Context(), user); err != nil {
			return caesar.NewError(400)
		}
	}

	if err := c.auth.Authenticate(ctx, *user); err != nil {
		return caesar.NewError(400)
	}

	return ctx.Redirect("/apps")
}
