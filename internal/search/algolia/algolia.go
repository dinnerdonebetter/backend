package algolia

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/pkg/types"

	algolia "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

var (
	_ search.Index[types.ValidIngredientSearchSubset]      = (*indexManager[types.ValidIngredientSearchSubset])(nil)
	_ search.Index[types.ValidInstrumentSearchSubset]      = (*indexManager[types.ValidInstrumentSearchSubset])(nil)
	_ search.Index[types.ValidMeasurementUnitSearchSubset] = (*indexManager[types.ValidMeasurementUnitSearchSubset])(nil)
	_ search.Index[types.ValidPreparationSearchSubset]     = (*indexManager[types.ValidPreparationSearchSubset])(nil)
	_ search.Index[types.ValidIngredientStateSearchSubset] = (*indexManager[types.ValidIngredientStateSearchSubset])(nil)
)

type (
	indexManager[T search.Searchable] struct {
		logger   logging.Logger
		tracer   tracing.Tracer
		client   *algolia.Index
		DataType *T
		timeout  time.Duration
	}
)

func ProvideIndexManager[T search.Searchable](
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *Config,
	indexName string,
) (search.Index[T], error) {
	im := &indexManager[T]{
		tracer:  tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("search_%s", indexName))),
		logger:  logging.EnsureLogger(logger).WithName(indexName),
		client:  algolia.NewClient(cfg.AppID, cfg.APIKey).InitIndex(indexName),
		timeout: cfg.Timeout,
	}

	return im, nil
}
