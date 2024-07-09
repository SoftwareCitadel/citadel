package ravelDriver

import (
	"citadel/internal/models"

	caesar "github.com/caesar-rocks/core"
)

// Ravel is a driver that does nothing
type Ravel struct{}

// New creates a new Ravel driver
func New() *Ravel {
	return &Ravel{}
}

// Init does nothing and returns nil
func (r *Ravel) Init() error {
	return nil
}

// CreateApplication does nothing and returns nil
func (r *Ravel) CreateApplication(app models.Application) error {
	return nil
}

// DeleteApplication does nothing and returns nil
func (r *Ravel) DeleteApplication(app models.Application) error {
	return nil
}

// CreateCertificate does nothing and returns empty slice and nil
func (r *Ravel) CreateCertificate(app models.Application, cert models.Certificate) ([]models.DnsEntry, error) {
	return []models.DnsEntry{}, nil
}

// CheckDnsConfig does nothing and returns false and nil
func (r *Ravel) CheckDnsConfig(app models.Application, cert models.Certificate) (bool, error) {
	return false, nil
}

// DeleteCertificate does nothing and returns nil
func (r *Ravel) DeleteCertificate(app models.Application, cert models.Certificate) error {
	return nil
}

// IgniteBuilder does nothing and returns nil
func (r *Ravel) IgniteBuilder(app models.Application, depl models.Deployment) error {
	return nil
}

// IgniteApplication does nothing and returns nil
func (r *Ravel) IgniteApplication(app models.Application, depl models.Deployment) error {
	return nil
}

// StreamLogs does nothing and returns nil
func (r *Ravel) StreamLogs(ctx *caesar.Context, app models.Application) error {
	return nil
}

// CreateDatabase does nothing and returns nil
func (r *Ravel) CreateDatabase(db models.Database) error {
	return nil
}

// DeleteDatabase does nothing and returns nil
func (r *Ravel) DeleteDatabase(db models.Database) error {
	return nil
}

// CreateStorageBucket does nothing and returns empty strings and nil
func (r *Ravel) CreateStorageBucket(bucket models.StorageBucket) (host string, keyId string, secretKey string, region string, err error) {
	return "", "", "", "", nil
}

// GetFilesAndTotalSize does nothing and returns 0, empty slice, and nil
func (r *Ravel) GetFilesAndTotalSize(bucket models.StorageBucket) (totalSize float64, files []models.StorageFile, err error) {
	return 0, []models.StorageFile{}, nil
}

// DeleteStorageBucket does nothing and returns nil
func (r *Ravel) DeleteStorageBucket(bucket models.StorageBucket) error {
	return nil
}
