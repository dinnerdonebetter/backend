package validpreparations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparations_service"
)

var _ types.ValidPreparationDataService = (*service)(nil)

type (
	// service handles valid preparations.
	service struct {
		cfg                         *Config
		logger                      logging.Logger
		validPreparationDataManager types.ValidPreparationDataManager
		validPreparationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher   func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher        messagequeue.Publisher
		encoderDecoder              encoding.ServerEncoderDecoder
		tracer                      tracing.Tracer
		searchIndex                 search.IndexSearcher[types.ValidPreparationSearchSubset]
	}
)

// ProvideService builds a new ValidPreparationsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	validPreparationDataManager types.ValidPreparationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidPreparationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid preparations service data changes publisher: %w", err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.ValidPreparationSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeValidPreparations)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid preparation index manager")
	}

	svc := &service{
		logger:                      logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		sessionContextDataFetcher:   authservice.FetchContextFromRequest,
		validPreparationDataManager: validPreparationDataManager,
		dataChangesPublisher:        dataChangesPublisher,
		cfg:                         cfg,
		encoderDecoder:              encoder,
		tracer:                      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:                 searchIndex,
	}

	return svc, nil
}
