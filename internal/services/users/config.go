package users

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/uploads"
)

type (
	// Config configures the users service.
	Config struct {
		DataChangesTopicName string         `json:"dataChangesTopicName,omitempty" mapstructure:"data_changes_topic_name" toml:"data_changes_topic_name,omitempty"`
		Uploads              uploads.Config `json:"uploads" mapstructure:"uploads" toml:"uploads,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.DataChangesTopicName, validation.Required),
		validation.Field(&cfg.Uploads, validation.Required),
	)
}
