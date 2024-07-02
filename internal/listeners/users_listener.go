package listeners

import (
	"citadel/internal/models"
	"citadel/internal/repositories"
	"citadel/util"
	"context"

	"github.com/caesar-rocks/vexillum"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

type UsersListener struct {
	stripe    *client.API
	usersRepo *repositories.UsersRepository
	vexillum  *vexillum.Vexillum
}

func NewUsersListener(stripe *client.API, usersRepo *repositories.UsersRepository, vexillum *vexillum.Vexillum) *UsersListener {
	return &UsersListener{stripe, usersRepo, vexillum}
}

func (usersListener *UsersListener) OnCreated(msg *message.Message) ([]*message.Message, error) {
	var user models.User
	if err := util.DecodeJSON(msg.Payload, &user); err != nil {
		return nil, err
	}

	if usersListener.vexillum.IsActive("billing") {
		if err := usersListener.assignStripeCustomer(&user); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (usersListener *UsersListener) assignStripeCustomer(user *models.User) error {
	cus, err := usersListener.stripe.Customers.New(&stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.FullName),
	})
	if err != nil {
		return err
	}

	user.StripeCustomerID = cus.ID

	if err := usersListener.usersRepo.UpdateOneWhere(context.Background(), user, "id", user.ID); err != nil {
		return err
	}

	return nil
}
