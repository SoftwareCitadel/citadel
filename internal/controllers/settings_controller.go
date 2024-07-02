package controllers

import (
	"citadel/internal/models"
	"citadel/internal/repositories"

	"citadel/views/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/events"
	"github.com/caesar-rocks/ui/toast"
)

type SettingsController struct {
	emitter *events.EventsEmitter
	repo    *repositories.UsersRepository
}

func NewSettingsController(emitter *events.EventsEmitter, repo *repositories.UsersRepository) *SettingsController {
	return &SettingsController{emitter, repo}
}

func (c *SettingsController) Edit(ctx *caesar.Context) error {
	return ctx.Render(pages.SettingsPage())
}

type SettingsValidator struct {
	Email           string `form:"email" validate:"required,email"`
	FullName        string `form:"full_name" validate:"required,min=3"`
	NewPassword     string `form:"new_password" validate:"omitempty,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"eqfield=NewPassword"`
}

func (c *SettingsController) Update(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	data, errors, ok := caesar.Validate[SettingsValidator](ctx)
	if !ok {
		return ctx.Render(pages.SettingsForm(errors))
	}

	user.Email = data.Email
	user.FullName = data.FullName
	if data.NewPassword != "" {
		hashedPassword, err := caesarAuth.HashPassword(data.NewPassword)
		if err != nil {
			return caesar.NewError(400)
		}

		user.Password = hashedPassword
	}

	if err := c.repo.UpdateOneWhere(ctx.Context(), user, "id", user.ID); err != nil {
		return caesar.NewError(500)
	}

	toast.Success(ctx, "Settings updated successfully.")

	return ctx.Render(pages.SettingsForm(nil))
}

func (c *SettingsController) Delete(ctx *caesar.Context) error {
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return caesar.NewError(400)
	}

	if err := c.repo.DeleteOneWhere(ctx.Context(), "id", user.ID); err != nil {
		return caesar.NewError(500)
	}

	return ctx.Redirect("/auth/sign_up")
}
