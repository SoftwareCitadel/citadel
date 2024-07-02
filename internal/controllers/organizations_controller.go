package controllers

import (
	"citadel/internal/middleware"
	"citadel/internal/models"
	"citadel/internal/repositories"
	orgsPages "citadel/views/concerns/orgs/pages"

	caesarAuth "github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/ui/toast"
)

type OrganizationsController struct {
	orgsRepo       *repositories.OrganizationsRepository
	orgMembersRepo *repositories.OrganizationMembersRepository
}

func NewOrganizationsController(orgsRepo *repositories.OrganizationsRepository, orgMembersRepo *repositories.OrganizationMembersRepository) *OrganizationsController {
	return &OrganizationsController{orgsRepo, orgMembersRepo}
}

type StoreOrgValidator struct {
	Name string `form:"name" validate:"required,min=3"`
}

func (c *OrganizationsController) Store(ctx *caesar.Context) error {
	// Retrieve the current user
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	// Validate the request data
	data, _, ok := caesar.Validate[StoreOrgValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	// Create the organization
	org := &models.Organization{
		Name: data.Name,
	}
	if err := c.orgsRepo.Create(ctx.Request.Context(), org); err != nil {
		return err
	}

	// Add the user as an owner of the organization
	member := &models.OrganizationMember{UserID: user.ID, OrganizationID: org.ID, Role: models.OrganizationMemberRoleOwner}
	if err := c.orgMembersRepo.Create(ctx.Context(), member); err != nil {
		return err
	}

	toast.Success(ctx, "Organization created successfully")

	return ctx.Redirect("/orgs/" + org.ID + "/apps")
}

func (c *OrganizationsController) Edit(ctx *caesar.Context) error {
	// Retrieve the current organization,
	// from the context (set in the view middleware)
	orgs, ok := ctx.Request.Context().Value(middleware.CTX_KEY_ORGS).([]models.Organization)
	if !ok {
		return caesar.NewError(500)
	}

	var org models.Organization
	for _, o := range orgs {
		if o.ID == ctx.PathValue("orgId") {
			org = o
			break
		}
	}

	if org.ID == "" {
		return caesar.NewError(404)
	}

	members, err := c.orgMembersRepo.FindAllFromOrganizationWithUser(ctx.Request.Context(), org.ID)
	if err != nil {
		return err
	}

	return ctx.Render(orgsPages.Edit(org, members, *org.OrganizationMembers[0]))
}

type UpdateOrgValidator struct {
	Name string `form:"name" validate:"required,min=3"`
}

func (c *OrganizationsController) Update(ctx *caesar.Context) error {
	org, err := c.getCurrentOrgAndCheckOwnership(ctx)
	if err != nil {
		return err
	}

	// Validate the request data
	data, errors, ok := caesar.Validate[UpdateOrgValidator](ctx)
	if !ok {
		return ctx.Render(orgsPages.UpdateOrgForm(*org, errors))
	}

	// Update the organization
	org.Name = data.Name
	if err := c.orgsRepo.UpdateOneWhere(ctx.Request.Context(), org, "id", org.ID); err != nil {
		return err
	}

	toast.Success(ctx, "Organization updated successfully")

	return ctx.Render(orgsPages.UpdateOrgForm(*org, nil))
}

func (c *OrganizationsController) Delete(ctx *caesar.Context) error {
	org, err := c.getCurrentOrgAndCheckOwnership(ctx)
	if err != nil {
		return err
	}

	if err := c.orgsRepo.DeleteOneWhere(ctx.Request.Context(), "id", org.ID); err != nil {
		return err
	}

	return ctx.Redirect("/")
}

func (c *OrganizationsController) getCurrentOrgAndCheckOwnership(ctx *caesar.Context) (*models.Organization, error) {
	// Retrieve the current user
	user, err := caesarAuth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve the org from the context, as set in the view middleware
	orgs, ok := ctx.Request.Context().Value(middleware.CTX_KEY_ORGS).([]models.Organization)
	if !ok {
		return nil, caesar.NewError(500)
	}

	var currentOrg models.Organization
	for _, o := range orgs {
		if o.ID == ctx.PathValue("orgId") {
			currentOrg = o
			break
		}
	}

	// Check if the user is an owner of the organization
	orgOwned := false
	for _, member := range currentOrg.OrganizationMembers {
		if member.UserID == user.ID {
			if member.Role == models.OrganizationMemberRoleOwner {
				orgOwned = true
			}
			break
		}
	}

	if !orgOwned {
		return nil, caesar.NewError(403)
	}

	return &currentOrg, nil
}
