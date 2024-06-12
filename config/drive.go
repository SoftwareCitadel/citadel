package config

import (
	"github.com/caesar-rocks/drive"
)

func ProvideDrive(env *EnvironmentVariables) *drive.Drive {
	return drive.NewDrive(map[string]drive.FileSystem{
		"s3": &drive.S3{
			Key:      env.S3_KEY,
			Secret:   env.S3_SECRET,
			Region:   env.S3_REGION,
			Endpoint: env.S3_ENDPOINT,
			Bucket:   env.S3_BUCKET,
		},
	})
}
