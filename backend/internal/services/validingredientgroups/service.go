package validingredientgroups

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
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidIngredientGroupDataService = (*service)(nil)

type (
	// service handles valid ingredients.
	service struct {
		logger                             logging.Logger
		validIngredientGroupDataManager    types.ValidIngredientGroupDataManager
		validIngredientGroupIDFetcher      func(*http.Request) string
		validIngredientGroupStateIDFetcher func(*http.Request) string
		validPreparationIDFetcher          func(*http.Request) string
		sessionContextDataFetcher          func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher               messagequeue.Publisher
		encoderDecoder                     encoding.ServerEncoderDecoder
		tracer                             tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientGroupsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	validIngredientGroupDataManager types.ValidIngredientGroupDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientGroupDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredients service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                             logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientGroupIDFetcher:      routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientGroupIDURIParamKey),
		validIngredientGroupStateIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientGroupStateIDURIParamKey),
		validPreparationIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		sessionContextDataFetcher:          authservice.FetchContextFromRequest,
		validIngredientGroupDataManager:    validIngredientGroupDataManager,
		dataChangesPublisher:               dataChangesPublisher,
		encoderDecoder:                     encoder,
		tracer:                             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
