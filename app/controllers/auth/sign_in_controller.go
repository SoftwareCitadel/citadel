package authControllers

import (
	"citadel/app/repositories"
	authPages "citadel/views/concerns/auth/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type SignInController struct {
	auth            *caesarAuth.Auth
	usersRepository *repositories.UsersRepository
}

func NewSignInController(auth *caesarAuth.Auth, usersRepository *repositories.UsersRepository) *SignInController {
	return &SignInController{
		auth:            auth,
		usersRepository: usersRepository,
	}
}

func (c *SignInController) Show(ctx *caesar.Context) error {
	return ctx.Render(authPages.SignInPage())
}

type SignInValidator struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
}

func (c *SignInController) Handle(ctx *caesar.Context) error {
	data, errors, ok := caesar.Validate[SignInValidator](ctx)
	if !ok {
		return ctx.Render(authPages.SignInForm(
			authPages.SignInInput{Email: data.Email, Password: data.Password},
			errors,
		))
	}

	user, _ := c.usersRepository.FindOneBy(ctx.Context(), "email", data.Email)
	if user == nil {
		return ctx.Render(authPages.SignInForm(
			authPages.SignInInput{Email: data.Email, Password: data.Password},
			map[string]string{"Auth": "Invalid credentials."},
		))
	}

	if !caesarAuth.CheckPasswordHash(data.Password, user.Password) {
		return ctx.Render(authPages.SignInForm(
			authPages.SignInInput{Email: data.Email, Password: data.Password},
			map[string]string{"Auth": "Invalid credentials."},
		))
	}

	if err := c.auth.Authenticate(ctx, *user); err != nil {
		return err
	}

	return ctx.Redirect("/apps")
}
