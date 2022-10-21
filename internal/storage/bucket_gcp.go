package storage

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// GCPCloudStorageProvider indicates we'd like to use the GCP adapter for blob storage.
	GCPCloudStorageProvider = "gcp"
)

type (
	// GCPConfig configures an GCP-based storage provider.
	GCPConfig struct {
		_ struct{}

		BucketName string `json:"bucketName" mapstructure:"bucket_name" toml:"bucket_name,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*GCPConfig)(nil)

// ValidateWithContext validates the GCPCloudStorageConfig.
func (c *GCPConfig) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.BucketName, validation.Required),
	)
}