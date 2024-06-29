package util

import (
	"citadel/internal/middleware"
	"context"
)

func Route(ctx context.Context, suffix string) string {
	return "/orgs/" + ctx.Value(middleware.CTX_KEY_ORG_ID).(string) + suffix
}
