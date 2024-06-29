package middleware

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	"context"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/vexillum"
)

type CTX_KEY string

const (
	CTX_KEY_URL                        CTX_KEY = "url"
	CTX_KEY_PATH                       CTX_KEY = "path"
	CTX_KEY_FLAGS                      CTX_KEY = "flags"
	CTX_KEY_ORGS                       CTX_KEY = "orgs"
	CTX_KEY_ORG_ID                     CTX_KEY = "orgId"
	CTX_KEY_SHOW_PAYMENT_METHOD_DIALOG CTX_KEY = "SHOW_PAYMENT_METHOD_DIALOG"
)

// ViewMiddleware is a middleware that injects data into the context
// (so that it can be used in the views).
func ViewMiddleware(vexillum *vexillum.Vexillum, orgsRepository *repositories.OrganizationsRepository) caesar.Handler {
	return func(ctx *caesar.Context) error {
		handlePaymentMethodDialog(vexillum, ctx)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), CTX_KEY_URL, ctx.Request.URL.String()),
		)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), CTX_KEY_PATH, ctx.Request.URL.Path),
		)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), CTX_KEY_FLAGS, vexillum.Flags),
		)

		if ctx.PathValue("orgId") != "" {
			user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
			if err != nil {
				return err
			}

			orgs, err := orgsRepository.FindAllWhereUserIsMember(ctx.Request.Context(), user.ID)
			if err != nil {
				return err
			}

			ctx.Request = ctx.Request.WithContext(
				context.WithValue(ctx.Request.Context(), CTX_KEY_ORG_ID, ctx.PathValue("orgId")),
			)

			ctx.Request = ctx.Request.WithContext(
				context.WithValue(ctx.Request.Context(), CTX_KEY_ORGS, orgs),
			)
		}

		ctx.Next()

		return nil
	}
}

// handlePaymentMethodDialog sets the "SHOW_PAYMENT_METHOD_DIALOG" param in the context,
// if the user has no active payment method and the billing feature is active.
func handlePaymentMethodDialog(vexillum *vexillum.Vexillum, ctx *caesar.Context) {
	if !vexillum.IsActive("billing") {
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), CTX_KEY_SHOW_PAYMENT_METHOD_DIALOG, false),
		)
		return
	}

	// We try to retrieve the user from the context.
	// If the user is not found, we return.
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return
	}

	ctx.Request = ctx.Request.WithContext(
		context.WithValue(ctx.Request.Context(), CTX_KEY_SHOW_PAYMENT_METHOD_DIALOG, !user.HasActivePaymentMethod()),
	)
}
