package util

import "context"

func IsFeatureActive(ctx context.Context, feature string) bool {
	flags := ctx.Value("flags").(map[string]bool)
	return flags[feature]
}
