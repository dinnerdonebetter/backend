package recipestepproducts

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipe_step_products_service"
)

var _ types.RecipeStepProductDataService = (*service)(nil)

type (
	// service handles recipe step products.
	service struct {
		logger                       logging.Logger
		recipeStepProductDataManager types.RecipeStepProductDataManager
		recipeIDFetcher              func(*http.Request) string
		recipeStepIDFetcher          func(*http.Request) string
		recipeStepProductIDFetcher   func(*http.Request) string
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher         messagequeue.Publisher
		encoderDecoder               encoding.ServerEncoderDecoder
		tracer                       tracing.Tracer
	}
)

// ProvideService builds a new RecipeStepProductsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	recipeStepProductDataManager types.RecipeStepProductDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.RecipeStepProductDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe step products service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:              routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:          routeParamManager.BuildRouteParamStringIDFetcher(recipestepsservice.RecipeStepIDURIParamKey),
		recipeStepProductIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepProductIDURIParamKey),
		sessionContextDataFetcher:    authservice.FetchContextFromRequest,
		recipeStepProductDataManager: recipeStepProductDataManager,
		dataChangesPublisher:         dataChangesPublisher,
		encoderDecoder:               encoder,
		tracer:                       tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
