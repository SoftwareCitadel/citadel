package dockerDriver

import (
	"citadel/app/models"
	"citadel/app/util"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
)

func (d *DockerDriver) CreateStorageBucket(bucket models.StorageBucket) (host string, keyId string, secretKey string, region string, err error) {
	exists, err := d.minioClient.BucketExists(context.Background(), bucket.Slug)
	if err != nil {
		return "", "", "", "", err
	}
	if exists {
		return "", "", "", "", errors.New("bucket already exists")
	}

	if err := d.minioClient.MakeBucket(context.Background(), bucket.Slug, minio.MakeBucketOptions{}); err != nil {
		return "", "", "", "", err
	}

	newAccessKey := bucket.ID
	newSecretKey, err := util.GenerateSecretKey()
	if err != nil {
		return "", "", "", "", err
	}

	if err := d.minioAdmin.AddUser(context.Background(), newAccessKey, newSecretKey); err != nil {
		return "", "", "", "", err
	}

	policyBytes := []byte(fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
		 {
		  "Effect": "Allow",
		  "Principal": {"AWS": ["arn:aws:iam::minio:user/%s"]},
		  "Action": [
		   "s3:PutObject",
		   "s3:GetObject",
		   "s3:ListBucket"
		  ],
		  "Resource": [
		   "arn:aws:s3:::%s/*"
		  ]
		 }
		]
	}`, newAccessKey, bucket.Slug))

	if err := d.minioAdmin.AddCannedPolicy(context.Background(), newAccessKey+"-policy", policyBytes); err != nil {
		log.Println("Failed to add canned policy")
		return "", "", "", "", err
	}

	if _, err := d.minioAdmin.AttachPolicy(context.Background(), madmin.PolicyAssociationReq{
		Policies: []string{newAccessKey + "-policy"},
		User:     newAccessKey,
	}); err != nil {
		return "", "", "", "", err
	}

	return os.Getenv("MINIO_HOST"), newAccessKey, newSecretKey, os.Getenv("MINIO_REGION"), nil
}

func (d *DockerDriver) DeleteStorageBucket(bucket models.StorageBucket) error {
	ctx := context.Background()

	objectCh := d.minioClient.ListObjects(ctx, bucket.Slug, minio.ListObjectsOptions{Recursive: true})

	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		if err := d.minioClient.RemoveObject(ctx, bucket.Slug, object.Key, minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}

	if err := d.minioClient.RemoveBucket(context.Background(), bucket.Slug); err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) GetFilesAndTotalSize(bucket models.StorageBucket) (float64, []models.StorageFile, error) {
	var totalSize float64
	files := make([]models.StorageFile, 0)

	objectCh := d.minioClient.ListObjects(context.Background(), bucket.Slug, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			return 0, nil, object.Err
		}
		files = append(files, models.StorageFile{
			Name:      object.Key,
			Size:      float64(object.Size),
			UpdatedAt: object.LastModified,
			Type:      object.ContentType,
		})
		totalSize += float64(object.Size)
	}

	return totalSize, files, nil
}
