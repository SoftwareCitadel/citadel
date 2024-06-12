package config

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
)

func ProvideGithub(env *EnvironmentVariables) *github.Client {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 1, 1, env.GITHUB_APP_PRIVATE_KEY_PATH)
	if err != nil {
		panic(err)
	}

	client := github.NewClient(&http.Client{Transport: itr})

	return client
}
