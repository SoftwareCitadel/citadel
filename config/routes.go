package config

import (
	"citadel/internal/controllers"
	apiControllers "citadel/internal/controllers/api"
	authControllers "citadel/internal/controllers/auth"
	"citadel/internal/middleware"
	"citadel/internal/models"
	"citadel/internal/repositories"

	mailsPages "citadel/views/concerns/mails/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/events"
	"github.com/caesar-rocks/vexillum"
)

func RegisterRoutes(
	auth *caesarAuth.Auth,
	signUpController *authControllers.SignUpController,
	signInController *authControllers.SignInController,
	signOutController *authControllers.SignOutController,
	forgotPwdController *authControllers.ForgotPwdController,
	resetPwdController *authControllers.ResetPwdController,
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
	analyticsWebsitesController *controllers.AnalyticsWebsitesController,
	orgsRepository *repositories.OrganizationsRepository,
	orgsController *controllers.OrganizationsController,
	emitter *events.EventsEmitter,
	vexillum *vexillum.Vexillum,
) *caesar.Router {
	router := caesar.NewRouter()

	// Middleware
	router.Use(auth.SilentMiddleware)
	router.Use(middleware.ViewMiddleware(vexillum, orgsRepository))

	// Home route
	router.Get("/", func(ctx *caesar.Context) error {
		user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
		if err != nil {
			return ctx.Redirect("/auth/sign_in")
		}

		// Redirect to owned organization if user has one
		org, err := orgsRepository.FindFirstOwnedByUser(ctx.Context(), user.ID)
		if err == nil {
			return ctx.Redirect("/orgs/" + org.ID + "/apps")
		}

		return ctx.Redirect("/orgs/no_org/apps")
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

	// Reset password routes
	router.Get("/auth/reset_password/{jwt}", resetPwdController.Show)
	router.Post("/auth/reset_password/{jwt}", resetPwdController.Handle)

	// Apps CRUD routes
	router.Get("/orgs/{orgId}/apps", appsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/apps", appsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.Get("/orgs/{orgId}/apps/{slug}", appsController.Show).Use(auth.AuthMiddleware)
	router.Get("/orgs/{orgId}/apps/{slug}/edit", appsController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/orgs/{orgId}/apps/{slug}", appsController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/orgs/{orgId}/apps/{slug}", appsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Apps-related routes
	router.Get("/orgs/{orgId}/apps/{slug}/certs", certsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/apps/{slug}/certs", certsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Post("/orgs/{orgId}/apps/{slug}/certs/{id}/check", certsController.Check).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/orgs/{orgId}/apps/{slug}/certs/{id}", certsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Databases-related routes
	router.
		Get("/orgs/{orgId}/databases", databasesController.Index).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Post("/orgs/{orgId}/databases", databasesController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/orgs/{orgId}/databases/{slug}", databasesController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Mails-related routes
	router.Render("/orgs/{orgId}/mails", mailsPages.OverviewPage())
	router.
		Get("/orgs/{orgId}/mails/domains", mailDomainsController.Index).
		Use(auth.AuthMiddleware)
	router.
		Get("/orgs/{orgId}/mails/domains/{id}", mailDomainsController.Show).
		Use(auth.AuthMiddleware)
	router.Post("/orgs/{orgId}/mails/domains", mailDomainsController.Store).Use(auth.AuthMiddleware)
	router.
		Delete("/orgs/{orgId}/mails/domains/{id}", mailDomainsController.Delete).
		Use(auth.AuthMiddleware)
	router.Post("/orgs/{orgId}/mails/domains/check_dns/{id}", mailDomainsController.CheckDNS).Use(auth.AuthMiddleware)

	router.
		Get("/orgs/{orgId}/mails/api_keys", mailApiKeysController.Index).
		Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/mails/api_keys", mailApiKeysController.Store).
		Use(auth.AuthMiddleware)
	router.
		Patch("/orgs/{orgId}/mails/api_keys/{id}", mailApiKeysController.Update).
		Use(auth.AuthMiddleware)
	router.
		Delete("/orgs/{orgId}/mails/api_keys/{id}", mailApiKeysController.Delete).
		Use(auth.AuthMiddleware)

	// Logs-related routes
	router.Get("/orgs/{orgId}/apps/{slug}/logs", logsController.Index).
		Use(auth.AuthMiddleware)
	router.Get("/orgs/{orgId}/apps/{slug}/logs/stream", logsController.Stream).
		Use(auth.AuthMiddleware)

	// Storage-related routes
	router.
		Get("/orgs/{orgId}/storage", storageController.Index).
		Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/storage", storageController.Store).
		Use(auth.AuthMiddleware)
	router.
		Get("/orgs/{orgId}/storage/{slug}", storageController.Show).
		Use(auth.AuthMiddleware)
	router.
		Get("/orgs/{orgId}/storage/{slug}/edit", storageController.Edit).
		Use(auth.AuthMiddleware)
	router.
		Put("/orgs/{orgId}/storage/{slug}/edit", storageController.Update).
		Use(auth.AuthMiddleware)
	router.
		Delete("/orgs/{orgId}/storage/{slug}", storageController.Delete).
		Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/storage/{slug}/upload", storageController.UploadFile).
		Use(auth.AuthMiddleware)

	// Environment variables-related routes
	router.Get("/orgs/{orgId}/apps/{slug}/env", envController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/orgs/{orgId}/apps/{slug}/env", envController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// Deployments-related routes
	router.Get("/orgs/{orgId}/apps/{slug}/deployments", deploymentsController.Index).Use(auth.AuthMiddleware)
	router.
		Post("/orgs/{orgId}/apps/{slug}/deployments", deploymentsController.Store).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.Get("/orgs/{orgId}/apps/{slug}/deployments/list", deploymentsController.List).Use(auth.AuthMiddleware)

	// Billing-related routes
	router.
		Get("/billing", billingController.Show).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))
	router.
		Get("/billing/manage", billingController.Manage).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))
	router.
		Get("/billing/payment_method", billingController.InitiatePaymentMethodChange).
		Use(auth.AuthMiddleware).
		Use(vexillum.EnsureFeatureEnabledMiddleware("billing"))

	// Organizations-related routes
	router.
		Post("/orgs", orgsController.Store).
		Use(auth.AuthMiddleware)
	router.
		Get("/orgs/{orgId}/edit", orgsController.Edit).
		Use(auth.AuthMiddleware)
	router.
		Patch("/orgs/{orgId}", orgsController.Update).
		Use(auth.AuthMiddleware)
	router.
		Delete("/orgs/{orgId}", orgsController.Delete).
		Use(auth.AuthMiddleware)

	// Settings-related routes
	router.Get("/orgs/{orgId}/settings", settingsController.Edit).Use(auth.AuthMiddleware)
	router.
		Patch("/orgs/{orgId}/settings", settingsController.Update).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))
	router.
		Delete("/orgs/{orgId}/settings", settingsController.Delete).
		Use(auth.AuthMiddleware).
		Use(middleware.PaymentMethodMiddleware(vexillum))

	// GitHub-related routes
	router.Get("/github/repositories", githubController.ListRepositories).Use(auth.AuthMiddleware)

	router.Post("/orgs/{orgId}/apps/{slug}/connect_github", appsController.ConnectGitHub).Use(auth.AuthMiddleware)
	router.Post("/orgs/{orgId}/apps/{slug}/disconnect_github", appsController.DisconnectGitHub).Use(auth.AuthMiddleware)

	// Webhooks routes
	router.Post("/webhooks/github", githubController.HandleWebhook)
	router.Post("/webhooks/stripe", stripeController.HandleWebhook)

	// Analytics-related routes
	router.Get("/orgs/{orgId}/analytics_websites", analyticsWebsitesController.Index).Use(auth.AuthMiddleware)
	router.Post("/orgs/{orgId}/analytics_websites", analyticsWebsitesController.Store).Use(auth.AuthMiddleware)
	router.Get("/orgs/{orgId}/analytics_websites/{id}", analyticsWebsitesController.Show).Use(auth.AuthMiddleware)
	router.Get("/orgs/{orgId}/analytics_websites/{id}/edit", analyticsWebsitesController.Edit).Use(auth.AuthMiddleware)
	router.Patch("/orgs/{orgId}/analytics_websites/{id}", analyticsWebsitesController.Update).Use(auth.AuthMiddleware)
	router.Delete("/orgs/{orgId}/analytics_websites/{id}", analyticsWebsitesController.Delete).Use(auth.AuthMiddleware)

	// API-related routes
	router.Get("/api/v1/emails", emailsController.Send).Use(auth.AuthMiddleware)
	router.Post("/api/v1/analytics_websites/{id}/track", analyticsWebsitesController.Track)

	return router
}
