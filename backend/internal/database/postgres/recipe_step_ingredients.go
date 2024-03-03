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
	_ types.RecipeStepIngredientDataManager = (*Querier)(nil)
)

// RecipeStepIngredientExists fetches whether a recipe step ingredient exists from the database.
func (q *Querier) RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (exists bool, err error) {
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

	if recipeStepIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	result, err := q.generatedQuerier.CheckRecipeStepIngredientExistence(ctx, q.db, &generated.CheckRecipeStepIngredientExistenceParams{
		RecipeStepID:           recipeStepID,
		RecipeStepIngredientID: recipeStepIngredientID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step ingredient existence check")
	}

	return result, nil
}

// GetRecipeStepIngredient fetches a recipe step ingredient from the database.
func (q *Querier) GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepIngredient, error) {
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

	if recipeStepIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	result, err := q.generatedQuerier.GetRecipeStepIngredient(ctx, q.db, &generated.GetRecipeStepIngredientParams{
		RecipeStepID:           recipeStepID,
		RecipeStepIngredientID: recipeStepIngredientID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step ingredient")
	}

	recipeStepIngredient := &types.RecipeStepIngredient{
		CreatedAt:                 result.CreatedAt,
		RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
		ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
		MaximumQuantity:           database.Float32PointerFromNullString(result.MaximumQuantityValue),
		VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
		ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
		RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
		QuantityNotes:             result.QuantityNotes,
		ID:                        result.ID,
		BelongsToRecipeStep:       result.BelongsToRecipeStep,
		IngredientNotes:           result.IngredientNotes,
		Name:                      result.Name,
		MeasurementUnit: types.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName,
			IconPath:      result.ValidMeasurementUnitIconPath,
			ID:            result.ValidMeasurementUnitID,
			Description:   result.ValidMeasurementUnitDescription,
			PluralName:    result.ValidMeasurementUnitPluralName,
			Slug:          result.ValidMeasurementUnitSlug,
			Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal,
			Metric:        result.ValidMeasurementUnitMetric,
			Imperial:      result.ValidMeasurementUnitImperial,
		},
		MinimumQuantity: database.Float32FromString(result.MinimumQuantityValue),
		OptionIndex:     uint16(result.OptionIndex),
		Optional:        result.Optional,
		ToTaste:         result.ToTaste,
	}

	if result.ValidIngredientID.Valid && result.ValidIngredientID.String != "" {
		recipeStepIngredient.Ingredient = &types.ValidIngredient{
			CreatedAt:                               result.ValidIngredientCreatedAt.Time,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.ValidIngredientIconPath.String,
			Warning:                                 result.ValidIngredientWarning.String,
			PluralName:                              result.ValidIngredientPluralName.String,
			StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
			Name:                                    result.ValidIngredientName.String,
			ID:                                      result.ValidIngredientID.String,
			Description:                             result.ValidIngredientDescription.String,
			Slug:                                    result.ValidIngredientSlug.String,
			ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
			ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
			IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
			IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
			ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
			ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
			ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
			ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
			AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
			RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
			ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
			ContainsFish:                            result.ValidIngredientContainsFish.Bool,
			ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
			ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
			ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
			AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
			IsStarch:                                result.ValidIngredientIsStarch.Bool,
			IsProtein:                               result.ValidIngredientIsProtein.Bool,
			IsGrain:                                 result.ValidIngredientIsGrain.Bool,
			IsFruit:                                 result.ValidIngredientIsFruit.Bool,
			IsSalt:                                  result.ValidIngredientIsSalt.Bool,
			IsFat:                                   result.ValidIngredientIsFat.Bool,
			IsAcid:                                  result.ValidIngredientIsAcid.Bool,
			IsHeat:                                  result.ValidIngredientIsHeat.Bool,
		}
	}

	return recipeStepIngredient, nil
}

// getRecipeStepIngredientsForRecipe fetches a list of recipe step ingredients from the database that meet a particular filter.
func (q *Querier) getRecipeStepIngredientsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.GetAllRecipeStepIngredientsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	recipeStepIngredients := []*types.RecipeStepIngredient{}
	for _, result := range results {
		recipeStepIngredient := &types.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
			MaximumQuantity:           database.Float32PointerFromNullString(result.MaximumQuantityValue),
			VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			MinimumQuantity: database.Float32FromString(result.MinimumQuantityValue),
			OptionIndex:     uint16(result.OptionIndex),
			Optional:        result.Optional,
			ToTaste:         result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			}
		}

		recipeStepIngredients = append(recipeStepIngredients, recipeStepIngredient)
	}

	return recipeStepIngredients, nil
}

