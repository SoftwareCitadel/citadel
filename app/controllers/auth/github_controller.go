package authControllers

import (
	"citadel/app/models"
	"citadel/app/repositories"
	"citadel/app/services"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/rs/xid"
)

type GithubController struct {
	auth    *auth.Auth
	repo    *repositories.UsersRepository
	service *services.UsersService
}

func NewGithubController(auth *auth.Auth, repo *repositories.UsersRepository, service *services.UsersService) *GithubController {
	return &GithubController{auth, repo, service}
}

func (c *GithubController) Redirect(ctx *caesar.CaesarCtx) error {
	return c.auth.Social.Use("github").Redirect(ctx)
}

func (c *GithubController) Callback(ctx *caesar.CaesarCtx) error {
	oauthUser, err := c.auth.Social.Use("github").Callback(ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	user, err := c.repo.FindOneBy(ctx.Context(), "github_user_id", oauthUser.UserID)
	if err != nil {
		user = &models.User{ID: xid.New().String(), Email: oauthUser.Email, FullName: oauthUser.Name, GitHubUserID: oauthUser.UserID}
		if err := c.service.CreateAndEmitEvent(ctx.Context(), user); err != nil {
			return caesar.NewError(400)
		}
	}

	if err := c.auth.Authenticate(ctx, *user); err != nil {
		return caesar.NewError(400)
	}

	return ctx.Redirect("/applications")
}
