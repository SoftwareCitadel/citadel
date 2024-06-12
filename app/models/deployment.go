package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Deployment struct {
	ID            string           `bun:"id,pk"`
	Origin        DeploymentOrigin `bun:"origin"`
	Status        DeploymentStatus `bun:"status"`
	ApplicationID string           `bun:"application_id"`
	Application   *Application     `bun:"rel:belongs-to,join:application_id=id"`
	CreatedAt     time.Time        `bun:"created_at"`
	UpdatedAt     time.Time        `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*Deployment)(nil)

func (deployment *Deployment) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		deployment.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		deployment.UpdatedAt = time.Now()
	}
	return nil
}

type DeploymentOrigin string

const (
	DeploymentOriginCli    DeploymentOrigin = "CLI"
	DeploymentOriginGithub DeploymentOrigin = "GitHub"
)

func (origin DeploymentOrigin) String() string {
	return string(origin)
}

type DeploymentStatus string

const (
	DeploymentStatusBuilding     DeploymentStatus = "Building"
	DeploymentStatusBuildFailed  DeploymentStatus = "Build Failed"
	DeploymentStatusDeploying    DeploymentStatus = "Deploying"
	DeploymentStatusDeployFailed DeploymentStatus = "Deployment Failed"
	DeploymentStatusSuccess      DeploymentStatus = "Success"
)

func (status DeploymentStatus) String() string {
	return string(status)
}
