package authControllers

import (
	"citadel/app/models"
	"citadel/app/repositories"
	"citadel/app/services"
	authPages "citadel/views/pages/auth"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/events"
	"github.com/rs/xid"
)

type SignUpController struct {
	auth    *caesarAuth.Auth
	repo    *repositories.UsersRepository
	service *services.UsersService
	emitter *events.EventsEmitter
}

func NewSignUpController(auth *caesarAuth.Auth, service *services.UsersService, repo *repositories.UsersRepository, emitter *events.EventsEmitter) *SignUpController {
	return &SignUpController{auth, repo, service, emitter}
}

func (c *SignUpController) Show(ctx *caesar.CaesarCtx) error {
	return ctx.Render(authPages.SignUpPage())
}

type SignUpValidator struct {
	Email    string `form:"email" validate:"required,email"`
	FullName string `form:"full_name" validate:"required,min=3"`
	Password string `form:"password" validate:"required,min=8"`
}

func (c *SignUpController) Handle(ctx *caesar.CaesarCtx) error {
	data, errors, ok := caesar.Validate[SignUpValidator](ctx)
	if !ok {
		return ctx.Render(authPages.SignUpForm(
			authPages.SignUpInput{Email: data.Email, FullName: data.FullName, Password: data.Password},
			errors,
		))
	}

	user, _ := c.repo.FindOneBy(ctx.Context(), "email", data.Email)
	if user != nil {
		return ctx.Render(authPages.SignUpForm(
			authPages.SignUpInput{Email: data.Email, FullName: data.FullName, Password: data.Password},
			map[string]string{"Email": "Email is already in use."},
		))
	}

	hashedPassword, err := caesarAuth.HashPassword(data.Password)
	if err != nil {
		return caesar.NewError(400)
	}

	user = &models.User{ID: xid.New().String(), Email: data.Email, FullName: data.FullName, Password: hashedPassword}
	if err := c.service.CreateAndEmitEvent(ctx.Context(), user); err != nil {
		return caesar.NewError(400)
	}

	if err := c.auth.Authenticate(ctx, *user); err != nil {
		return caesar.NewError(400)
	}

	return ctx.Redirect("/applications")
}
