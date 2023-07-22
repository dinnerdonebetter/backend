package jaeger

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// Jaeger represents the open source tracing server.
	Jaeger = "jaeger"
)

type (
	// Config contains settings related to tracing.
	Config struct {
		_ struct{}

		CollectorEndpoint         string  `json:"collector_endpoint,omitempty"        toml:"collector_endpoint,omitempty"`
		ServiceName               string  `json:"service_name,omitempty"              toml:"service_name,omitempty"`
		SpanCollectionProbability float64 `json:"spanCollectionProbability,omitempty" toml:"span_collection_probability,omitempty"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the config struct.
func (c *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, c,
		validation.Field(&c.CollectorEndpoint, validation.Required),
		validation.Field(&c.ServiceName, validation.Required),
	)
}
