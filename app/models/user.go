package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	ID                                string          `bun:"id,pk"`
	Email                             string          `bun:"email,unique"`
	FullName                          string          `bun:"full_name,notnull"`
	Password                          string          `bun:"password,notnull"`
	GitHubUserID                      string          `bun:"github_user_id,unique,default:NULL"`
	GitHubInstallationIDs             json.RawMessage `bun:"github_installation_ids,type:jsonb,default:'[]'"`
	StripeCustomerID                  string          `bun:"stripe_customer_id,unique,default:NULL"`
	StripePaymentMethodID             string          `bun:"stripe_payment_method_id,unique,default:NULL"`
	StripePaymentMethodExpirationDate time.Time       `bun:"stripe_payment_method_expiration_date,default:NULL"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (user *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		user.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		user.UpdatedAt = time.Now()
	}
	return nil
}

// GetGitHubInstallationIDs returns the GitHub installation IDs of the user.
func (user *User) GetGitHubInstallationIDs() []int64 {
	var ids []int64
	if err := json.Unmarshal(user.GitHubInstallationIDs, &ids); err != nil {
		return nil
	}
	return ids
}

// AddGitHubInstallationId adds the given GitHub installation IDs to the user.
func (user *User) AddGitHubInstallationId(id int64) error {
	installationIds := user.GetGitHubInstallationIDs()
	installationIds = append(installationIds, id)

	data, err := json.Marshal(installationIds)
	if err != nil {
		return err
	}
	user.GitHubInstallationIDs = data
	return nil
}

// RemoveGitHubInstallationId removes the given GitHub installation IDs from the user.
func (user *User) RemoveGitHubInstallationId(id int64) error {
	installationIds := user.GetGitHubInstallationIDs()
	for i, installationId := range installationIds {
		if installationId == id {
			installationIds = append(installationIds[:i], installationIds[i+1:]...)
			break
		}
	}

	data, err := json.Marshal(installationIds)
	if err != nil {
		return err
	}
	user.GitHubInstallationIDs = data
	return nil
}

func (u *User) HasActivePaymentMethod() bool {
	if u.StripePaymentMethodID == "" {
		return false
	}

	return u.StripePaymentMethodExpirationDate.After(time.Now())
}
