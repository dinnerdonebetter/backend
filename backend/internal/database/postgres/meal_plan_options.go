package postgres

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"math/rand/v2"
	"resenje.org/schulze"
)

var (
	_ types.MealPlanOptionDataManager = (*Querier)(nil)
)

// MealPlanOptionExists fetches whether a meal plan option exists from the database.
func (q *Querier) MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (exists bool, err error) {
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

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := q.generatedQuerier.CheckMealPlanOptionExistence(ctx, q.db, &generated.CheckMealPlanOptionExistenceParams{
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
		MealPlanOptionID: mealPlanOptionID,
		MealPlanID:       mealPlanID,
	})
	if err != nil {
		logger.Error(err, "performing meal plan option existence check")
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan option existence check")
	}

	return result, nil
}

// GetMealPlanOption fetches a meal plan option from the database.
func (q *Querier) GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*types.MealPlanOption, error) {
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

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := q.generatedQuerier.GetMealPlanOption(ctx, q.db, &generated.GetMealPlanOptionParams{
		MealPlanID:       mealPlanID,
		MealPlanEventID:  database.NullStringFromString(mealPlanEventID),
		MealPlanOptionID: mealPlanOptionID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	mealPlanOption := &types.MealPlanOption{
		CreatedAt:              result.CreatedAt,
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
		Notes:                  result.Notes,
		BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
		ID:                     result.ID,
		Meal: types.Meal{
			CreatedAt:                result.MealCreatedAt,
			ArchivedAt:               database.TimePointerFromNullTime(result.MealArchivedAt),
			LastUpdatedAt:            database.TimePointerFromNullTime(result.MealLastUpdatedAt),
			MaximumEstimatedPortions: database.Float32PointerFromNullString(result.MealMaxEstimatedPortions),
			ID:                       result.MealID,
			Description:              result.MealDescription,
			CreatedByUser:            result.MealCreatedByUser,
			Name:                     result.MealName,
			MinimumEstimatedPortions: database.Float32FromString(result.MealMinEstimatedPortions),
			EligibleForMealPlans:     result.MealEligibleForMealPlans,
		},
		MealScale: database.Float32FromString(result.MealScale),
		Chosen:    result.Chosen,
		TieBroken: result.Tiebroken,
	}

	return mealPlanOption, nil
}

// getMealPlanOptionByID fetches a meal plan option from the database by its ID.
func (q *Querier) getMealPlanOptionByID(ctx context.Context, mealPlanOptionID string) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanOptionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	result, err := q.generatedQuerier.GetMealPlanOptionByID(ctx, q.db, mealPlanOptionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	mealPlanOption := &types.MealPlanOption{
		CreatedAt:              result.CreatedAt,
		LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
		AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
		ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
		AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
		Notes:                  result.Notes,
		BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
		ID:                     result.ID,
		Votes:                  nil,
		Meal: types.Meal{
			CreatedAt:                result.MealCreatedAt,
			ArchivedAt:               database.TimePointerFromNullTime(result.MealArchivedAt),
			LastUpdatedAt:            database.TimePointerFromNullTime(result.MealLastUpdatedAt),
			MaximumEstimatedPortions: database.Float32PointerFromNullString(result.MealMaxEstimatedPortions),
			ID:                       result.MealID,
			Description:              result.MealDescription,
			CreatedByUser:            result.MealCreatedByUser,
			Name:                     result.MealName,
			Components:               []*types.MealComponent{},
			MinimumEstimatedPortions: database.Float32FromString(result.MealMinEstimatedPortions),
			EligibleForMealPlans:     result.MealEligibleForMealPlans,
		},
		MealScale: database.Float32FromString(result.MealScale),
		Chosen:    result.Chosen,
		TieBroken: result.Tiebroken,
	}

	return mealPlanOption, nil
}

// getMealPlanOptionsForMealPlanEvent fetches a list of meal plan options from the database that meet a particular filter.
func (q *Querier) getMealPlanOptionsForMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) ([]*types.MealPlanOption, error) {
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
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanEventID)

	results, err := q.generatedQuerier.GetAllMealPlanOptionsForMealPlanEvent(ctx, q.db, &generated.GetAllMealPlanOptionsForMealPlanEventParams{
		MealPlanID:      mealPlanID,
		MealPlanEventID: database.NullStringFromString(mealPlanEventID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	x := []*types.MealPlanOption{}
	for _, result := range results {
		x = append(x, &types.MealPlanOption{
			CreatedAt:              result.CreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
			Notes:                  result.Notes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
			ID:                     result.ID,
			Votes:                  nil,
			Meal: types.Meal{
				ID: result.MealID,
			},
			MealScale: database.Float32FromString(result.MealScale),
			Chosen:    result.Chosen,
			TieBroken: result.Tiebroken,
		})
	}

	for i, opt := range x {
		votes, voteFetchErr := q.GetMealPlanOptionVotesForMealPlanOption(ctx, mealPlanID, mealPlanEventID, opt.ID)
		if voteFetchErr != nil {
			return nil, observability.PrepareError(voteFetchErr, span, "fetching meal plan option votes for meal plan option")
		}
		x[i].Votes = votes

		m, mealFetchErr := q.GetMeal(ctx, opt.Meal.ID)
		if mealFetchErr != nil {
			return nil, observability.PrepareAndLogError(mealFetchErr, logger, span, "scanning meal plan options for meal plan event")
		}
		x[i].Meal = *m
	}

	logger.WithValue("quantity", len(x)).Info("fetched meal plan options for meal plan event")

	return x, nil
}

// GetMealPlanOptions fetches a list of meal plan options from the database that meet a particular filter.
func (q *Querier) GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.MealPlanOption], err error) {
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

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.MealPlanOption]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetMealPlanOptions(ctx, q.db, &generated.GetMealPlanOptionsParams{
		MealPlanEventID: database.NullStringFromString(mealPlanEventID),
		MealPlanID:      mealPlanID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan option query")
	}

	for _, result := range results {
		mealPlanOption := &types.MealPlanOption{
			CreatedAt:              result.CreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(result.LastUpdatedAt),
			AssignedCook:           database.StringPointerFromNullString(result.AssignedCook),
			ArchivedAt:             database.TimePointerFromNullTime(result.ArchivedAt),
			AssignedDishwasher:     database.StringPointerFromNullString(result.AssignedDishwasher),
			Notes:                  result.Notes,
			BelongsToMealPlanEvent: database.StringFromNullString(result.BelongsToMealPlanEvent),
			ID:                     result.ID,
			Votes:                  nil,
			Meal: types.Meal{
				ID: result.MealID,
			},
			MealScale: database.Float32FromString(result.MealScale),
			Chosen:    result.Chosen,
			TieBroken: result.Tiebroken,
		}
		x.Data = append(x.Data, mealPlanOption)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createMealPlanOption creates a meal plan option in the database.
func (q *Querier) createMealPlanOption(ctx context.Context, db database.SQLQueryExecutor, input *types.MealPlanOptionDatabaseCreationInput, markAsChosen bool) (*types.MealPlanOption, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, input.ID)
	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, input.ID)

	// create the meal plan option.
	if err := q.generatedQuerier.CreateMealPlanOption(ctx, db, &generated.CreateMealPlanOptionParams{
		ID:                     input.ID,
		AssignedCook:           database.NullStringFromStringPointer(input.AssignedCook),
		AssignedDishwasher:     database.NullStringFromStringPointer(input.AssignedDishwasher),
		MealID:                 input.MealID,
		Notes:                  input.Notes,
		MealScale:              database.StringFromFloat32(input.MealScale),
		BelongsToMealPlanEvent: database.NullStringFromString(input.BelongsToMealPlanEvent),
		Chosen:                 markAsChosen,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan option")
	}

	x := &types.MealPlanOption{
		ID:                     input.ID,
		AssignedCook:           input.AssignedCook,
		Meal:                   types.Meal{ID: input.MealID},
		Notes:                  input.Notes,
		BelongsToMealPlanEvent: input.BelongsToMealPlanEvent,
		CreatedAt:              q.currentTime(),
		MealScale:              input.MealScale,
		Votes:                  []*types.MealPlanOptionVote{},
	}

	logger.Info("meal plan option created")

	return x, nil
}

// CreateMealPlanOption creates a meal plan option in the database.
func (q *Querier) CreateMealPlanOption(ctx context.Context, input *types.MealPlanOptionDatabaseCreationInput) (*types.MealPlanOption, error) {
	return q.createMealPlanOption(ctx, q.db, input, false)
}

// UpdateMealPlanOption updates a particular meal plan option.
func (q *Querier) UpdateMealPlanOption(ctx context.Context, updated *types.MealPlanOption) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanOptionIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateMealPlanOption(ctx, q.db, &generated.UpdateMealPlanOptionParams{
		MealID:             updated.Meal.ID,
		Notes:              updated.Notes,
		MealScale:          database.StringFromFloat32(updated.MealScale),
		MealPlanOptionID:   updated.ID,
		AssignedCook:       database.NullStringFromStringPointer(updated.AssignedCook),
		AssignedDishwasher: database.NullStringFromStringPointer(updated.AssignedDishwasher),
		MealPlanEventID:    database.NullStringFromString(updated.BelongsToMealPlanEvent),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan option")
	}

	logger.Info("meal plan option updated")

	return nil
}

// ArchiveMealPlanOption archives a meal plan option from the database by its ID.
func (q *Querier) ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error {
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

	if mealPlanOptionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if _, err := q.generatedQuerier.ArchiveMealPlanOption(ctx, q.db, &generated.ArchiveMealPlanOptionParams{
		ID:                     mealPlanOptionID,
		BelongsToMealPlanEvent: sql.NullString{String: mealPlanEventID, Valid: true},
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan option")
	}

	logger.Info("meal plan option archived")

	return nil
}

func (q *Querier) determineWinner(winners []schulze.Result[string]) string {
	var (
		highestScore int
		scoreWinners []string
	)

	for _, winner := range winners {
		if winner.Wins == highestScore {
			scoreWinners = append(scoreWinners, winner.Choice)
		} else if winner.Wins > highestScore {
			highestScore = winner.Wins
			scoreWinners = []string{winner.Choice}
		}
	}

	/* #nosec: G404 */
	return scoreWinners[rand.N(len(scoreWinners))]
}

func (q *Querier) decideOptionWinner(ctx context.Context, options []*types.MealPlanOption) (_ string, _, _ bool) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	candidateMap := map[string]struct{}{}
	votesByUser := map[string]schulze.Ballot[string]{}

	logger := q.logger.WithValue("options.count", len(options))

	for _, option := range options {
		for _, v := range option.Votes {
			if votesByUser[v.ByUser] == nil {
				votesByUser[v.ByUser] = schulze.Ballot[string]{}
			}

			if !v.Abstain {
				votesByUser[v.ByUser][v.BelongsToMealPlanOption] = int(v.Rank)
			}

			candidateMap[v.BelongsToMealPlanOption] = struct{}{}
		}
	}

	candidates := []string{}
	for c := range candidateMap {
		candidates = append(candidates, c)
	}

	e := schulze.NewVoting(candidates)
	for _, vote := range votesByUser {
		if _, err := e.Vote(vote); err != nil {
			// this actually can never happen because we use uints for ranks, lol
			observability.AcknowledgeError(err, logger, span, "an invalid vote was received")
		}
	}

	winners, _, tie := e.Compute()
	if tie {
		return q.determineWinner(winners), true, true
	}

	if len(winners) > 0 {
		return winners[0].Choice, false, true
	}

	return "", false, false
}

// FinalizeMealPlanOption archives a meal plan option vote from the database by its ID.
func (q *Querier) FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, householdID string) (changed bool, err error) {
	_, span := q.tracer.StartSpan(ctx)
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

	if mealPlanOptionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanOptionIDKey, mealPlanOptionID)
	tracing.AttachToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	mealPlan, err := q.GetMealPlan(ctx, mealPlanID, householdID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching meal plan")
	}

	var (
		mealPlanEvent  *types.MealPlanEvent
		mealPlanOption *types.MealPlanOption
	)
	for _, event := range mealPlan.Events {
		if event.ID == mealPlanEventID {
			mealPlanEvent = event
			for _, option := range event.Options {
				if option.ID == mealPlanOptionID {
					mealPlanOption = option
					break
				}
			}
		}
	}

	// fetch household data
	household, err := q.GetHousehold(ctx, mealPlan.BelongsToHousehold)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "fetching household")
	}

	// go through all the votes for this meal plan option and determine if they're all there
	for _, member := range household.Members {
		memberVoteFound := false
		for _, vote := range mealPlanOption.Votes {
			if vote.ByUser == member.BelongsToUser.ID {
				memberVoteFound = true
				break
			}
		}

		if !memberVoteFound {
			return false, nil
		}
	}

	winner, tiebroken, chosen := q.decideOptionWinner(ctx, mealPlanEvent.Options)
	if chosen {
		if err = q.generatedQuerier.FinalizeMealPlanOption(ctx, q.db, &generated.FinalizeMealPlanOptionParams{
			MealPlanEventID: database.NullStringFromString(mealPlanEventID),
			ID:              winner,
			Tiebroken:       tiebroken,
		}); err != nil {
			return false, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan option")
		}

		logger.Debug("finalized meal plan option")
	}

	return chosen, nil
}
