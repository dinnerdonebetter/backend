package mealplanfinalizerfunction

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"

	"github.com/prixfixeco/backend/internal/config"
	customerdataconfig "github.com/prixfixeco/backend/internal/customerdata/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/features/grocerylistpreparation"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("InitializeGroceryListsItemsForMealPlans", InitializeGroceryListsItemsForMealPlans)
}

// InitializeGroceryListsItemsForMealPlans is our cloud function entrypoint.
func InitializeGroceryListsItemsForMealPlans(ctx context.Context, _ event.Event) error {
	logger := zerolog.NewZerologLogger()
	logger.SetLevel(logging.DebugLevel)

	cfg, err := config.GetMealPlanGroceryListInitializerWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.Initialize(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	ctx, span := tracing.NewTracer(tracerProvider.Tracer("meal_plan_grocery_list_items_init_job")).StartSpan(ctx)
	defer span.End()

	customerDataCollector, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring customer data collector")
	}

	defer customerDataCollector.Close()

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	defer dataManager.Close()

	if !dataManager.IsReady(ctx, 50) {
		return observability.PrepareAndLogError(errors.New("database is not ready"), logger, span, "pinging database")
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	dataChangesPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring data changes publisher")
	}

	defer dataChangesPublisher.Stop()

	mealPlanGroceryListInitializationWorker := workers.ProvideMealPlanGroceryListInitializer(
		logger,
		dataManager,
		recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider),
		dataChangesPublisher,
		customerDataCollector,
		tracerProvider,
		grocerylistpreparation.NewGroceryListCreator(logger, tracerProvider),
	)

	if err = mealPlanGroceryListInitializationWorker.HandleMessage(ctx, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "handling message")
	}

	return nil
}
