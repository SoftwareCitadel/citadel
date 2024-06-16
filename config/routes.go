package config

import (
	"citadel/app/controllers"
	apiControllers "citadel/app/controllers/api"
	authControllers "citadel/app/controllers/auth"
	"citadel/app/middleware"
	"citadel/app/vexillum"
	mailsPages "citadel/views/pages/mails"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/events"
)

func RegisterRoutes(
	auth *caesarAuth.Auth,
	signUpController *authControllers.SignUpController,
	signInController *authControllers.SignInController,
	signOutController *authControllers.SignOutController,
	forgotPwdController *authControllers.ForgotPwdController,
	authGithubController *authControllers.GithubController,
	logsController *controllers.LogsController,
	appsController *controllers.AppsController,
	databasesController *controllers.DatabasesController,
	envController *controllers.EnvController,
	deploymentsController *controllers.DeploymentsController,
	certsController *controllers.CertsController,
	billingController *controllers.BillingController,
	settingsController *controllers.SettingsController,
	githubController *controllers.GithubController,
	authCliController *authControllers.CliController,
	stripeController *controllers.StripeController,
	mailDomainsController *controllers.MailDomainsController,
	mailApiKeysController *controllers.MailApiKeysController,
	storageController *controllers.StorageController,
	emailsController *apiControllers.EmailsController,
	emitter *events.EventsEmitter,
	vexillum *vexillum.Vexillum,
) *caesar.Router {
	router := caesar.NewRouter()

	// Middleware
	router.Use(auth.SilentMiddleware)
	router.Use(middleware.ViewMiddleware(vexillum))

	// Home route
	router.Get("/", func(ctx *caesar.CaesarCtx) error {
		return ctx.Redirect("/apps")
	})

	// Auth routes
	router.Get("/auth/sign_up", signUpController.Show)
	router.Post("/auth/sign_up", signUpController.Handle)

	router.Get("/auth/sign_in", signInController.Show)
	router.Post("/auth/sign_in", signInController.Handle)

	router.Post("/auth/sign_out", signOutController.Handle).Use(auth.AuthMiddleware)

	// OAuth-related routes
	router.Get("/auth/github/redirect", authGithubController.Redirect)
	router.Get("/auth/github/callback", authGithubController.Callback)

	// CLI authentication-related routes
	router.Get("/auth/cli", authCliController.GetSession)
	router.Get("/auth/cli/{sessionId}", authCliController.Show).Use(auth.AuthMiddleware)
	router.Post("/auth/cli/{sessionId}", authCliController.Handle).Use(auth.AuthMiddleware)
	router.Get("/auth/cli/check", authCliController.Check)
	router.Get("/auth/cli/wait/{sessionId}", authCliController.Wait)

	// Forgot password routes
	router.Get("/auth/forgot_password", forgotPwdController.Show)
	router.Post("/auth/forgot_password", forgotPwdController.Handle)

	// Apps CRUD routes
	router.Get("/apps", appsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/apps", appsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.Get("/apps/{slug}", appsController.Show).Use(auth.AuthMiddleware)
	router.Get("/apps/{slug}/edit", appsController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/apps/{slug}", appsController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/apps/{slug}", appsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Apps-related routes
	router.Get("/apps/{slug}/certs", certsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/apps/{slug}/certs", certsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Post("/apps/{slug}/certs/{id}/check", certsController.Check).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/apps/{slug}/certs/{id}", certsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Databases-related routes
	router.Get("/databases", databasesController.Index).Use(auth.AuthMiddleware)
	router.Post("/databases", databasesController.Store).Use(auth.AuthMiddleware)

	// Mails-related routes
	router.Render("/mails", mailsPages.OverviewPage())
	router.
		Get("/mails/domains", mailDomainsController.Index).
		Use(auth.AuthMiddleware)
	router.
		Get("/mails/domains/{domain}", mailDomainsController.Show).
		Use(auth.AuthMiddleware)
	router.
		Get("/mails/api_keys", mailApiKeysController.Index).
		Use(auth.AuthMiddleware)

	// Logs-related routes
	router.Get("/apps/{slug}/logs", logsController.Index).Use(auth.AuthMiddleware)
	router.Get("/apps/{slug}/logs/stream", logsController.Stream).Use(auth.AuthMiddleware)

	// Storage-related routes
	router.
		Get("/storage", storageController.Index).
		Use(auth.AuthMiddleware)
	router.
		Post("/storage", storageController.Store).
		Use(auth.AuthMiddleware)
	router.
		Get("/storage/{slug}", storageController.Show).
		Use(auth.AuthMiddleware)
	router.
		Get("/storage/{slug}/edit", storageController.Edit).
		Use(auth.AuthMiddleware)
	router.
		Put("/storage/{slug}/edit", storageController.Update).
		Use(auth.AuthMiddleware)
	router.
		Delete("/storage/{slug}", storageController.Delete).
		Use(auth.AuthMiddleware)

	// Environment variables-related routes
	router.Get("/apps/{slug}/env", envController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/apps/{slug}/env", envController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Deployments-related routes
	router.Get("/apps/{slug}/deployments", deploymentsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/apps/{slug}/deployments", deploymentsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.Get("/apps/{slug}/deployments/list", deploymentsController.List).Use(auth.AuthMiddleware)

	// Billing-related routes
	router.
		Get("/billing", billingController.Show).
		Use(auth.AuthMiddleware).
		Use(vexillum.Middleware("billing"))
	router.
		Get("/billing/manage", billingController.Manage).
		Use(auth.AuthMiddleware).
		Use(vexillum.Middleware("billing"))
	router.
		Get("/billing/payment_method", billingController.InitiatePaymentMethodChange).
		Use(auth.AuthMiddleware).
		Use(vexillum.Middleware("billing"))

	// Settings-related routes
	router.Get("/settings", settingsController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/settings", settingsController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/settings", settingsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// GitHub-related routes
	router.Get("/github/repositories", githubController.ListRepositories).Use(auth.AuthMiddleware)

	router.Post("/apps/{slug}/connect_github", appsController.ConnectGitHub).Use(auth.AuthMiddleware)
	router.Post("/apps/{slug}/disconnect_github", appsController.DisconnectGitHub).Use(auth.AuthMiddleware)

	// Webhooks routes
	router.Post("/webhooks/github", githubController.HandleWebhook)
	router.Post("/webhooks/stripe", stripeController.HandleWebhook)

	// API-related routes
	router.Get("/api/v1/emails", emailsController.Send).Use(auth.AuthMiddleware)

	return router
}
