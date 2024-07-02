package util

import (
	"citadel/internal/models"
	"context"
	"strings"

	"github.com/caesar-rocks/auth"
)

func RetrieveInitials(ctx context.Context) string {
	ctxValue := ctx.Value(auth.USER_CONTEXT_KEY)
	if ctxValue == nil {
		return ""
	}
	user, ok := ctxValue.(*models.User)
	if !ok {
		return ""
	}

	initials := ""

	for _, part := range strings.Split(user.FullName, " ") {
		if len(part) == 0 {
			continue
		}
		initials += string(part[0])
	}

	return initials
}
