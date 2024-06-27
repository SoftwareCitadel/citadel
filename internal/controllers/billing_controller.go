package controllers

import (
	"citadel/internal/models"
	"citadel/views/pages"
	"os"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

type BillingController struct {
	stripeClientApi *client.API
}

func NewBillingController(stripeClientApi *client.API) *BillingController {
	return &BillingController{stripeClientApi}
}

func (c *BillingController) Show(ctx *caesar.Context) error {
	return ctx.Render(pages.BillingPage())
}

func (c *BillingController) Manage(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	s, err := c.stripeClientApi.BillingPortalSessions.New(&stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(os.Getenv("APP_URL") + "/billing"),
	})
	if err != nil {
		return err
	}

	return ctx.Redirect(s.URL)
}

func (c *BillingController) InitiatePaymentMethodChange(ctx *caesar.Context) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	redirectUrl := os.Getenv("APP_URL") + "/billing"
	session, err := c.stripeClientApi.CheckoutSessions.New(&stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("setup"),
		Customer:           stripe.String(user.StripeCustomerID),
		SuccessURL:         stripe.String(redirectUrl),
		CancelURL:          stripe.String(redirectUrl),
	})
	if err != nil {
		return err
	}

	return ctx.Redirect(session.URL)
}
