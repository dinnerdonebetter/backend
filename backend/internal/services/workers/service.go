package workers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/features/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/internal/workers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "workers_service"
)

var _ types.WorkerService = (*service)(nil)

type (
	// service handles valid vessels.
	service struct {
		cfg                            *Config
		logger                         logging.Logger
		dataManager                    database.DataManager
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher           messagequeue.Publisher
		encoderDecoder                 encoding.ServerEncoderDecoder
		tracer                         tracing.Tracer
		mealPlanFinalizationWorker     workers.MealPlanFinalizationWorker
		mealPlanGroceryListInitializer workers.MealPlanGroceryListInitializer
		mealPlanTaskCreatorWorker      workers.MealPlanTaskCreatorWorker
	}
)

// ProvideService builds a new ValidVesselsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	analyticsEventReporter analytics.EventReporter,
	grapher recipeanalysis.RecipeAnalyzer,
) (types.WorkerService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up workers service data changes publisher: %w", err)
	}

	mealPlanFinalizationWorker := workers.ProvideMealPlanFinalizationWorker(
		logger,
		dataManager,
		dataChangesPublisher,
		tracerProvider,
	)

	mealPlanGroceryListInitializer := workers.ProvideMealPlanGroceryListInitializer(
		logger,
		dataManager,
		grapher,
		dataChangesPublisher,
		analyticsEventReporter,
		tracerProvider,
		grocerylistpreparation.NewGroceryListCreator(logger, tracerProvider),
	)

	mealPlanTaskCreatorWorker := workers.ProvideMealPlanTaskCreationEnsurerWorker(
		logger,
		dataManager,
		grapher,
		dataChangesPublisher,
		analyticsEventReporter,
		tracerProvider,
	)

	svc := &service{
		cfg:                            cfg,
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		dataChangesPublisher:           dataChangesPublisher,
		encoderDecoder:                 encoder,
		dataManager:                    dataManager,
		mealPlanFinalizationWorker:     mealPlanFinalizationWorker,
		mealPlanGroceryListInitializer: mealPlanGroceryListInitializer,
		mealPlanTaskCreatorWorker:      mealPlanTaskCreatorWorker,
		tracer:                         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
