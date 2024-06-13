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

	StreamLogs(ctx *caesar.CaesarCtx, app models.Application) error

	CreateStorageBucket(bucket models.StorageBucket) (host string, keyId string, secretKey string, err error)
	DeleteStorageBucket(bucket models.StorageBucket) error
}
