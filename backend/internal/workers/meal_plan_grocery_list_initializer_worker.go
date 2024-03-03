package workers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/features/grocerylistpreparation"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/hashicorp/go-multierror"
)

const (
	mealPlanGroceryListInitializerName = "meal_plan_grocery_list_initializer"
)

type (
	// MealPlanGroceryListInitializer initializes grocery lists for finalized meal plans.
	MealPlanGroceryListInitializer interface {
		InitializeGroceryListsForFinalizedMealPlans(ctx context.Context, _ []byte) error
	}

	// mealPlanGroceryListInitializer ensurers meal plan tasks are created.
	mealPlanGroceryListInitializer struct {
		logger                 logging.Logger
		tracer                 tracing.Tracer
		analyzer               recipeanalysis.RecipeAnalyzer
		encoder                encoding.ClientEncoder
		dataManager            database.DataManager
		postUpdatesPublisher   messagequeue.Publisher
		analyticsEventReporter analytics.EventReporter
		groceryListCreator     grocerylistpreparation.GroceryListCreator
	}
)

// ProvideMealPlanGroceryListInitializer provides a mealPlanGroceryListInitializer.
func ProvideMealPlanGroceryListInitializer(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	analyticsEventReporter analytics.EventReporter,
	tracerProvider tracing.TracerProvider,
	groceryListCreator grocerylistpreparation.GroceryListCreator,
) MealPlanGroceryListInitializer {
	return &mealPlanGroceryListInitializer{
		logger:                 logging.EnsureLogger(logger).WithName(mealPlanGroceryListInitializerName),
		tracer:                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(mealPlanGroceryListInitializerName)),
		encoder:                encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:            dataManager,
		analyzer:               grapher,
		postUpdatesPublisher:   postUpdatesPublisher,
		analyticsEventReporter: analyticsEventReporter,
		groceryListCreator:     groceryListCreator,
	}
}

// InitializeGroceryListsForFinalizedMealPlans handles a pending write.
func (w *mealPlanGroceryListInitializer) InitializeGroceryListsForFinalizedMealPlans(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlans, err := w.dataManager.GetFinalizedMealPlansWithUninitializedGroceryLists(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting finalized meal plan data")
	}

	logger = logger.WithValue("meal_plan_quantity", len(mealPlans))

	if len(mealPlans) > 0 {
		logger.Info("attempting to initialize grocery lists for meal plans")
	}

	var errorResult *multierror.Error

	for _, mealPlan := range mealPlans {
		l := logger.WithValue(keys.MealPlanIDKey, mealPlan.ID)

		dbInputs, groceryListCreationErr := w.groceryListCreator.GenerateGroceryListInputs(ctx, mealPlan)
		if groceryListCreationErr != nil {
			errorResult = multierror.Append(errorResult, groceryListCreationErr)
			l.Error(groceryListCreationErr, "failed to generate grocery list inputs for meal plan")
			continue
		}

		l = l.WithValue("grocery_list_items_to_create", len(dbInputs))
		l.Info("creating grocery list items for meal plan")

		for _, dbInput := range dbInputs {
			if _, err = w.dataManager.CreateMealPlanGroceryListItem(ctx, dbInput); err != nil {
				errorResult = multierror.Append(errorResult, err)
				l.Error(groceryListCreationErr, "failed to create grocery list for meal plan")
				continue
			}
		}
	}

	return errorResult.ErrorOrNil()
}
