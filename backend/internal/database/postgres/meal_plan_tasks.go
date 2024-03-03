package postgres

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealPlanTaskDataManager = (*Querier)(nil)
)

// MealPlanTaskExists checks if a meal plan task exists.
func (q *Querier) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanTaskID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	result, err := q.generatedQuerier.CheckMealPlanTaskExistence(ctx, q.db, &generated.CheckMealPlanTaskExistenceParams{
		MealPlanID:     mealPlanID,
		MealPlanTaskID: mealPlanTaskID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan task existence check")
	}

	logger.Info("meal plan task existence retrieved")

	return result, nil
}

// GetMealPlanTask fetches a meal plan task.
func (q *Querier) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanTaskID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, mealPlanTaskID)
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)

	result, err := q.generatedQuerier.GetMealPlanTask(ctx, q.db, mealPlanTaskID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan task existence check")
	}

	mealPlanTask := &types.MealPlanTask{
		RecipePrepTask:      types.RecipePrepTask{},
		CreatedAt:           result.CreatedAt,
		LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
		CompletedAt:         database.TimePointerFromNullTime(result.CompletedAt),
		AssignedToUser:      database.StringPointerFromNullString(result.AssignedToUser),
		ID:                  result.ID,
		Status:              string(result.Status),
		CreationExplanation: result.CreationExplanation,
		StatusExplanation:   result.StatusExplanation,
		MealPlanOption: types.MealPlanOption{
			CreatedAt:              result.MealPlanOptionCreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.MealPlanOptionLastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.MealPlanOptionAssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.MealPlanOptionArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.MealPlanOptionAssignedDishwasher),
			Notes:                  result.MealPlanOptionNotes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.MealPlanOptionBelongsToMealPlanEvent),
			ID:                     result.MealPlanOptionID,
			Votes:                  nil,
			Meal: types.Meal{
				ID: result.MealPlanOptionMealID,
			},
			MealScale: database.Float32FromString(result.MealPlanOptionMealScale),
			Chosen:    result.MealPlanOptionChosen,
			TieBroken: result.MealPlanOptionTiebroken,
		},
	}

	return mealPlanTask, nil
}

// createMealPlanTask creates a meal plan task.
func (q *Querier) createMealPlanTask(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	// create the meal plan task.
	if err := q.generatedQuerier.CreateMealPlanTask(ctx, querier, &generated.CreateMealPlanTaskParams{
		ID:                      input.ID,
		Status:                  types.MealPlanTaskStatusUnfinished,
		StatusExplanation:       input.StatusExplanation,
		CreationExplanation:     input.CreationExplanation,
		BelongsToMealPlanOption: input.MealPlanOptionID,
		BelongsToRecipePrepTask: input.RecipePrepTaskID,
		AssignedToUser:          database.NullStringFromStringPointer(input.AssignedToUser),
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	x := &types.MealPlanTask{
		CreatedAt:           q.timeFunc(),
		ID:                  input.ID,
		AssignedToUser:      input.AssignedToUser,
		Status:              types.MealPlanTaskStatusUnfinished,
		CreationExplanation: input.CreationExplanation,
		StatusExplanation:   input.StatusExplanation,
		MealPlanOption: types.MealPlanOption{
			ID: input.MealPlanOptionID,
		},
		RecipePrepTask: types.RecipePrepTask{
			ID: input.RecipePrepTaskID,
		},
	}

	tracing.AttachToSpan(span, keys.MealPlanIDKey, x.ID)
	logger.Info("meal plan task created")

	return x, nil
}

// CreateMealPlanTask creates a meal plan task.
func (q *Querier) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	x, err := q.createMealPlanTask(ctx, tx, input)
	if err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan task")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("meal plan task created")

	return x, nil
}

// GetMealPlanTasksForMealPlan fetches a list of meal plan tasks.
func (q *Querier) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanTask, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := q.generatedQuerier.ListAllMealPlanTasksByMealPlan(ctx, q.db, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan tasks list retrieval query")
	}

	x = []*types.MealPlanTask{}
	for _, result := range results {
		mealPlanTask := &types.MealPlanTask{
			RecipePrepTask:      types.RecipePrepTask{},
			CreatedAt:           result.CreatedAt,
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			CompletedAt:         database.TimePointerFromNullTime(result.CompletedAt),
			AssignedToUser:      database.StringPointerFromNullString(result.AssignedToUser),
			ID:                  result.ID,
			Status:              string(result.Status),
			CreationExplanation: result.CreationExplanation,
			StatusExplanation:   result.StatusExplanation,
			MealPlanOption: types.MealPlanOption{
				CreatedAt:              result.MealPlanOptionCreatedAt,
				LastUpdatedAt:          database.TimePointerFromNullTime(result.MealPlanOptionLastUpdatedAt),
				AssignedCook:           database.StringPointerFromNullString(result.MealPlanOptionAssignedCook),
				ArchivedAt:             database.TimePointerFromNullTime(result.MealPlanOptionArchivedAt),
				AssignedDishwasher:     database.StringPointerFromNullString(result.MealPlanOptionAssignedDishwasher),
				Notes:                  result.MealPlanOptionNotes,
				BelongsToMealPlanEvent: database.StringFromNullString(result.MealPlanOptionBelongsToMealPlanEvent),
				ID:                     result.MealPlanOptionID,
				Votes:                  nil,
				Meal: types.Meal{
					ID: result.MealPlanOptionMealID,
				},
				MealScale: database.Float32FromString(result.MealPlanOptionMealScale),
				Chosen:    result.MealPlanOptionChosen,
				TieBroken: result.MealPlanOptionTiebroken,
			},
		}

		x = append(x, mealPlanTask)
	}

	logger.Info("meal plan tasks retrieved")

	return x, nil
}

