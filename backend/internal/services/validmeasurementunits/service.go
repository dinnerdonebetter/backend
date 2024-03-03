package validmeasurementunits

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
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidMeasurementUnitDataService = (*service)(nil)

type (
	// service handles valid ingredients.
	service struct {
		cfg                             *Config
		logger                          logging.Logger
		validMeasurementUnitDataManager types.ValidMeasurementUnitDataManager
		validMeasurementUnitIDFetcher   func(*http.Request) string
		validIngredientIDFetcher        func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
		searchIndex                     search.IndexSearcher[types.ValidMeasurementUnitSearchSubset]
	}
)

// ProvideService builds a new ValidMeasurementUnitsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	validMeasurementUnitDataManager types.ValidMeasurementUnitDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidMeasurementUnitDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients service data changes publisher: %w", err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeValidMeasurementUnits)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid preparation index manager")
	}

	svc := &service{
		cfg:                             cfg,
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		validMeasurementUnitIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		validIngredientIDFetcher:        routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		validMeasurementUnitDataManager: validMeasurementUnitDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:                     searchIndex,
	}

	return svc, nil
}
