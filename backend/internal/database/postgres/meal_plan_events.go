package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.MealPlanEventDataManager = (*Querier)(nil)
)

// MealPlanEventExists fetches whether a meal plan event exists from the database.
func (q *Querier) MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	result, err := q.generatedQuerier.CheckMealPlanEventExistence(ctx, q.db, &generated.CheckMealPlanEventExistenceParams{
		ID:         mealPlanEventID,
		MealPlanID: mealPlanID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

// GetMealPlanEvent fetches a meal plan event from the database.
func (q *Querier) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	result, err := q.generatedQuerier.GetMealPlanEvent(ctx, q.db, &generated.GetMealPlanEventParams{
		ID:                mealPlanEventID,
		BelongsToMealPlan: mealPlanID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan event retrieval query")
	}

	mealPlanEvent := &types.MealPlanEvent{
		CreatedAt:         result.CreatedAt,
		StartsAt:          result.StartsAt,
		EndsAt:            result.EndsAt,
		ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
		MealName:          string(result.MealName),
		Notes:             result.Notes,
		BelongsToMealPlan: result.BelongsToMealPlan,
		ID:                result.ID,
	}

	options, err := q.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting meal plan options for meal plan event")
	}
	mealPlanEvent.Options = options

	return mealPlanEvent, nil
}

// getMealPlanEventsForMealPlan fetches a list of mealPlanEvents from the database that meet a particular filter.
func (q *Querier) getMealPlanEventsForMealPlan(ctx context.Context, mealPlanID string) (x []*types.MealPlanEvent, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	x = []*types.MealPlanEvent{}

	results, err := q.generatedQuerier.GetAllMealPlanEventsForMealPlan(ctx, q.db, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	for _, result := range results {
		event := &types.MealPlanEvent{
			CreatedAt:         result.CreatedAt,
			StartsAt:          result.StartsAt,
			EndsAt:            result.EndsAt,
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			MealName:          string(result.MealName),
			Notes:             result.Notes,
			BelongsToMealPlan: result.BelongsToMealPlan,
			ID:                result.ID,
		}

		mealPlanOptions, mealPlanOptionsErr := q.getMealPlanOptionsForMealPlanEvent(ctx, mealPlanID, event.ID)
		if mealPlanOptionsErr != nil {
			return nil, observability.PrepareAndLogError(mealPlanOptionsErr, logger, span, "fetching options for meal plan events")
		}

		event.Options = mealPlanOptions

		x = append(x, event)
	}

	return x, nil
}

// GetMealPlanEvents fetches a list of meal plan events from the database that meet a particular filter.
func (q *Querier) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.MealPlanEvent], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.MealPlanEvent]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetMealPlanEvents(ctx, q.db, &generated.GetMealPlanEventsParams{
		MealPlanID:    mealPlanID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan events list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.MealPlanEvent{
			CreatedAt:         result.CreatedAt,
			StartsAt:          result.StartsAt,
			EndsAt:            result.EndsAt,
			ArchivedAt:        database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:     database.TimePointerFromNullTime(result.LastUpdatedAt),
			MealName:          string(result.MealName),
			Notes:             result.Notes,
			BelongsToMealPlan: result.BelongsToMealPlan,
			ID:                result.ID,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// MealPlanEventIsEligibleForVoting returns if a meal plan can be voted on.
func (q *Querier) MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	result, err := q.generatedQuerier.MealPlanEventIsEligibleForVoting(ctx, q.db, &generated.MealPlanEventIsEligibleForVotingParams{
		MealPlanID:      mealPlanID,
		MealPlanEventID: mealPlanEventID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan event existence check")
	}

	return result, nil
}

// createMealPlanEvent creates a meal plan event in the database.
func (q *Querier) createMealPlanEvent(ctx context.Context, querier database.SQLQueryExecutorAndTransactionManager, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanEventIDKey, input.ID)

	// create the meal plan event.
	if err := q.generatedQuerier.CreateMealPlanEvent(ctx, querier, &generated.CreateMealPlanEventParams{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          generated.MealName(input.MealName),
		BelongsToMealPlan: input.BelongsToMealPlan,
	}); err != nil {
		q.rollbackTransaction(ctx, querier)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan event creation query")
	}

	x := &types.MealPlanEvent{
		ID:                input.ID,
		Notes:             input.Notes,
		StartsAt:          input.StartsAt,
		EndsAt:            input.EndsAt,
		MealName:          input.MealName,
		BelongsToMealPlan: input.BelongsToMealPlan,
		CreatedAt:         q.currentTime(),
	}

	logger.WithValue("quantity", len(input.Options)).Info("creating options for meal plan event")
	for _, option := range input.Options {
		option.BelongsToMealPlanEvent = x.ID
		opt, createErr := q.createMealPlanOption(ctx, querier, option, len(input.Options) == 1)
		if createErr != nil {
			q.rollbackTransaction(ctx, querier)
			return nil, observability.PrepareError(createErr, span, "creating meal plan option for meal plan event")
		}
		x.Options = append(x.Options, opt)
	}

	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, x.ID)
	logger.Info("meal plan event created")

	return x, nil
}

// CreateMealPlanEvent creates a meal plan event in the database.
func (q *Querier) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, span, "beginning transaction")
	}

	x, err := q.createMealPlanEvent(ctx, tx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan event")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, span, "committing transaction")
	}

	return x, nil
}

// UpdateMealPlanEvent updates a particular meal plan event.
func (q *Querier) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanEventIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateMealPlanEvent(ctx, q.db, &generated.UpdateMealPlanEventParams{
		Notes:             updated.Notes,
		StartsAt:          updated.StartsAt,
		EndsAt:            updated.EndsAt,
		MealName:          generated.MealName(updated.MealName),
		BelongsToMealPlan: updated.BelongsToMealPlan,
		ID:                updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan event")
	}

	logger.Info("meal plan event updated")

	return nil
}

// ArchiveMealPlanEvent archives a meal plan event from the database by its ID.
func (q *Querier) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanEventIDKey, mealPlanEventID)
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	if _, err := q.generatedQuerier.ArchiveMealPlanEvent(ctx, q.db, &generated.ArchiveMealPlanEventParams{
		ID:                mealPlanEventID,
		BelongsToMealPlan: mealPlanID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan event")
	}

	logger.Info("meal plan event archived")

	return nil
}
