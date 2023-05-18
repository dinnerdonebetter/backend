package config

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	"github.com/dinnerdonebetter/backend/internal/search/elasticsearch"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ElasticsearchProvider represents the elasticsearch search index provider.
	ElasticsearchProvider = "elasticsearch"
	// AlgoliaProvider represents the algolia search index provider.
	AlgoliaProvider = "algolia"
)

// Config contains settings regarding search indices.
type Config struct {
	_ struct{}

	Algolia       *algolia.Config       `json:"algolia" mapstructure:"algolia" toml:"algolia,omitempty"`
	Elasticsearch *elasticsearch.Config `json:"elasticsearch" mapstructure:"elasticsearch" toml:"elasticsearch,omitempty"`
	Provider      string                `json:"provider" mapstructure:"provider" toml:"provider,omitempty"`
}

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates a Config struct.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, cfg,
		validation.Field(&cfg.Provider, validation.In(ElasticsearchProvider, AlgoliaProvider)),
		validation.Field(&cfg.Algolia, validation.When(cfg.Provider == AlgoliaProvider, validation.Required)),
		validation.Field(&cfg.Elasticsearch, validation.When(cfg.Provider == ElasticsearchProvider, validation.Required)),
	)
}

// ProvideIndexManager validates a Config struct.
func ProvideIndexManager[T search.Searchable](ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config, indexName string) (search.IndexManager[T], error) {
	switch cfg.Provider {
	case ElasticsearchProvider:
		return elasticsearch.ProvideIndexManager[T](ctx, logger, tracerProvider, cfg.Elasticsearch, indexName)
	case AlgoliaProvider:
		return algolia.ProvideIndexManager[T](ctx, logger, tracerProvider, cfg.Algolia, indexName)
	default:
		return nil, fmt.Errorf("invalid search provider: %q", cfg.Provider)
	}
}
