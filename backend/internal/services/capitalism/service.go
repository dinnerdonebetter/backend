package capitalism

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/capitalism"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "capitalism_service"
)

var _ types.CapitalismService = (*service)(nil)

type (
	// service handles valid instruments.
	service struct {
		cfg                       *Config
		logger                    logging.Logger
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher      messagequeue.Publisher
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
		paymentManager            capitalism.PaymentManager
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	encoder encoding.ServerEncoderDecoder,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	paymentManager capitalism.PaymentManager,
) (types.CapitalismService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid instruments service data changes publisher: %w", err)
	}

	svc := &service{
		cfg:                       cfg,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		paymentManager:            paymentManager,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
