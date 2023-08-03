package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	_ search.Index[types.UserSearchSubset] = (*indexManager[types.UserSearchSubset])(nil)
)

type (
	indexManager[T search.Searchable] struct {
		logger                logging.Logger
		tracer                tracing.Tracer
		esclient              *elasticsearch.Client
		indexName             string
		indexOperationTimeout time.Duration
	}
)

func ProvideIndexManager[T search.Searchable](ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config, indexName string) (search.Index[T], error) {
	c, err := cfg.provideElasticsearchClient()
	if err != nil {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	if ready := elasticsearchIsReady(ctx, cfg, logger, 10); !ready {
		return nil, fmt.Errorf("initializing search client: %w", err)
	}

	im := &indexManager[T]{
		tracer:                tracing.NewTracer(tracerProvider.Tracer(fmt.Sprintf("search_%s", indexName))),
		logger:                logging.EnsureLogger(logger).WithName(indexName),
		esclient:              c,
		indexOperationTimeout: cfg.IndexOperationTimeout,
		indexName:             indexName,
	}

	if indexErr := im.ensureIndices(ctx); indexErr != nil {
		return nil, indexErr
	}

	return im, nil
}

func elasticsearchIsReady(
	ctx context.Context,
	cfg *Config,
	l logging.Logger,
	maxAttempts uint8,
) (ready bool) {
	attemptCount := 0

	logger := l.WithValues(map[string]any{
		"interval": time.Second.String(),
		"address":  cfg.Address,
	})

	logger.Debug("checking if elasticsearch is ready")

	for !ready {
		c, err := cfg.provideElasticsearchClient()
		if err != nil {
			logger.WithValue("attempt_count", attemptCount).Debug("client setup failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		}

		if res, infoReqErr := (esapi.InfoRequest{}).Do(ctx, c); infoReqErr != nil && !res.IsError() {
			logger.WithValue("attempt_count", attemptCount).Debug("ping failed, waiting for elasticsearch")
			time.Sleep(time.Second)

			attemptCount++
			if attemptCount >= int(maxAttempts) {
				break
			}
		} else {
			ready = true
			logger.Debug("elasticsearch is ready")
			return ready
		}
	}

	logger.Debug("elasticsearch is ready")

	return false
}

func (sm *indexManager[T]) ensureIndices(ctx context.Context) error {
	_, span := sm.tracer.StartSpan(ctx)
	defer span.End()

	res, err := esapi.IndicesExistsRequest{
		Index:             []string{sm.indexName},
		IgnoreUnavailable: esapi.BoolPtr(false),
		ErrorTrace:        false,
		FilterPath:        nil,
	}.Do(ctx, sm.esclient)
	if err != nil {
		return observability.PrepareError(err, span, "checking index existence successfully")
	}

	if res.StatusCode == http.StatusNotFound {
		if _, err = (esapi.IndicesCreateRequest{Index: strings.ToLower(sm.indexName)}).Do(ctx, sm.esclient); err != nil {
			return observability.PrepareError(err, span, "checking index existence")
		}
	}

	return nil
}
