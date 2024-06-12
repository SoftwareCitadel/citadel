package config

import (
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func ProvideStripe(env *EnvironmentVariables) *client.API {
	config := &stripe.BackendConfig{}
	sc := &client.API{}
	sc.Init(env.STRIPE_SECRET_KEY, &stripe.Backends{
		API:     stripe.GetBackendWithConfig(stripe.APIBackend, config),
		Uploads: stripe.GetBackendWithConfig(stripe.UploadsBackend, config),
	})
	return sc
}
