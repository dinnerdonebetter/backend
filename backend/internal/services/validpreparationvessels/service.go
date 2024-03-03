package validpreparationvessels

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparation_vessels_service"
)

var _ types.ValidPreparationVesselDataService = (*service)(nil)

type (
	// service handles valid preparation vessels.
	service struct {
		logger                            logging.Logger
		validPreparationVesselDataManager types.ValidPreparationVesselDataManager
		validPreparationVesselIDFetcher   func(*http.Request) string
		validPreparationIDFetcher         func(*http.Request) string
		validVesselIDFetcher              func(*http.Request) string
		sessionContextDataFetcher         func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher              messagequeue.Publisher
		encoderDecoder                    encoding.ServerEncoderDecoder
		tracer                            tracing.Tracer
	}
)

// ProvideService builds a new ValidPreparationVesselsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validPreparationVesselDataManager types.ValidPreparationVesselDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidPreparationVesselDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid prepraration vessels service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                            logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationVesselIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationVesselIDURIParamKey),
		validPreparationIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		validVesselIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ValidVesselIDURIParamKey),
		sessionContextDataFetcher:         authservice.FetchContextFromRequest,
		validPreparationVesselDataManager: validPreparationVesselDataManager,
		dataChangesPublisher:              dataChangesPublisher,
		encoderDecoder:                    encoder,
		tracer:                            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
