package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.RecipeStepCompletionConditionDataManager = (*Querier)(nil)
)

// RecipeStepCompletionConditionExists fetches whether a recipe step completion condition exists from the database.
func (q *Querier) RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepCompletionConditionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	result, err := q.generatedQuerier.CheckRecipeStepCompletionConditionExistence(ctx, q.db, &generated.CheckRecipeStepCompletionConditionExistenceParams{
		RecipeStepID:                    recipeStepID,
		RecipeStepCompletionConditionID: recipeStepCompletionConditionID,
		RecipeID:                        recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step completion condition existence check")
	}

	return result, nil
}

// GetRecipeStepCompletionCondition fetches a recipe step completion condition from the database.
func (q *Querier) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepCompletionConditionID string) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepCompletionConditionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	results, err := q.generatedQuerier.GetRecipeStepCompletionConditionWithIngredients(ctx, q.db, &generated.GetRecipeStepCompletionConditionWithIngredientsParams{
		RecipeID:                        recipeID,
		RecipeStepID:                    recipeStepID,
		RecipeStepCompletionConditionID: recipeStepCompletionConditionID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for recipe step completion condition")
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	recipeStepCompletionCondition := &types.RecipeStepCompletionCondition{}
	for _, result := range results {
		if recipeStepCompletionCondition.ID == "" {
			recipeStepCompletionCondition = &types.RecipeStepCompletionCondition{
				CreatedAt:     result.CreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
				IngredientState: types.ValidIngredientState{
					CreatedAt:     result.ValidIngredientStateCreatedAt,
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
					PastTense:     result.ValidIngredientStatePastTense,
					Description:   result.ValidIngredientStateDescription,
					IconPath:      result.ValidIngredientStateIconPath,
					ID:            result.ValidIngredientStateID,
					Name:          result.ValidIngredientStateName,
					AttributeType: string(result.ValidIngredientStateAttributeType),
					Slug:          result.ValidIngredientStateSlug,
				},
				ID:                  result.ID,
				BelongsToRecipeStep: result.BelongsToRecipeStep,
				Notes:               result.Notes,
				Ingredients:         []*types.RecipeStepCompletionConditionIngredient{},
				Optional:            result.Optional,
			}
		}

		recipeStepCompletionCondition.Ingredients = append(recipeStepCompletionCondition.Ingredients, &types.RecipeStepCompletionConditionIngredient{
			CreatedAt:                              result.RecipeStepCompletionConditionIngredientCreatedAt,
			ArchivedAt:                             database.TimePointerFromNullTime(result.RecipeStepCompletionConditionIngredientArchivedAt),
			LastUpdatedAt:                          database.TimePointerFromNullTime(result.RecipeStepCompletionConditionIngredientLastUpdatedAt),
			ID:                                     result.RecipeStepCompletionConditionIngredientID,
			BelongsToRecipeStepCompletionCondition: result.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
			RecipeStepIngredient:                   result.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
		})
	}

	return recipeStepCompletionCondition, nil
}

// getRecipeStepCompletionConditionsForRecipe fetches a recipe step completion condition from the database.
func (q *Querier) getRecipeStepCompletionConditionsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepCompletionCondition, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.GetAllRecipeStepCompletionConditionsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for recipe step completion condition")
	}

	idOrder := []string{}
	byID := map[string]*types.RecipeStepCompletionCondition{}
	for _, result := range results {
		recipeStepCompletionCondition := &types.RecipeStepCompletionCondition{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			IngredientState: types.ValidIngredientState{
				CreatedAt:     time.Time{},
				ArchivedAt:    nil,
				LastUpdatedAt: nil,
				PastTense:     "",
				Description:   "",
				IconPath:      "",
				ID:            "",
				Name:          "",
				AttributeType: "",
				Slug:          "",
			},
			ID:                  result.ID,
			BelongsToRecipeStep: result.BelongsToRecipeStep,
			Notes:               result.Notes,
			Optional:            result.Optional,
		}

		if byID[recipeStepCompletionCondition.ID] == nil {
			byID[recipeStepCompletionCondition.ID] = recipeStepCompletionCondition
		} else {
			idOrder = append(idOrder, recipeStepCompletionCondition.ID)
			byID[recipeStepCompletionCondition.ID].Ingredients = append(byID[recipeStepCompletionCondition.ID].Ingredients, &types.RecipeStepCompletionConditionIngredient{
				ID:                                     result.RecipeStepCompletionConditionIngredientID,
				BelongsToRecipeStepCompletionCondition: result.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
				RecipeStepIngredient:                   result.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			})
			byID[recipeStepCompletionCondition.ID] = recipeStepCompletionCondition
		}
	}

	recipeStepConditions := []*types.RecipeStepCompletionCondition{}
	for _, id := range idOrder {
		recipeStepConditions = append(recipeStepConditions, byID[id])
	}

	return recipeStepConditions, nil
}

