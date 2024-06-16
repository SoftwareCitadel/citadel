package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"
	"github.com/uptrace/bun"
)

type Application struct {
	ID             string          `bun:"id,pk"`
	Name           string          `bun:"name"`
	ReleaseCommand string          `bun:"release_command"`
	Slug           string          `bun:"slug,unique"`
	Env            json.RawMessage `bun:"env,type:jsonb,default:'[]'"`
	CpuConfig      string          `bun:"cpu_cfg"`
	RamConfig      string          `bun:"ram_cfg"`

	GitHubRepository     string `bun:"github_repository"`
	GitHubBranch         string `bun:"github_branch"`
	GitHubInstallationID int64  `bun:"github_installation_id,default:-1"`

	Certificates []*Certificate `bun:"rel:has-many,join:id=application_id"`

	UserID string `bun:"user_id"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var _ bun.BeforeAppendModelHook = (*Application)(nil)

func (app *Application) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		app.ID = xid.New().String()
		app.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		app.UpdatedAt = time.Now()
	}
	return nil
}

func (app *Application) GetEnv() map[string]string {
	var rawEnv []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(app.Env, &rawEnv); err != nil {
		return nil
	}

	env := make(map[string]string)
	for _, e := range rawEnv {
		env[e.Key] = e.Value
	}

	return env
}

func (app *Application) GetEnvVar(key string, defaultValues ...string) string {
	env := app.GetEnv()
	if val, ok := env[key]; ok {
		return val
	}

	if len(defaultValues) > 0 {
		return defaultValues[0]
	}

	return ""
}
