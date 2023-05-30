package validingredientstates

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
	serviceName string = "valid_preparations_service"
)

var _ types.ValidIngredientStateDataService = (*service)(nil)

type (
	// service handles valid ingredient states.
	service struct {
		logger                          logging.Logger
		validIngredientStateDataManager types.ValidIngredientStateDataManager
		validIngredientStateIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientStatesService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	validIngredientStateDataManager types.ValidIngredientStateDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientStateDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient states service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientStateIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		validIngredientStateDataManager: validIngredientStateDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
