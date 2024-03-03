package validmeasurementunitconversions

import (
	"context"
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
	serviceName string = "valid_measurement_conversion_service"
)

var _ types.ValidMeasurementUnitConversionDataService = (*service)(nil)

type (
	// service handles valid measurement conversions.
	service struct {
		logger                                    logging.Logger
		validMeasurementUnitConversionDataManager types.ValidMeasurementUnitConversionDataManager
		validMeasurementUnitConversionIDFetcher   func(*http.Request) string
		validMeasurementUnitIDFetcher             func(*http.Request) string
		sessionContextDataFetcher                 func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                      messagequeue.Publisher
		encoderDecoder                            encoding.ServerEncoderDecoder
		tracer                                    tracing.Tracer
	}
)

// ProvideService builds a new ValidMeasurementUnitConversionsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	validMeasurementUnitConversionDataManager types.ValidMeasurementUnitConversionDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidMeasurementUnitConversionDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid measurement conversions service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                                    logging.EnsureLogger(logger).WithName(serviceName),
		validMeasurementUnitConversionIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitConversionIDURIParamKey),
		validMeasurementUnitIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		sessionContextDataFetcher:                 authservice.FetchContextFromRequest,
		validMeasurementUnitConversionDataManager: validMeasurementUnitConversionDataManager,
		dataChangesPublisher:                      dataChangesPublisher,
		encoderDecoder:                            encoder,
		tracer:                                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
