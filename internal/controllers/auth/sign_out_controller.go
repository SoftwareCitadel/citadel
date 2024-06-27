package authControllers

// Uncomment the following import statement once you implement your first controller method.
import (
	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type SignOutController struct {
	auth *auth.Auth
}

func NewSignOutController(auth *auth.Auth) *SignOutController {
	return &SignOutController{
		auth: auth,
	}
}

func (c *SignOutController) Handle(ctx *caesar.Context) error {
	if err := c.auth.SignOut(ctx); err != nil {
		return err
	}

	return ctx.Redirect("/auth/sign_in")
}
