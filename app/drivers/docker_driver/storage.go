package dockerDriver

import (
	"citadel/app/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/minio/minio-go/v7"
)

func (d *DockerDriver) CreateStorageBucket(bucket models.StorageBucket) (host string, keyId string, secretKey string, err error) {
	exists, err := d.minioClient.BucketExists(context.Background(), bucket.Slug)
	if err != nil {
		return "", "", "", err
	}
	if exists {
		return "", "", "", errors.New("bucket already exists")
	}

	if err := d.minioClient.MakeBucket(context.Background(), bucket.Slug, minio.MakeBucketOptions{}); err != nil {
		return "", "", "", err
	}

	newAccessKey := bucket.ID
	newSecretKey, err := generateSecretKey()
	if err != nil {
		return "", "", "", err
	}

	if err := d.minioAdmin.AddUser(context.Background(), newAccessKey, newSecretKey); err != nil {
		return "", "", "", err
	}

	return d.ipv4, newAccessKey, newSecretKey, nil
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

// generateSecretKey generates a random secret key of the specified length.
func generateSecretKey() (string, error) {
	const (
		MIN_LENGTH = 32
		MAX_LENGTH = 64
	)

	// Generate a random number between min and max.
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(MAX_LENGTH-MIN_LENGTH+1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random length: %v", err)
	}

	// Convert to int and adjust for the range offset.
	length := int(randomBigInt.Int64()) + MIN_LENGTH

	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer")
	}

	// Create a byte slice to hold the random bytes.
	key := make([]byte, length)

	// Read random bytes from the crypto/rand reader.
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate secret key: %v", err)
	}

	// Convert the byte slice to a hexadecimal string.
	secretKey := hex.EncodeToString(key)
	return secretKey, nil
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
