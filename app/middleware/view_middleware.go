package middleware

import (
	"citadel/app/models"
	"citadel/app/vexillum"
	"context"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

// ViewMiddleware is a middleware that injects data into the context
// (so that it can be used in the views).
func ViewMiddleware(vexillum *vexillum.Vexillum) caesar.Handler {
	return func(ctx *caesar.Context) error {
		handlePaymentMethodDialog(vexillum, ctx)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), "url", ctx.Request.URL.String()),
		)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), "path", ctx.Request.URL.Path),
		)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), "flags", vexillum.Flags),
		)

		ctx.Next()

		return nil
	}
}

// handlePaymentMethodDialog sets the "SHOW_PAYMENT_METHOD_DIALOG" param in the context,
// if the user has no active payment method and the billing feature is active.
func handlePaymentMethodDialog(vexillum *vexillum.Vexillum, ctx *caesar.Context) {
	if !vexillum.IsActive("billing") {
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), "SHOW_PAYMENT_METHOD_DIALOG", false),
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
		context.WithValue(ctx.Request.Context(), "SHOW_PAYMENT_METHOD_DIALOG", !user.HasActivePaymentMethod()),
	)
}
