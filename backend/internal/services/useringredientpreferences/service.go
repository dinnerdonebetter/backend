package useringredientpreferences

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

var _ types.UserIngredientPreferenceDataService = (*service)(nil)

type (
	// service handles user ingredient preferences.
	service struct {
		logger                              logging.Logger
		userIngredientPreferenceDataManager types.UserIngredientPreferenceDataManager
		userIngredientPreferenceIDFetcher   func(*http.Request) string
		sessionContextDataFetcher           func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                messagequeue.Publisher
		encoderDecoder                      encoding.ServerEncoderDecoder
		tracer                              tracing.Tracer
	}
)

// ProvideService builds a new UserIngredientPreferencesService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	userIngredientPreferenceDataManager types.UserIngredientPreferenceDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.UserIngredientPreferenceDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up user ingredient preferences service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                              logging.EnsureLogger(logger).WithName(serviceName),
		userIngredientPreferenceIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(UserIngredientPreferenceIDURIParamKey),
		sessionContextDataFetcher:           authservice.FetchContextFromRequest,
		userIngredientPreferenceDataManager: userIngredientPreferenceDataManager,
		dataChangesPublisher:                dataChangesPublisher,
		encoderDecoder:                      encoder,
		tracer:                              tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
