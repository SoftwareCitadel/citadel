package middleware

import (
	"citadel/internal/models"

	"github.com/caesar-rocks/vexillum"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

// PaymentMethodMiddleware is a middleware that checks if the user has a payment method set.
func PaymentMethodMiddleware(vexillum *vexillum.Vexillum) caesar.Handler {
	return func(ctx *caesar.Context) error {
		user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
		if err != nil {
			return err
		}

		if vexillum.IsActive("billing") && !user.HasActivePaymentMethod() {
			if ctx.WantsJSON() {
				return ctx.SendJSON(map[string]interface{}{
					"error": "Payment method is required",
				}, 400)
			}
			return ctx.Redirect("/orgs/" + ctx.PathValue("orgId") + "/apps")
		}

		ctx.Next()

		return nil
	}
}
