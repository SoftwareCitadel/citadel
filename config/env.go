package config

import (
	"github.com/caesar-rocks/core"
	"github.com/caesar-rocks/orm"
)

// EnvironmentVariables is a struct that holds all the environment variables that need to be validated.
// For full reference, see: https://github.com/go-playground/validator.
type EnvironmentVariables struct {
	// APP_KEY is the key used for encryption and decryption.
	APP_KEY string `validate:"required"`

	// Addr is the address to listen on for incoming requests.
	ADDR string `validate:"required"`

	// APP_URL is the URL of the application.
	APP_URL string `validate:"required"`

	// DBMS is the database management system to use ("postgres", "mysql", "sqlite").
	DBMS orm.DBMS `validate:"oneof=postgres mysql sqlite"`

	// DSN is the data source name, which is a connection string for the database.
	DSN string `validate:"required"`

	// GITHUB_OAUTH_KEY is the key for the GitHub OAuth application.
	GITHUB_OAUTH_KEY string

	// GITHUB_OAUTH_SECRET is the secret for the GitHub OAuth application.
	GITHUB_OAUTH_SECRET string

	// GITHUB_OAUTH_CALLBACK_URL is the callback URL for the GitHub OAuth application.
	GITHUB_OAUTH_CALLBACK_URL string

	// STRIPE_SECRET_KEY is the key for the Stripe API.
	STRIPE_SECRET_KEY string

	// STRIPE_PUBLIC_KEY is the public key for the Stripe API.
	STRIPE_PUBLIC_KEY string

	// S3_KEY is the key for the S3 API.
	S3_KEY string `validate:"required"`

	// S3_SECRET is the secret for the S3 API.
	S3_SECRET string `validate:"required"`

	// S3_REGION is the region for the S3 API.
	S3_REGION string `validate:"required"`

	// S3_ENDPOINT is the endpoint for the S3 API.
	S3_ENDPOINT string `validate:"required"`

	// S3_BUCKET is the bucket for the S3 API.
	S3_BUCKET string `validate:"required"`

	// RESEND_KEY is the key for the Resend API.
	RESEND_KEY string `validate:"required"`

	// REDIS_ADDR is the address for the Redis server.
	REDIS_ADDR string `validate:"required"`

	// REDIS_PASSWORD is the password for the Redis server.
	REDIS_PASSWORD string

	// BUILDER_IMAGE is the image to use for the builder.
	BUILDER_IMAGE string `validate:"required"`

	// REGISTRY_HOST is the registry to use for the builder.
	REGISTRY_HOST string `validate:"required"`

	// REGISTRY_TOKEN is the token to use for the registry.
	REGISTRY_TOKEN string

	// 	GITHUB_APP_ID is the ID for the GitHub App.
	GITHUB_APP_ID string

	// GITHUB_APP_PRIVATE_KEY is the private key for the GitHub App.
	GITHUB_APP_PRIVATE_KEY string

	// GITHUB_APP_WEBHOOK_SECRET is the secret for the GitHub App webhook.
	GITHUB_APP_WEBHOOK_SECRET string

	// GITHUB_APP_PRIVATE_KEY_PATH is the path to the private key for the GitHub App.
	GITHUB_APP_PRIVATE_KEY_PATH string

	// DB_HOST is the host for the database.
	DB_HOST string
}

func ProvideEnvironmentVariables() *EnvironmentVariables {
	return core.ValidateEnvironmentVariables[EnvironmentVariables]()
}
