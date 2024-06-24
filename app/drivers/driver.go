package drivers

import (
	"citadel/app/models"

	caesar "github.com/caesar-rocks/core"
)

type Driver interface {
	Init() error

	CreateApplication(app models.Application) error
	DeleteApplication(app models.Application) error

	CreateCertificate(app models.Application, cert models.Certificate) ([]models.DnsEntry, error)
	CheckDnsConfig(app models.Application, cert models.Certificate) (bool, error)
	DeleteCertificate(app models.Application, cert models.Certificate) error

	IgniteBuilder(app models.Application, depl models.Deployment) error
	IgniteApplication(app models.Application, depl models.Deployment) error

	StreamLogs(ctx *caesar.Context, app models.Application) error

	// Database-related methods
	CreateDatabase(db models.Database) error
	DeleteDatabase(db models.Database) error

	// Storage-related methods
	CreateStorageBucket(bucket models.StorageBucket) (host string, keyId string, secretKey string, region string, err error)
	GetFilesAndTotalSize(bucket models.StorageBucket) (totalSize float64, files []models.StorageFile, err error)
	DeleteStorageBucket(bucket models.StorageBucket) error
}
