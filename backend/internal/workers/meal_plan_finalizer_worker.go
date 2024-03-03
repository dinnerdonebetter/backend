package workers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

type (
	MealPlanFinalizationWorker interface {
		FinalizeExpiredMealPlans(ctx context.Context, _ []byte) (int, error)
		FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context, _ []byte) error
	}

	// mealPlanFinalizationWorker finalizes meal plans.
	mealPlanFinalizationWorker struct {
		logger               logging.Logger
		tracer               tracing.Tracer
		encoder              encoding.ClientEncoder
		dataManager          database.DataManager
		postUpdatesPublisher messagequeue.Publisher
	}
)

// ProvideMealPlanFinalizationWorker provides a mealPlanFinalizationWorker.
func ProvideMealPlanFinalizationWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	tracerProvider tracing.TracerProvider,
) MealPlanFinalizationWorker {
	n := "meal_plan_finalizer"

	return &mealPlanFinalizationWorker{
		logger:               logging.EnsureLogger(logger).WithName(n),
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(n)),
		encoder:              encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:          dataManager,
		postUpdatesPublisher: postUpdatesPublisher,
	}
}

// finalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) finalizeExpiredMealPlans(ctx context.Context) (int, error) {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()
	logger.Info("beginning finalization of expired meal plans")

	mealPlans, err := w.dataManager.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
	if err != nil {
		return -1, observability.PrepareAndLogError(err, logger, span, "fetching unfinalized and expired meal plan")
	}

	if len(mealPlans) > 0 {
		logger.WithValue("quantity", len(mealPlans)).Info("finalizing expired meal plans")
	}

	var changedCount int
	for _, mealPlan := range mealPlans {
		var changed bool
		changed, err = w.dataManager.AttemptToFinalizeMealPlan(ctx, mealPlan.ID, mealPlan.BelongsToHousehold)
		if err != nil {
			return -1, observability.PrepareError(err, span, "finalizing meal plan")
		}

		if changed {
			changedCount++
		}
	}

	logger.WithValue("changed_count", changedCount).Info("finalized expired meal plans")

	return changedCount, nil
}

// FinalizeExpiredMealPlans handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) FinalizeExpiredMealPlans(ctx context.Context, _ []byte) (int, error) {
	return w.finalizeExpiredMealPlans(ctx)
}

// FinalizeExpiredMealPlansWithoutReturningCount handles a message ordering the finalization of expired meal plans.
func (w *mealPlanFinalizationWorker) FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context, _ []byte) error {
	_, err := w.finalizeExpiredMealPlans(ctx)
	return err
}
