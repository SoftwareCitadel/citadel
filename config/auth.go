package config

import (
	"citadel/internal/repositories"
	"context"
	"time"

	"github.com/caesar-rocks/auth"
)

func ProvideAuth(env *EnvironmentVariables, usersRepository *repositories.UsersRepository) *auth.Auth {
	return auth.NewAuth(&auth.AuthCfg{
		Key:           env.APP_KEY,
		JWTSigningKey: []byte(env.APP_KEY),
		MaxAge:        time.Hour * 24 * 30,
		JWTExpiration: time.Hour * 24 * 30,
		UserProvider: func(ctx context.Context, userID any) (any, error) {
			return usersRepository.FindOneBy(ctx, "id", userID)
		},
		RedirectTo: "/auth/sign_in",
		SocialProviders: &map[string]auth.SocialAuthProvider{
			"github": {
				Key:         env.GITHUB_OAUTH_KEY,
				Secret:      env.GITHUB_OAUTH_SECRET,
				CallbackURL: env.GITHUB_OAUTH_CALLBACK_URL,
				Scopes:      []string{"user:email", "read:user"},
			},
		},
	})
}
