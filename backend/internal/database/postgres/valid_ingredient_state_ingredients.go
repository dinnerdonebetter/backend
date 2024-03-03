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
	_ types.ValidIngredientStateIngredientDataManager = (*Querier)(nil)
)

// ValidIngredientStateIngredientExists fetches whether a valid ingredient state ingredient exists from the database.
func (q *Querier) ValidIngredientStateIngredientExists(ctx context.Context, validIngredientStateIngredientID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	result, err := q.generatedQuerier.CheckValidIngredientStateIngredientExistence(ctx, q.db, validIngredientStateIngredientID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state ingredient existence check")
	}

	return result, nil
}

// GetValidIngredientStateIngredient fetches a valid ingredient state ingredient from the database.
func (q *Querier) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	result, err := q.generatedQuerier.GetValidIngredientStateIngredient(ctx, q.db, validIngredientStateIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid ingredient state ingredient")
	}

	validIngredientStateIngredient := &types.ValidIngredientStateIngredient{
		CreatedAt:     result.ValidIngredientStateIngredientCreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateIngredientLastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateIngredientArchivedAt),
		Notes:         result.ValidIngredientStateIngredientNotes,
		ID:            result.ValidIngredientStateIngredientID,
		IngredientState: types.ValidIngredientState{
			CreatedAt:     result.ValidIngredientStateCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
			PastTense:     result.ValidIngredientStatePastTense,
			Description:   result.ValidIngredientStateDescription,
			IconPath:      result.ValidIngredientStateIconPath,
			ID:            result.ValidIngredientStateID,
			Name:          result.ValidIngredientStateName,
			AttributeType: string(result.ValidIngredientStateAttributeType),
			Slug:          result.ValidIngredientStateSlug,
		},
		Ingredient: types.ValidIngredient{
			CreatedAt:                               result.ValidIngredientCreatedAt,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.ValidIngredientIconPath,
			Warning:                                 result.ValidIngredientWarning,
			PluralName:                              result.ValidIngredientPluralName,
			StorageInstructions:                     result.ValidIngredientStorageInstructions,
			Name:                                    result.ValidIngredientName,
			ID:                                      result.ValidIngredientID,
			Description:                             result.ValidIngredientDescription,
			Slug:                                    result.ValidIngredientSlug,
			ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions,
			ContainsShellfish:                       result.ValidIngredientContainsShellfish,
			IsMeasuredVolumetrically:                result.ValidIngredientVolumetric,
			IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:                          result.ValidIngredientContainsPeanut,
			ContainsTreeNut:                         result.ValidIngredientContainsTreeNut,
			ContainsEgg:                             result.ValidIngredientContainsEgg,
			ContainsWheat:                           result.ValidIngredientContainsWheat,
			ContainsSoy:                             result.ValidIngredientContainsSoy,
			AnimalDerived:                           result.ValidIngredientAnimalDerived,
			RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations,
			ContainsSesame:                          result.ValidIngredientContainsSesame,
			ContainsFish:                            result.ValidIngredientContainsFish,
			ContainsGluten:                          result.ValidIngredientContainsGluten,
			ContainsDairy:                           result.ValidIngredientContainsDairy,
			ContainsAlcohol:                         result.ValidIngredientContainsAlcohol,
			AnimalFlesh:                             result.ValidIngredientAnimalFlesh,
			IsStarch:                                result.ValidIngredientIsStarch,
			IsProtein:                               result.ValidIngredientIsProtein,
			IsGrain:                                 result.ValidIngredientIsGrain,
			IsFruit:                                 result.ValidIngredientIsFruit,
			IsSalt:                                  result.ValidIngredientIsSalt,
			IsFat:                                   result.ValidIngredientIsFat,
			IsAcid:                                  result.ValidIngredientIsAcid,
			IsHeat:                                  result.ValidIngredientIsHeat,
		},
	}

	return validIngredientStateIngredient, nil
}

