package util

import (
	"citadel/internal/middleware"
	"context"
)

func IsFeatureActive(ctx context.Context, feature string) bool {
	flags := ctx.Value(middleware.CTX_KEY_FLAGS).(map[string]bool)
	return flags[feature]
}
