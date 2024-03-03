package validvessels

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
	serviceName string = "valid_vessels_service"
)

var _ types.ValidVesselDataService = (*service)(nil)

type (
	// service handles valid vessels.
	service struct {
		cfg                       *Config
		logger                    logging.Logger
		validVesselDataManager    types.ValidVesselDataManager
		validVesselIDFetcher      func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
		searchIndex               search.IndexSearcher[types.ValidVesselSearchSubset]
	}
)

// ProvideService builds a new ValidVesselsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	validVesselDataManager types.ValidVesselDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidVesselDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid vessels service data changes publisher: %w", err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.ValidVesselSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeValidVessels)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid vessel index manager")
	}

	svc := &service{
		cfg:                       cfg,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		validVesselIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(ValidVesselIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		validVesselDataManager:    validVesselDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:               searchIndex,
	}

	return svc, nil
}
