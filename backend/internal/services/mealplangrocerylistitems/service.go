package mealplangrocerylistitems

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "meal_plan_grocery_list_items_service"
)

var _ types.MealPlanGroceryListItemDataService = (*service)(nil)

type (
	// service handles meal plan grocery list items.
	service struct {
		logger                             logging.Logger
		mealPlanGroceryListItemDataManager types.MealPlanGroceryListItemDataManager
		mealPlanIDFetcher                  func(*http.Request) string
		mealPlanEventIDFetcher             func(*http.Request) string
		mealPlanGroceryListItemIDFetcher   func(*http.Request) string
		sessionContextDataFetcher          func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher               messagequeue.Publisher
		encoderDecoder                     encoding.ServerEncoderDecoder
		tracer                             tracing.Tracer
	}
)

// ProvideService builds a new MealPlanGroceryListItem.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	mealPlanGroceryListItemDataManager types.MealPlanGroceryListItemDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.MealPlanGroceryListItemDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up meal plans service data changes publisher: %w", err)
	}

	svc := &service{
		logger:                             logging.EnsureLogger(logger).WithName(serviceName),
		mealPlanIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(mealplansservice.MealPlanIDURIParamKey),
		mealPlanEventIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(mealplaneventsservice.MealPlanEventIDURIParamKey),
		mealPlanGroceryListItemIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(MealPlanGroceryListItemIDURIParamKey),
		sessionContextDataFetcher:          authservice.FetchContextFromRequest,
		mealPlanGroceryListItemDataManager: mealPlanGroceryListItemDataManager,
		dataChangesPublisher:               dataChangesPublisher,
		encoderDecoder:                     encoder,
		tracer:                             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