// GetValidIngredientStateIngredients fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientStateIngredients(ctx, q.db, &generated.GetValidIngredientStateIngredientsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	for _, result := range results {
		validIngredientStateIngredient := &types.ValidIngredientStateIngredient{
			CreatedAt:     result.ValidIngredientStateIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateIngredientArchivedAt),
			Notes:         result.ValidIngredientStateIngredientNotes,
			ID:            result.ValidIngredientStateIngredientID,
			IngredientState: types.ValidIngredientState{
				CreatedAt:     result.ValidIngredientStateCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
				PastTense:     result.ValidIngredientStatePastTense,
				Description:   result.ValidIngredientStateDescription,
				IconPath:      result.ValidIngredientStateIconPath,
				ID:            result.ValidIngredientStateID,
				Name:          result.ValidIngredientStateName,
				AttributeType: string(result.ValidIngredientStateAttributeType),
				Slug:          result.ValidIngredientStateSlug,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt,
				LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath,
				Warning:                                 result.ValidIngredientWarning,
				PluralName:                              result.ValidIngredientPluralName,
				StorageInstructions:                     result.ValidIngredientStorageInstructions,
				Name:                                    result.ValidIngredientName,
				ID:                                      result.ValidIngredientID,
				Description:                             result.ValidIngredientDescription,
				Slug:                                    result.ValidIngredientSlug,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric,
				IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut,
				ContainsEgg:                             result.ValidIngredientContainsEgg,
				ContainsWheat:                           result.ValidIngredientContainsWheat,
				ContainsSoy:                             result.ValidIngredientContainsSoy,
				AnimalDerived:                           result.ValidIngredientAnimalDerived,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations,
				ContainsSesame:                          result.ValidIngredientContainsSesame,
				ContainsFish:                            result.ValidIngredientContainsFish,
				ContainsGluten:                          result.ValidIngredientContainsGluten,
				ContainsDairy:                           result.ValidIngredientContainsDairy,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh,
				IsStarch:                                result.ValidIngredientIsStarch,
				IsProtein:                               result.ValidIngredientIsProtein,
				IsGrain:                                 result.ValidIngredientIsGrain,
				IsFruit:                                 result.ValidIngredientIsFruit,
				IsSalt:                                  result.ValidIngredientIsSalt,
				IsFat:                                   result.ValidIngredientIsFat,
				IsAcid:                                  result.ValidIngredientIsAcid,
				IsHeat:                                  result.ValidIngredientIsHeat,
			},
		}

		x.Data = append(x.Data, validIngredientStateIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientStateIngredientsForIngredientState fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientStateID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, ingredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, ingredientStateID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientStateIngredientsForIngredientState(ctx, q.db, &generated.GetValidIngredientStateIngredientsForIngredientStateParams{
		CreatedBefore:        database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:         database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:        database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:         database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:          database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:           database.NullInt32FromUint8Pointer(filter.Limit),
		ValidIngredientState: ingredientStateID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	for _, result := range results {
		validIngredientStateIngredient := &types.ValidIngredientStateIngredient{
			CreatedAt:     result.ValidIngredientStateIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateIngredientArchivedAt),
			Notes:         result.ValidIngredientStateIngredientNotes,
			ID:            result.ValidIngredientStateIngredientID,
			IngredientState: types.ValidIngredientState{
				CreatedAt:     result.ValidIngredientStateCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
				PastTense:     result.ValidIngredientStatePastTense,
				Description:   result.ValidIngredientStateDescription,
				IconPath:      result.ValidIngredientStateIconPath,
				ID:            result.ValidIngredientStateID,
				Name:          result.ValidIngredientStateName,
				AttributeType: string(result.ValidIngredientStateAttributeType),
				Slug:          result.ValidIngredientStateSlug,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt,
				LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath,
				Warning:                                 result.ValidIngredientWarning,
				PluralName:                              result.ValidIngredientPluralName,
				StorageInstructions:                     result.ValidIngredientStorageInstructions,
				Name:                                    result.ValidIngredientName,
				ID:                                      result.ValidIngredientID,
				Description:                             result.ValidIngredientDescription,
				Slug:                                    result.ValidIngredientSlug,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric,
				IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut,
				ContainsEgg:                             result.ValidIngredientContainsEgg,
				ContainsWheat:                           result.ValidIngredientContainsWheat,
				ContainsSoy:                             result.ValidIngredientContainsSoy,
				AnimalDerived:                           result.ValidIngredientAnimalDerived,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations,
				ContainsSesame:                          result.ValidIngredientContainsSesame,
				ContainsFish:                            result.ValidIngredientContainsFish,
				ContainsGluten:                          result.ValidIngredientContainsGluten,
				ContainsDairy:                           result.ValidIngredientContainsDairy,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh,
				IsStarch:                                result.ValidIngredientIsStarch,
				IsProtein:                               result.ValidIngredientIsProtein,
				IsGrain:                                 result.ValidIngredientIsGrain,
				IsFruit:                                 result.ValidIngredientIsFruit,
				IsSalt:                                  result.ValidIngredientIsSalt,
				IsFat:                                   result.ValidIngredientIsFat,
				IsAcid:                                  result.ValidIngredientIsAcid,
				IsHeat:                                  result.ValidIngredientIsHeat,
			},
		}

		x.Data = append(x.Data, validIngredientStateIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientStateIngredientsForIngredient fetches a list of valid ingredient state ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientStateIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, ingredientID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredientStateIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientStateIngredientsForIngredient(ctx, q.db, &generated.GetValidIngredientStateIngredientsForIngredientParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		ValidIngredient: ingredientID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state ingredients list retrieval query")
	}

	for _, result := range results {
		validIngredientStateIngredient := &types.ValidIngredientStateIngredient{
			CreatedAt:     result.ValidIngredientStateIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateIngredientArchivedAt),
			Notes:         result.ValidIngredientStateIngredientNotes,
			ID:            result.ValidIngredientStateIngredientID,
			IngredientState: types.ValidIngredientState{
				CreatedAt:     result.ValidIngredientStateCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientStateLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientStateArchivedAt),
				PastTense:     result.ValidIngredientStatePastTense,
				Description:   result.ValidIngredientStateDescription,
				IconPath:      result.ValidIngredientStateIconPath,
				ID:            result.ValidIngredientStateID,
				Name:          result.ValidIngredientStateName,
				AttributeType: string(result.ValidIngredientStateAttributeType),
				Slug:          result.ValidIngredientStateSlug,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt,
				LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath,
				Warning:                                 result.ValidIngredientWarning,
				PluralName:                              result.ValidIngredientPluralName,
				StorageInstructions:                     result.ValidIngredientStorageInstructions,
				Name:                                    result.ValidIngredientName,
				ID:                                      result.ValidIngredientID,
				Description:                             result.ValidIngredientDescription,
				Slug:                                    result.ValidIngredientSlug,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric,
				IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut,
				ContainsEgg:                             result.ValidIngredientContainsEgg,
				ContainsWheat:                           result.ValidIngredientContainsWheat,
				ContainsSoy:                             result.ValidIngredientContainsSoy,
				AnimalDerived:                           result.ValidIngredientAnimalDerived,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations,
				ContainsSesame:                          result.ValidIngredientContainsSesame,
				ContainsFish:                            result.ValidIngredientContainsFish,
				ContainsGluten:                          result.ValidIngredientContainsGluten,
				ContainsDairy:                           result.ValidIngredientContainsDairy,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh,
				IsStarch:                                result.ValidIngredientIsStarch,
				IsProtein:                               result.ValidIngredientIsProtein,
				IsGrain:                                 result.ValidIngredientIsGrain,
				IsFruit:                                 result.ValidIngredientIsFruit,
				IsSalt:                                  result.ValidIngredientIsSalt,
				IsFat:                                   result.ValidIngredientIsFat,
				IsAcid:                                  result.ValidIngredientIsAcid,
				IsHeat:                                  result.ValidIngredientIsHeat,
			},
		}

		x.Data = append(x.Data, validIngredientStateIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateValidIngredientStateIngredient creates a valid ingredient state ingredient in the database.
func (q *Querier) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientDatabaseCreationInput) (*types.ValidIngredientStateIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientStateIngredientIDKey, input.ID)

	// create the valid ingredient state ingredient.
	if err := q.generatedQuerier.CreateValidIngredientStateIngredient(ctx, q.db, &generated.CreateValidIngredientStateIngredientParams{
		ID:                   input.ID,
		Notes:                input.Notes,
		ValidIngredientState: input.ValidIngredientStateID,
		ValidIngredient:      input.ValidIngredientID,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state ingredient creation query")
	}

	x := &types.ValidIngredientStateIngredient{
		ID:              input.ID,
		Notes:           input.Notes,
		IngredientState: types.ValidIngredientState{ID: input.ValidIngredientStateID},
		Ingredient:      types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedAt:       q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, x.ID)
	logger.Info("valid ingredient state ingredient created")

	return x, nil
}

// UpdateValidIngredientStateIngredient updates a particular valid ingredient state ingredient.
func (q *Querier) UpdateValidIngredientStateIngredient(ctx context.Context, updated *types.ValidIngredientStateIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientStateIngredientIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredientStateIngredient(ctx, q.db, &generated.UpdateValidIngredientStateIngredientParams{
		Notes:                updated.Notes,
		ValidIngredientState: updated.IngredientState.ID,
		ValidIngredient:      updated.Ingredient.ID,
		ID:                   updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state ingredient")
	}

	logger.Info("valid ingredient state ingredient updated")

	return nil
}

// ArchiveValidIngredientStateIngredient archives a valid ingredient state ingredient from the database by its ID.
func (q *Querier) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)

	if _, err := q.generatedQuerier.ArchiveValidIngredientStateIngredient(ctx, q.db, validIngredientStateIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state ingredient")
	}

	logger.Info("valid ingredient state ingredient archived")

	return nil
}
