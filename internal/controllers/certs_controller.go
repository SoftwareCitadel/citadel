package controllers

import (
	"citadel/internal/drivers"
	"citadel/internal/models"
	"citadel/internal/repositories"
	"citadel/internal/services"
	appsPages "citadel/views/concerns/apps/pages"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type CertsController struct {
	auth        *auth.Auth
	appsRepo    *repositories.ApplicationsRepository
	certsRepo   *repositories.CertificatesRepository
	driver      drivers.Driver
	appsService *services.AppsService
}

func NewCertsController(auth *auth.Auth, appsRepo *repositories.ApplicationsRepository, certsRepo *repositories.CertificatesRepository, driver drivers.Driver, appsService *services.AppsService) *CertsController {
	return &CertsController{auth, appsRepo, certsRepo, driver, appsService}
}

func (c *CertsController) Index(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	certs, err := c.certsRepo.FindAllFromApp(ctx.Context(), app.ID)
	if err != nil {
		return err
	}

	return ctx.Render(appsPages.CertsPage(*app, certs))
}

type StoreCertInput struct {
	Domain string `form:"domain" validate:"required"`
}

func (c *CertsController) Store(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreCertInput](ctx)
	if !ok {
		return err
	}

	cert := &models.Certificate{Domain: data.Domain, ApplicationID: app.ID}

	dnsEntries, err := c.driver.CreateCertificate(*app, *cert)
	if err != nil {
		return err
	}
	cert.DnsEntries = dnsEntries

	if ok, _ = c.driver.CheckDnsConfig(*app, *cert); ok {
		cert.Status = models.CertificateStatusVerified
	} else {
		cert.Status = models.CertificateStatusPending
	}

	if err := c.certsRepo.Create(ctx.Context(), cert); err != nil {
		return err
	}

	return ctx.RedirectBack()
}

func (c *CertsController) Delete(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	cert, err := c.certsRepo.FindOneBy(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"application_id", app.ID,
	)
	if err != nil {
		return err
	}

	if err := c.certsRepo.DeleteOneWhere(ctx.Context(), "id", cert.ID); err != nil {
		return err
	}

	return ctx.RedirectBack()
}

func (c *CertsController) Check(ctx *caesar.Context) error {
	app, err := c.appsService.GetAppOwnedByCurrentOrg(ctx)
	if err != nil {
		return err
	}

	cert, err := c.certsRepo.FindOneBy(
		ctx.Context(),
		"id", ctx.PathValue("id"),
		"application_id", app.ID,
	)
	if err != nil {
		return err
	}

	if ok, _ := c.driver.CheckDnsConfig(*app, *cert); ok {
		cert.Status = models.CertificateStatusVerified
	} else {
		cert.Status = models.CertificateStatusPending
	}

	if err := c.certsRepo.UpdateOneWhere(ctx.Context(), cert, "id", cert.ID); err != nil {
		return err
	}

	return ctx.RedirectBack()
}