// GetRecipeStepCompletionConditions fetches a list of recipe step completion conditions from the database that meet a particular filter.
func (q *Querier) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepCompletionCondition], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeStepCompletionCondition]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepCompletionConditions(ctx, q.db, &generated.GetRecipeStepCompletionConditionsParams{
		RecipeStepID:  recipeStepID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for recipe step completion conditions")
	}

	idOrder := []string{}
	byID := map[string]*types.RecipeStepCompletionCondition{}
	for _, result := range results {
		recipeStepCompletionCondition := &types.RecipeStepCompletionCondition{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			IngredientState: types.ValidIngredientState{
				CreatedAt:     result.ValidIngredientStateCreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
				PastTense:     result.ValidIngredientStatePastTense,
				Description:   result.ValidIngredientStateDescription,
				IconPath:      result.ValidIngredientStateIconPath,
				ID:            result.ValidIngredientStateID,
				Name:          result.ValidIngredientStateName,
				AttributeType: string(result.ValidIngredientStateAttributeType),
				Slug:          result.ValidIngredientStateSlug,
			},
			ID:                  result.ID,
			BelongsToRecipeStep: result.BelongsToRecipeStep,
			Notes:               result.Notes,
			Optional:            result.Optional,
		}

		if byID[recipeStepCompletionCondition.ID] == nil {
			byID[recipeStepCompletionCondition.ID] = recipeStepCompletionCondition
			idOrder = append(idOrder, recipeStepCompletionCondition.ID)
		} else {
			byID[recipeStepCompletionCondition.ID].Ingredients = append(byID[recipeStepCompletionCondition.ID].Ingredients, &types.RecipeStepCompletionConditionIngredient{
				ID:                                     result.RecipeStepCompletionConditionIngredientID,
				BelongsToRecipeStepCompletionCondition: result.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
				RecipeStepIngredient:                   result.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			})
			byID[recipeStepCompletionCondition.ID] = recipeStepCompletionCondition
		}

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	for _, id := range idOrder {
		x.Data = append(x.Data, byID[id])
	}

	return x, nil
}

// createRecipeStepCompletionCondition creates a recipe step completion condition in the database.
func (q *Querier) createRecipeStepCompletionCondition(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step completion condition.
	if err := q.generatedQuerier.CreateRecipeStepCompletionCondition(ctx, db, &generated.CreateRecipeStepCompletionConditionParams{
		ID:                  input.ID,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		IngredientState:     input.IngredientStateID,
		Notes:               input.Notes,
		Optional:            input.Optional,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step completion condition creation query")
	}

	x := &types.RecipeStepCompletionCondition{
		ID:                  input.ID,
		Notes:               input.Notes,
		Optional:            input.Optional,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		IngredientState:     types.ValidIngredientState{ID: input.IngredientStateID},
		CreatedAt:           q.currentTime(),
	}

	for _, ingredient := range input.Ingredients {
		ingredient.BelongsToRecipeStepCompletionCondition = x.ID
		completionConditionIngredient, err := q.createRecipeStepCompletionConditionIngredient(ctx, db, ingredient)
		if err != nil {
			return nil, observability.PrepareError(err, span, "creating ingredient for recipe step completion condition")
		}

		x.Ingredients = append(x.Ingredients, completionConditionIngredient)
	}

	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, x.ID)

	return x, nil
}

// createRecipeStepCompletionConditionIngredient creates a recipe step completion condition ingredient in the database.
func (q *Querier) createRecipeStepCompletionConditionIngredient(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepCompletionConditionIngredientDatabaseCreationInput) (*types.RecipeStepCompletionConditionIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step completion condition.
	if err := q.generatedQuerier.CreateRecipeStepCompletionConditionIngredient(ctx, db, &generated.CreateRecipeStepCompletionConditionIngredientParams{
		ID:                                     input.ID,
		BelongsToRecipeStepCompletionCondition: input.BelongsToRecipeStepCompletionCondition,
		RecipeStepIngredient:                   input.RecipeStepIngredient,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step completion condition ingredient creation query")
	}

	x := &types.RecipeStepCompletionConditionIngredient{
		ID:                                     input.ID,
		BelongsToRecipeStepCompletionCondition: input.BelongsToRecipeStepCompletionCondition,
		RecipeStepIngredient:                   input.RecipeStepIngredient,
		CreatedAt:                              q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, x.ID)

	return x, nil
}

// CreateRecipeStepCompletionCondition creates a recipe step completion condition in the database.
func (q *Querier) CreateRecipeStepCompletionCondition(ctx context.Context, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	return q.createRecipeStepCompletionCondition(ctx, q.db, input)
}

// UpdateRecipeStepCompletionCondition updates a particular recipe step completion condition.
func (q *Querier) UpdateRecipeStepCompletionCondition(ctx context.Context, updated *types.RecipeStepCompletionCondition) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepCompletionConditionIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateRecipeStepCompletionCondition(ctx, q.db, &generated.UpdateRecipeStepCompletionConditionParams{
		Optional:            updated.Optional,
		Notes:               updated.Notes,
		BelongsToRecipeStep: updated.BelongsToRecipeStep,
		IngredientState:     updated.IngredientState.ID,
		ID:                  updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	logger.Info("recipe step completion condition updated")

	return nil
}

// ArchiveRecipeStepCompletionCondition archives a recipe step completion condition from the database by its ID.
func (q *Querier) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepCompletionConditionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepCompletionConditionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
	tracing.AttachToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepCompletionCondition(ctx, q.db, &generated.ArchiveRecipeStepCompletionConditionParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepCompletionConditionID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step completion condition")
	}

	logger.Info("recipe step completion condition archived")

	return nil
}
