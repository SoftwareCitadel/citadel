package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	"citadel/internal/types"
	appsPages "citadel/views/concerns/apps/pages"
	"fmt"
	"net/http"
	"os"

	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/google/go-github/v62/github"
)

type GithubController struct {
	repo *repositories.UsersRepository
}

func NewGithubController(repo *repositories.UsersRepository) *GithubController {
	return &GithubController{repo}
}

func (c *GithubController) HandleWebhook(ctx *caesar.Context) error {
	payload, err := github.ValidatePayload(ctx.Request, []byte(os.Getenv("GITHUB_APP_WEBHOOK_SECRET")))
	if err != nil {
		return err
	}

	event, err := github.ParseWebHook(github.WebHookType(ctx.Request), payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *github.PushEvent:
		return c.processPushEvent(ctx, event)
	case *github.InstallationEvent:
		return c.processInstallationEvent(ctx, event)
	}

	return nil
}

func (c *GithubController) processPushEvent(ctx *caesar.Context, event *github.PushEvent) error {
	ghRepo := event.Repo.Owner.GetLogin() + "/" + event.Repo.GetName()

	// TODO: Retrrigger the deployment for the given repository.
	fmt.Println("Push event received for repository:", ghRepo)

	// TODO: Check if the branches match
	branch := strings.Replace(event.GetRef(), "refs/heads/", "", 1)
	fmt.Println("Branch:", branch)

	// TODO: Mark deploying (GitHub check on the commit)

	// TODO: Download the commit, and put it in the S3 bucket

	// TODO: Create the deployment in the db

	// TODO: Ignite the builder

	return nil
}

func (c *GithubController) processInstallationEvent(ctx *caesar.Context, event *github.InstallationEvent) error {
	user, _ := c.repo.FindOneBy(ctx.Context(), "github_user_id", event.Sender.ID)
	if user == nil {
		return nil
	}

	if event.GetAction() == "created" {
		user.AddGitHubInstallationId(event.Installation.GetID())
	}

	if event.GetAction() == "deleted" {
		user.RemoveGitHubInstallationId(event.Installation.GetID())
	}

	if err := c.repo.UpdateOneWhere(ctx.Context(), user, "id", user.ID); err != nil {
		return caesar.NewError(500)
	}

	return nil
}

func (c *GithubController) ListRepositories(ctx *caesar.Context) error {
	reposMap := make(map[string][]*github.Repository)
	owners := []types.GithubOwner{}

	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	githubAppId, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		return err
	}

	for _, installationID := range user.GetGitHubInstallationIDs() {
		ownerAdded := false

		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, githubAppId, installationID, os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
		if err != nil {
			return err
		}

		client := github.NewClient(&http.Client{Transport: itr})

		opt := &github.ListOptions{}
		for {
			repos, resp, err := client.Apps.ListRepos(ctx.Context(), opt)
			if err != nil {
				return err
			}

			for _, repo := range repos.Repositories {
				if !ownerAdded {
					owners = append(owners, types.GithubOwner{User: repo.GetOwner(), InstallationID: installationID})
					ownerAdded = true
				}

				orgName := repo.GetOwner().GetLogin()
				reposMap[orgName] = append(reposMap[orgName], repo)
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(reposMap)
	}

	return ctx.Render(appsPages.ConnectGitHubRepositoryList(owners, reposMap))
}