// CreateMealPlanTasksForMealPlanOption creates meal plan tasks.
func (q *Querier) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	outputs := []*types.MealPlanTask{}
	for _, input := range inputs {
		mealPlanTask, createMealPlanTaskErr := q.createMealPlanTask(ctx, tx, input)
		if createMealPlanTaskErr != nil {
			q.rollbackTransaction(ctx, tx)
			return nil, observability.PrepareAndLogError(createMealPlanTaskErr, logger, span, "creating meal plan task")
		}

		outputs = append(outputs, mealPlanTask)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, observability.PrepareAndLogError(commitErr, logger, span, "committing transaction")
	}

	logger.Info("meal plan tasks created")

	return outputs, nil
}

// MarkMealPlanAsHavingTasksCreated marks a meal plan as having all its tasks created.
func (q *Querier) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsPrepTasksCreated(ctx, q.db, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	logger.Info("meal plan tasks created")

	return nil
}

// MarkMealPlanAsHavingGroceryListInitialized marks a meal plan as having all its tasks created.
func (q *Querier) MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, mealPlanID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if err := q.generatedQuerier.MarkMealPlanAsGroceryListInitialized(ctx, q.db, mealPlanID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking meal plan as having tasks created")
	}

	logger.Info("meal plan marked as grocery list initialized")

	return nil
}

// ChangeMealPlanTaskStatus changes a meal plan task's status.
func (q *Querier) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanTaskIDKey, input.ID)
	logger = logger.WithValue(keys.MealPlanTaskIDKey, input.ID)

	var settledAt *time.Time
	if input.Status != nil && *input.Status == types.MealPlanTaskStatusFinished {
		t := q.timeFunc()
		settledAt = &t
	}

	var newStatus string
	if input.Status != nil {
		newStatus = *input.Status
	}

	if err := q.generatedQuerier.ChangeMealPlanTaskStatus(ctx, q.db, &generated.ChangeMealPlanTaskStatusParams{
		ID:                input.ID,
		Status:            generated.PrepStepStatus(newStatus),
		StatusExplanation: input.StatusExplanation,
		CompletedAt:       database.NullTimeFromTimePointer(settledAt),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "changing meal plan task status")
	}

	logger.Info("meal plan task status changed")

	return nil
}
