package authControllers

import (
	"citadel/internal/repositories"
	authPages "citadel/views/concerns/auth/pages"
	"os"
	"time"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/golang-jwt/jwt"
)

type ResetPwdController struct {
	usersRepo *repositories.UsersRepository
}

func NewResetPwdController(usersRepo *repositories.UsersRepository) *ResetPwdController {
	return &ResetPwdController{usersRepo}
}

func (c *ResetPwdController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.ResetPasswordPage())
}

type ResetPwdValidator struct {
	Password        string `form:"password" validate:"required,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
}

func (c *ResetPwdController) Handle(ctx *caesar.Context) error {
	// Fetch the JWT token from the URL parameter
	tokenString := ctx.PathValue("jwt")

	// Parse and validate the JWT token
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_KEY")), nil
	})
	if err != nil || !validateTokenClaims(claims) {
		return err
	}

	// Validate the submitted form data
	data, errors, ok := caesar.Validate[ResetPwdValidator](ctx)
	if !ok {
		return ctx.Render(authPages.ResetPasswordForm(errors))
	}

	// Fetch the user by ID from the token claims
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return err
	}

	user, err := c.usersRepo.FindOneBy(ctx.Context(), "id", userID)
	if err != nil || user == nil {
		return err
	}

	// Update the user's password in the database
	pwd, err := auth.HashPassword(data.Password)
	if err != nil {
		return err
	}

	user.Password = pwd
	if err := c.usersRepo.UpdateOneWhere(ctx.Context(), user, "id", userID); err != nil {
		return err
	}

	return ctx.Render(authPages.ResetPasswordSuccessAlert())
}

func validateTokenClaims(claims *jwt.MapClaims) bool {
	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return false
		}
	}
	return true
}
