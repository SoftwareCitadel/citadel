package types

import (
	"github.com/google/go-github/v62/github"
)

type GithubOwner struct {
	*github.User
	InstallationID int64
}