// GetRecipeStepIngredients fetches a list of recipe step ingredients from the database that meet a particular filter.
func (q *Querier) GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepIngredient], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepIngredients(ctx, q.db, &generated.GetRecipeStepIngredientsParams{
		RecipeID:      recipeID,
		RecipeStepID:  recipeStepID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step ingredients list retrieval query")
	}

	for _, result := range results {
		recipeStepIngredient := &types.RecipeStepIngredient{
			CreatedAt:                 result.CreatedAt,
			RecipeStepProductID:       database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:                database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:             database.TimePointerFromNullTime(result.LastUpdatedAt),
			MaximumQuantity:           database.Float32PointerFromNullString(result.MaximumQuantityValue),
			VesselIndex:               database.Uint16PointerFromNullInt32(result.VesselIndex),
			ProductPercentageToUse:    database.Float32PointerFromNullString(result.ProductPercentageToUse),
			RecipeStepProductRecipeID: database.StringPointerFromNullString(result.RecipeStepProductRecipeID),
			QuantityNotes:             result.QuantityNotes,
			ID:                        result.ID,
			BelongsToRecipeStep:       result.BelongsToRecipeStep,
			IngredientNotes:           result.IngredientNotes,
			Name:                      result.Name,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			MinimumQuantity: database.Float32FromString(result.MinimumQuantityValue),
			OptionIndex:     uint16(result.OptionIndex),
			Optional:        result.Optional,
			ToTaste:         result.ToTaste,
		}

		if result.ValidIngredientID.Valid {
			recipeStepIngredient.Ingredient = &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:                              database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// createRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) createRecipeStepIngredient(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step ingredient.
	if err := q.generatedQuerier.CreateRecipeStepIngredient(ctx, db, &generated.CreateRecipeStepIngredientParams{
		QuantityNotes:             input.QuantityNotes,
		Name:                      input.Name,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		IngredientNotes:           input.IngredientNotes,
		ID:                        input.ID,
		MinimumQuantityValue:      database.StringFromFloat32(input.MinimumQuantity),
		RecipeStepProductID:       database.NullStringFromStringPointer(input.RecipeStepProductID),
		MaximumQuantityValue:      database.NullStringFromFloat32Pointer(input.MaximumQuantity),
		MeasurementUnit:           database.NullStringFromString(input.MeasurementUnitID),
		IngredientID:              database.NullStringFromStringPointer(input.IngredientID),
		ProductPercentageToUse:    database.NullStringFromFloat32Pointer(input.ProductPercentageToUse),
		RecipeStepProductRecipeID: database.NullStringFromStringPointer(input.RecipeStepProductRecipeID),
		VesselIndex:               database.NullInt32FromUint16Pointer(input.VesselIndex),
		OptionIndex:               int32(input.OptionIndex),
		ToTaste:                   input.ToTaste,
		Optional:                  input.Optional,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step ingredient creation query")
	}

	x := &types.RecipeStepIngredient{
		ID:                        input.ID,
		Name:                      input.Name,
		Optional:                  input.Optional,
		MeasurementUnit:           types.ValidMeasurementUnit{ID: input.MeasurementUnitID},
		MinimumQuantity:           input.MinimumQuantity,
		MaximumQuantity:           input.MaximumQuantity,
		QuantityNotes:             input.QuantityNotes,
		IngredientNotes:           input.IngredientNotes,
		BelongsToRecipeStep:       input.BelongsToRecipeStep,
		RecipeStepProductID:       input.RecipeStepProductID,
		OptionIndex:               input.OptionIndex,
		ToTaste:                   input.ToTaste,
		ProductPercentageToUse:    input.ProductPercentageToUse,
		VesselIndex:               input.VesselIndex,
		RecipeStepProductRecipeID: input.RecipeStepProductRecipeID,
		CreatedAt:                 q.currentTime(),
	}

	if input.IngredientID != nil {
		x.Ingredient = &types.ValidIngredient{ID: *input.IngredientID}
	}

	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, x.ID)

	return x, nil
}

// CreateRecipeStepIngredient creates a recipe step ingredient in the database.
func (q *Querier) CreateRecipeStepIngredient(ctx context.Context, input *types.RecipeStepIngredientDatabaseCreationInput) (*types.RecipeStepIngredient, error) {
	return q.createRecipeStepIngredient(ctx, q.db, input)
}

// UpdateRecipeStepIngredient updates a particular recipe step ingredient.
func (q *Querier) UpdateRecipeStepIngredient(ctx context.Context, updated *types.RecipeStepIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepIngredientIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, updated.ID)

	var ingredientID *string
	if updated.Ingredient != nil {
		ingredientID = &updated.Ingredient.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepIngredient(ctx, q.db, &generated.UpdateRecipeStepIngredientParams{
		IngredientID:              database.NullStringFromStringPointer(ingredientID),
		Name:                      updated.Name,
		Optional:                  updated.Optional,
		MeasurementUnit:           database.NullStringFromString(updated.MeasurementUnit.ID),
		MinimumQuantityValue:      database.StringFromFloat32(updated.MinimumQuantity),
		MaximumQuantityValue:      database.NullStringFromFloat32Pointer(updated.MaximumQuantity),
		QuantityNotes:             updated.QuantityNotes,
		RecipeStepProductID:       database.NullStringFromStringPointer(updated.RecipeStepProductID),
		IngredientNotes:           updated.IngredientNotes,
		OptionIndex:               int32(updated.OptionIndex),
		ToTaste:                   updated.ToTaste,
		ProductPercentageToUse:    database.NullStringFromFloat32Pointer(updated.ProductPercentageToUse),
		VesselIndex:               database.NullInt32FromUint16Pointer(updated.VesselIndex),
		RecipeStepProductRecipeID: database.NullStringFromStringPointer(updated.RecipeStepProductRecipeID),
		BelongsToRecipeStep:       updated.BelongsToRecipeStep,
		ID:                        updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step ingredient")
	}

	logger.Info("recipe step ingredient updated")

	return nil
}

// ArchiveRecipeStepIngredient archives a recipe step ingredient from the database by its ID.
func (q *Querier) ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
	tracing.AttachToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepIngredient(ctx, q.db, &generated.ArchiveRecipeStepIngredientParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepIngredientID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step ingredient")
	}

	logger.Info("recipe step ingredient archived")

	return nil
}
