package controllers

import (
	"citadel/app/drivers"
	"citadel/app/models"
	"citadel/app/repositories"
	storagePages "citadel/views/pages/storage"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
)

type StorageController struct {
	storageBucketsRepo *repositories.StorageBucketsRepository
	driver             drivers.Driver
}

func NewStorageController(storageBucketsRepo *repositories.StorageBucketsRepository, driver drivers.Driver) *StorageController {
	return &StorageController{storageBucketsRepo, driver}
}

func (c *StorageController) Index(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	storageBuckets, err := c.storageBucketsRepo.FindAllFromUser(ctx.Context(), user.ID)
	if err != nil {
		return caesar.NewError(400)
	}

	if ctx.WantsJSON() {
		return ctx.SendJSON(storageBuckets)
	}

	return ctx.Render(storagePages.Index(storageBuckets))
}

type StoreStorageBucketValidator struct {
	Name string `form:"name" validate:"required,min=3,lowercase"`
}

func (c *StorageController) Store(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	data, _, ok := caesar.Validate[StoreStorageBucketValidator](ctx)
	if !ok {
		return ctx.Redirect("/storage")
	}

	bucket := &models.StorageBucket{Name: data.Name, UserID: user.ID}
	if err := c.storageBucketsRepo.Create(ctx.Context(), bucket); err != nil {
		return err
	}

	host, keyId, secretKey, err := c.driver.CreateStorageBucket(*bucket)
	if err != nil {
		return err
	}
	bucket.Host = host
	bucket.KeyId = keyId
	bucket.SecretKey = secretKey

	if err := c.storageBucketsRepo.UpdateOneWhere(ctx.Context(), "id", bucket.ID, bucket); err != nil {
		return err
	}

	return ctx.Redirect("/storage")
}

func (c *StorageController) Show(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	bucket, err := c.storageBucketsRepo.FindOneBy(ctx.Context(), "slug", ctx.PathValue("slug"))
	if err != nil {
		return err
	}

	if bucket.UserID != user.ID {
		return caesar.NewError(403)
	}

	return ctx.Render(storagePages.Show(*bucket))
}

func (c *StorageController) Edit(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	bucket, err := c.storageBucketsRepo.FindOneBy(ctx.Context(), "slug", ctx.PathValue("slug"))
	if err != nil {
		return err
	}

	if bucket.UserID != user.ID {
		return caesar.NewError(403)
	}

	return ctx.Render(storagePages.Edit(*bucket))
}

func (c *StorageController) Delete(ctx *caesar.CaesarCtx) error {
	user, err := auth.RetrieveUserFromCtx[models.User](ctx)
	if err != nil {
		return err
	}

	bucket, err := c.storageBucketsRepo.FindOneBy(ctx.Context(), "slug", ctx.PathValue("slug"))
	if err != nil {
		return err
	}

	if bucket.UserID != user.ID {
		return caesar.NewError(403)
	}

	if err := c.storageBucketsRepo.DeleteOneWhere(ctx.Context(), "id", bucket.ID); err != nil {
		return err
	}

	return ctx.Render(storagePages.Edit(*bucket))
}
