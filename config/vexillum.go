package config

import (
	"github.com/caesar-rocks/vexillum"
)

func ProvideVexillum(env *EnvironmentVariables) *vexillum.Vexillum {
	v := vexillum.New()

	if env.GITHUB_OAUTH_KEY != "" && env.GITHUB_OAUTH_SECRET != "" && env.GITHUB_OAUTH_CALLBACK_URL != "" {
		v.Activate("github_oauth")
	}

	if env.GITHUB_APP_ID != "" && env.GITHUB_APP_PRIVATE_KEY != "" && env.GITHUB_APP_WEBHOOK_SECRET != "" && env.GITHUB_APP_PRIVATE_KEY_PATH != "" {
		v.Activate("github_deployments")
	}

	if env.STRIPE_PUBLIC_KEY != "" && env.STRIPE_SECRET_KEY != "" {
		v.Activate("billing")
	}

	return v
}
