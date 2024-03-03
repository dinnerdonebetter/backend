package validpreparationinstruments

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
	serviceName string = "valid_preparation_instruments_service"
)

var _ types.ValidPreparationInstrumentDataService = (*service)(nil)

type (
	// service handles valid preparation instruments.
	service struct {
		logger                                logging.Logger
		validPreparationInstrumentDataManager types.ValidPreparationInstrumentDataManager
		validPreparationInstrumentIDFetcher   func(*http.Request) string
		validPreparationIDFetcher             func(*http.Request) string
		validInstrumentIDFetcher              func(*http.Request) string
		sessionContextDataFetcher             func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                  messagequeue.Publisher
		encoderDecoder                        encoding.ServerEncoderDecoder
		tracer                                tracing.Tracer
	}
)

// ProvideService builds a new ValidPreparationInstrumentsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validPreparationInstrumentDataManager types.ValidPreparationInstrumentDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidPreparationInstrumentDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid prepraration instruments service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                logging.EnsureLogger(logger).WithName(serviceName),
		validPreparationInstrumentIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationVesselIDURIParamKey),
		validPreparationIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		validInstrumentIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(ValidInstrumentIDURIParamKey),
		sessionContextDataFetcher:             authservice.FetchContextFromRequest,
		validPreparationInstrumentDataManager: validPreparationInstrumentDataManager,
		dataChangesPublisher:                  dataChangesPublisher,
		encoderDecoder:                        encoder,
		tracer:                                tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
