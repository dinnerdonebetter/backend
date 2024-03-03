package oauth2clients

import (
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"

	"github.com/google/wire"
)

var (
	// Providers are what we provide for dependency injection.
	Providers = wire.NewSet(
		ProvideConfig,
		ProvideOAuth2ClientsService,
	)
)

// ProvideConfig converts an auth Config to a local Config.
func ProvideConfig(cfg *authservice.Config) *Config {
	return &Config{
		DataChangesTopicName: cfg.DataChangesTopicName,
	}
}
