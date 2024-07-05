package util

import (
	"context"

	"github.com/caesar-rocks/auth"

	"citadel/internal/models"
)

func RetrieveUser(ctx context.Context) *models.User {
	ctxValue := ctx.Value(auth.USER_CONTEXT_KEY)
	if ctxValue == nil {
		return nil
	}
	user, ok := ctxValue.(*models.User)
	if !ok {
		return nil
	}

	return user
}
