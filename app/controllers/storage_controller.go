package controllers

import (
	"bytes"
	"citadel/app/drivers"
	"citadel/app/models"
	"citadel/app/repositories"
	storagePages "citadel/views/concerns/storage/pages"
	"io"

	"github.com/caesar-rocks/auth"
	caesar "github.com/caesar-rocks/core"
	"github.com/caesar-rocks/drive"
	"github.com/caesar-rocks/ui/toast"
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

	host, keyId, secretKey, region, err := c.driver.CreateStorageBucket(*bucket)
	if err != nil {
		return err
	}
	bucket.Host = host
	bucket.KeyId = keyId
	bucket.SecretKey = secretKey
	bucket.Region = region

	if err := c.storageBucketsRepo.UpdateOneWhere(ctx.Context(), "id", bucket.ID, bucket); err != nil {
		return err
	}

	return ctx.Redirect("/storage/" + bucket.Slug)
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

	bucketSize, storageFiles, err := c.driver.GetFilesAndTotalSize(*bucket)
	if err != nil {
		return err
	}

	return ctx.Render(storagePages.Show(*bucket, storageFiles, bucketSize))
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

func (c *StorageController) Update(ctx *caesar.CaesarCtx) error {
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

	data, _, ok := caesar.Validate[StoreStorageBucketValidator](ctx)
	if !ok {
		return ctx.RedirectBack()
	}

	bucket.Name = data.Name
	if err := c.storageBucketsRepo.UpdateOneWhere(ctx.Context(), "id", bucket.ID, bucket); err != nil {
		return err
	}

	toast.Success(ctx, "Storage bucket updated successfully.")

	return ctx.Render(storagePages.EditForm(*bucket))
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

	if err := c.driver.DeleteStorageBucket(*bucket); err != nil {
		return err
	}

	return ctx.Redirect("/storage")
}

func (c *StorageController) UploadFile(ctx *caesar.CaesarCtx) error {
	// Retrieve the bucket owned by the current user
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

	// Parse the file, and pass its contents into a buffer.
	ctx.Request.ParseMultipartForm(10 << 20) // 10 MB
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return err
	}

	// Create a new Drive instance.
	myDrive := drive.NewDrive(map[string]drive.FileSystem{
		"s3": &drive.S3{
			Key:            bucket.KeyId,
			Secret:         bucket.SecretKey,
			Region:         bucket.Region,
			Endpoint:       bucket.Host,
			Bucket:         bucket.Slug,
			ForcePathStyle: true,
		},
	})

	// Upload the file to the S3 bucket.
	if err := myDrive.Use("s3").Put(fileHeader.Filename, buf.Bytes()); err != nil {
		return err
	}

	return ctx.RedirectBack()
}
