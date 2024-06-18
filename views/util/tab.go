package util

import "context"

func GetClassForTab(tabPath string, currentPath string) string {
	if tabPath == currentPath {
		return "text-yellow-300 whitespace-nowrap"
	}

	return "hover:text-yellow-300 transition-colors whitespace-nowrap"
}

func RetrievePath(ctx context.Context) string {
	ctxValue := ctx.Value("path")
	if ctxValue == nil {
		return ""
	}
	path, ok := ctxValue.(string)
	if !ok {
		return ""
	}

	return path
}
