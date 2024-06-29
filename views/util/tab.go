package util

import (
	"citadel/internal/middleware"
	"context"
)

func GetClassForTab(tabPath string, currentPath string) string {
	if tabPath == currentPath {
		return "text-yellow-300 whitespace-nowrap"
	}

	return "hover:text-yellow-300 transition-colors whitespace-nowrap"
}

func RetrievePath(ctx context.Context) string {
	ctxValue := ctx.Value(middleware.CTX_KEY_PATH)
	if ctxValue == nil {
		return ""
	}
	path, ok := ctxValue.(string)
	if !ok {
		return ""
	}

	return path
}
