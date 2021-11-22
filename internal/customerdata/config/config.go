package config

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/customerdata/segment"
	"github.com/prixfixeco/api_server/internal/observability/logging"
)

const (
	providerSegment = "segment"
)

type (
	// Config is the configuration structure.
	Config struct {
		Provider string
		APIToken string
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.APIToken, validation.When(strings.EqualFold(strings.TrimSpace(cfg.Provider), providerSegment), validation.Required)),
	)
}

// ProvideCollector provides a collector.
func (cfg *Config) ProvideCollector(logger logging.Logger) (customerdata.Collector, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case providerSegment:
		return segment.NewSegmentCustomerDataCollector(logger, cfg.APIToken)
	default:
		return customerdata.NewNoopCollector()
	}
}
