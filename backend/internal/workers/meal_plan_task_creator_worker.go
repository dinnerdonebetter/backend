package workers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/hashicorp/go-multierror"
)

const (
	mealPlanTaskCreationEnsurerWorkerName = "meal_plan_task_creation_ensurer"
)

type (
	// MealPlanTaskCreatorWorker ensures meal plan tasks are created.
	MealPlanTaskCreatorWorker interface {
		CreateMealPlanTasksForFinalizedMealPlans(ctx context.Context, _ []byte) error
	}

	// mealPlanTaskCreatorWorker ensurers meal plan tasks are created.
	mealPlanTaskCreatorWorker struct {
		logger                 logging.Logger
		tracer                 tracing.Tracer
		analyzer               recipeanalysis.RecipeAnalyzer
		encoder                encoding.ClientEncoder
		dataManager            database.DataManager
		postUpdatesPublisher   messagequeue.Publisher
		analyticsEventReporter analytics.EventReporter
	}
)

// ProvideMealPlanTaskCreationEnsurerWorker provides a mealPlanTaskCreatorWorker.
func ProvideMealPlanTaskCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	grapher recipeanalysis.RecipeAnalyzer,
	postUpdatesPublisher messagequeue.Publisher,
	analyticsEventReporter analytics.EventReporter,
	tracerProvider tracing.TracerProvider,
) *mealPlanTaskCreatorWorker {
	return &mealPlanTaskCreatorWorker{
		logger:                 logging.EnsureLogger(logger).WithName(mealPlanTaskCreationEnsurerWorkerName),
		tracer:                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(mealPlanTaskCreationEnsurerWorkerName)),
		encoder:                encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:            dataManager,
		analyzer:               grapher,
		postUpdatesPublisher:   postUpdatesPublisher,
		analyticsEventReporter: analyticsEventReporter,
	}
}

// CreateMealPlanTasksForFinalizedMealPlans does the main thing.
func (w *mealPlanTaskCreatorWorker) CreateMealPlanTasksForFinalizedMealPlans(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	mealPlansAndSteps, err := w.determineCreatableMealPlanTasks(ctx)
	if err != nil {
		return observability.PrepareError(err, nil, "determining creatable steps")
	}

	logger = logger.WithValue("creatable_steps_qty", len(mealPlansAndSteps))

	var result *multierror.Error
	for mealPlanID, steps := range mealPlansAndSteps {
		l := logger.Clone().WithValue(keys.MealPlanIDKey, mealPlanID).WithValue("creatable_prep_step_qty", len(steps))

		createdMealPlanTasks, creationErr := w.dataManager.CreateMealPlanTasksForMealPlanOption(ctx, steps)
		if creationErr != nil {
			result = multierror.Append(result, creationErr)
			observability.AcknowledgeError(creationErr, l, span, "creating meal plan tasks for meal plan option")
			continue
		}

		for _, createdStep := range createdMealPlanTasks {
			if publishErr := w.postUpdatesPublisher.Publish(ctx, &types.DataChangeMessage{
				EventType:      types.MealPlanTaskCreatedCustomerEventType,
				MealPlanTask:   createdStep,
				MealPlanTaskID: createdStep.ID,
				MealPlanID:     mealPlanID,
			}); publishErr != nil {
				observability.AcknowledgeError(publishErr, l, span, "publishing data change event")
			}
		}

		if err = w.dataManager.MarkMealPlanAsHavingTasksCreated(ctx, mealPlanID); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

// determineCreatableMealPlanTasks determines which meal plan tasks are creatable for a recipe.
func (w *mealPlanTaskCreatorWorker) determineCreatableMealPlanTasks(ctx context.Context) (map[string][]*types.MealPlanTaskDatabaseCreationInput, error) {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()
	logger.Info("fetching finalized meal plan IDs to determine creatable steps")

	results, err := w.dataManager.GetFinalizedMealPlanIDsForTheNextWeek(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting finalized meal plan data for the next week")
	}

	if len(results) > 0 {
		logger = logger.WithValue("steps_to_create", len(results))
		logger.Info("determining creatable steps")
	}

	inputs := map[string][]*types.MealPlanTaskDatabaseCreationInput{}
	for _, result := range results {
		l := logger.Clone().WithValues(map[string]any{
			keys.MealPlanIDKey:       result.MealPlanID,
			keys.MealPlanEventIDKey:  result.MealPlanEventID,
			keys.MealPlanOptionIDKey: result.MealPlanOptionID,
			keys.MealIDKey:           result.MealID,
			"recipe_ids":             result.RecipeIDs,
		})
		l.Info("fetching meal plan event")

		if _, ok := inputs[result.MealPlanID]; !ok {
			inputs[result.MealPlanID] = []*types.MealPlanTaskDatabaseCreationInput{}
		}

		for _, recipeID := range result.RecipeIDs {
			recipe, getRecipeErr := w.dataManager.GetRecipe(ctx, recipeID)
			if getRecipeErr != nil {
				return nil, observability.PrepareAndLogError(getRecipeErr, l, span, "fetching recipe")
			}

			creatableSteps, determineStepsErr := w.analyzer.GenerateMealPlanTasksForRecipe(ctx, result.MealPlanOptionID, recipe)
			if determineStepsErr != nil {
				return nil, observability.PrepareAndLogError(determineStepsErr, l, span, "fetching recipe")
			}

			inputs[result.MealPlanID] = append(inputs[result.MealPlanID], creatableSteps...)
		}
	}

	return inputs, nil
}
