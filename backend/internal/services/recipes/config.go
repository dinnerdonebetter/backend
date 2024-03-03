package recipes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/uploads"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config configures the service.
type Config struct {
	_ struct{} `json:"-"`

	DataChangesTopicName string         `json:"dataChangesTopicName,omitempty" toml:"data_changes_topic_name,omitempty"`
	PublicMediaURLPrefix string         `json:"mediaUploadPrefix"              toml:"media_upload_prefix"`
	Uploads              uploads.Config `json:"uploads"                        toml:"uploads,omitempty"`
	UseSearchService     bool           `json:"searchFromDatabase"             toml:"search_from_database,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.PublicMediaURLPrefix, validation.Required),
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
		validation.Field(&cfg.Uploads, validation.Required),
	)
}
