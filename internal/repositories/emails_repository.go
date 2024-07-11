package repositories

import (
	"citadel/internal/models"
	"context"

	"github.com/Squwid/go-randomizer"
	"github.com/caesar-rocks/orm"
	"github.com/gosimple/slug"
)

type EmailRepository struct {
	*orm.Repository[models.Email]
}

func NewEmailRepository(db *orm.Database) *EmailRepository {
	return &EmailRepository{Repository: &orm.Repository[models.Email]{Database: db}}
}
func (r EmailRepository) Create(ctx context.Context, e *models.Email) error {
	slug := slug.Make(e.Subject)

	for {
		_, err := r.FindOneBy(ctx, "slug", slug)
		if err != nil {
			break
		}

		slug = slug + "-" + randomizer.Noun()
	}
	e.Slug = slug

	_, err := r.NewInsert().Model(e).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve emails from the database
// This function allows the user to retrieve all emails, emails from a specific sender, and emails to a specific recipient
// Limiting and offset have been implemented to provide efficiency and performance for querying the database
func (r EmailRepository) FindEmails(ctx context.Context, orgID string, filter map[string]string, limit, offset int) ([]models.Email, error) {
    var emails []models.Email
    query := r.NewSelect().Model(&emails).Where("organization_id = ?", orgID)

    // Apply filters
    if sender, ok := filter["sender"]; ok && sender != "" {
        query = query.Where("sender = ?", sender)
    }
    if recipient, ok := filter["recipient"]; ok && recipient != "" {
        query = query.Where("recipient = ?", recipient)
    }

    // Apply sorting, limit, and offset
    err := query.
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Scan(ctx)

    if err != nil {
        return nil, err
    }
    return emails, nil
}
