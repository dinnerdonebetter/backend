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
	_ types.MealPlanGroceryListItemDataManager = (*Querier)(nil)
)

// MealPlanGroceryListItemExists fetches whether a meal plan grocery list exists from the database.
func (q *Querier) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	result, err := q.generatedQuerier.CheckMealPlanGroceryListItemExistence(ctx, q.db, &generated.CheckMealPlanGroceryListItemExistenceParams{
		MealPlanID:                mealPlanID,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list existence check")
	}

	return result, nil
}

func (q *Querier) fleshOutMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItem *types.MealPlanGroceryListItem) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItem == nil {
		return nil, ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItem.ID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItem.ID)

	validIngredient, err := q.GetValidIngredient(ctx, mealPlanGroceryListItem.Ingredient.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item ingredient")
	}
	mealPlanGroceryListItem.Ingredient = *validIngredient

	validMeasurementUnit, err := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.MeasurementUnit.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching grocery list item measurement unit")
	}
	mealPlanGroceryListItem.MeasurementUnit = *validMeasurementUnit

	if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
		if getPurchasedMeasurementUnitErr != nil {
			return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
		}
		mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
	}

	return mealPlanGroceryListItem, nil
}

// GetMealPlanGroceryListItem fetches a meal plan grocery list from the database.
func (q *Querier) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanGroceryListItemID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	result, err := q.generatedQuerier.GetMealPlanGroceryListItem(ctx, q.db, &generated.GetMealPlanGroceryListItemParams{
		MealPlanID:                mealPlanID,
		MealPlanGroceryListItemID: mealPlanGroceryListItemID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meal plan grocery list item")
	}

	mealPlanGroceryListItem := &types.MealPlanGroceryListItem{
		CreatedAt:             result.CreatedAt,
		MaximumQuantityNeeded: database.Float32PointerFromNullString(result.MaximumQuantityNeeded),
		LastUpdatedAt:         database.TimePointerFromNullTime(result.LastUpdatedAt),
		PurchasePrice:         database.Float32PointerFromNullString(result.PurchasePrice),
		PurchasedUPC:          database.StringPointerFromNullString(result.PurchasedUpc),
		ArchivedAt:            database.TimePointerFromNullTime(result.ArchivedAt),
		QuantityPurchased:     database.Float32PointerFromNullString(result.QuantityPurchased),

		BelongsToMealPlan:     result.BelongsToMealPlan,
		Status:                string(result.Status),
		StatusExplanation:     result.StatusExplanation,
		ID:                    result.ID,
		MinimumQuantityNeeded: database.Float32FromString(result.MinimumQuantityNeeded),
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

	if result.PurchasedMeasurementUnit.Valid {
		mealPlanGroceryListItem.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{
			ID: result.PurchasedMeasurementUnit.String,
		}
	}

	if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
		if getPurchasedMeasurementUnitErr != nil {
			return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
		}
		mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
	}

	return mealPlanGroceryListItem, nil
}

// GetMealPlanGroceryListItemsForMealPlan fetches a list of meal plan grocery lists from the database that meet a particular filter.
func (q *Querier) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	results, err := q.generatedQuerier.GetMealPlanGroceryListItemsForMealPlan(ctx, q.db, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing meal plan grocery list items list retrieval query")
	}

	x := []*types.MealPlanGroceryListItem{}
	for _, result := range results {
		mealPlanGroceryListItem := &types.MealPlanGroceryListItem{
			CreatedAt:             result.CreatedAt,
			MaximumQuantityNeeded: database.Float32PointerFromNullString(result.MaximumQuantityNeeded),
			LastUpdatedAt:         database.TimePointerFromNullTime(result.LastUpdatedAt),
			PurchasePrice:         database.Float32PointerFromNullString(result.PurchasePrice),
			PurchasedUPC:          database.StringPointerFromNullString(result.PurchasedUpc),
			ArchivedAt:            database.TimePointerFromNullTime(result.ArchivedAt),
			QuantityPurchased:     database.Float32PointerFromNullString(result.QuantityPurchased),

			BelongsToMealPlan:     result.BelongsToMealPlan,
			Status:                string(result.Status),
			StatusExplanation:     result.StatusExplanation,
			ID:                    result.ID,
			MinimumQuantityNeeded: database.Float32FromString(result.MinimumQuantityNeeded),
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

		if result.PurchasedMeasurementUnit.Valid {
			mealPlanGroceryListItem.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{
				ID: result.PurchasedMeasurementUnit.String,
			}
		}

		if mealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
			purchasedMeasurementUnit, getPurchasedMeasurementUnitErr := q.GetValidMeasurementUnit(ctx, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
			if getPurchasedMeasurementUnitErr != nil {
				return nil, observability.PrepareAndLogError(getPurchasedMeasurementUnitErr, logger, span, "fetching grocery list item purchased measurement unit")
			}
			mealPlanGroceryListItem.PurchasedMeasurementUnit = purchasedMeasurementUnit
		}

		x = append(x, mealPlanGroceryListItem)
	}

	for i := range x {
		x[i], err = q.fleshOutMealPlanGroceryListItem(ctx, x[i])
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "augmenting grocery list item data")
		}
	}

	return x, nil
}

// createMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) createMealPlanGroceryListItem(ctx context.Context, querier database.SQLQueryExecutor, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, input.ID)

	// create the meal plan grocery list.
	if err := q.generatedQuerier.CreateMealPlanGroceryListItem(ctx, querier, &generated.CreateMealPlanGroceryListItemParams{
		ID:                       input.ID,
		BelongsToMealPlan:        input.BelongsToMealPlan,
		ValidIngredient:          input.ValidIngredientID,
		ValidMeasurementUnit:     input.ValidMeasurementUnitID,
		MinimumQuantityNeeded:    database.StringFromFloat32(input.MinimumQuantityNeeded),
		StatusExplanation:        input.StatusExplanation,
		Status:                   generated.GroceryListItemStatus(input.Status),
		MaximumQuantityNeeded:    database.NullStringFromFloat32Pointer(input.MaximumQuantityNeeded),
		QuantityPurchased:        database.NullStringFromFloat32Pointer(input.QuantityPurchased),
		PurchasedMeasurementUnit: database.NullStringFromStringPointer(input.PurchasedMeasurementUnitID),
		PurchasedUpc:             database.NullStringFromStringPointer(input.PurchasedUPC),
		PurchasePrice:            database.NullStringFromFloat32Pointer(input.PurchasePrice),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing meal plan grocery list creation query")
	}

	x := &types.MealPlanGroceryListItem{
		ID:                    input.ID,
		BelongsToMealPlan:     input.BelongsToMealPlan,
		Ingredient:            types.ValidIngredient{ID: input.ValidIngredientID},
		MeasurementUnit:       types.ValidMeasurementUnit{ID: input.ValidMeasurementUnitID},
		MinimumQuantityNeeded: input.MinimumQuantityNeeded,
		MaximumQuantityNeeded: input.MaximumQuantityNeeded,
		QuantityPurchased:     input.QuantityPurchased,
		PurchasedUPC:          input.PurchasedUPC,
		PurchasePrice:         input.PurchasePrice,
		StatusExplanation:     input.StatusExplanation,
		Status:                input.Status,
		CreatedAt:             q.currentTime(),
	}

	if input.PurchasedMeasurementUnitID != nil {
		x.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{ID: *input.PurchasedMeasurementUnitID}
	}

	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, x.ID)
	logger.Info("meal plan grocery list created")

	return x, nil
}

// CreateMealPlanGroceryListItem creates a meal plan grocery list in the database.
func (q *Querier) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	return q.createMealPlanGroceryListItem(ctx, q.db, input)
}

// UpdateMealPlanGroceryListItem updates a particular meal plan grocery list.
func (q *Querier) UpdateMealPlanGroceryListItem(ctx context.Context, updated *types.MealPlanGroceryListItem) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.MealPlanGroceryListItemIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, updated.ID)

	var purchasedMeasurementUnitID *string
	if updated.PurchasedMeasurementUnit != nil {
		purchasedMeasurementUnitID = &updated.PurchasedMeasurementUnit.ID
	}

	if _, err := q.generatedQuerier.UpdateMealPlanGroceryListItem(ctx, q.db, &generated.UpdateMealPlanGroceryListItemParams{
		BelongsToMealPlan:        updated.BelongsToMealPlan,
		ValidIngredient:          updated.Ingredient.ID,
		ValidMeasurementUnit:     updated.MeasurementUnit.ID,
		MinimumQuantityNeeded:    database.StringFromFloat32(updated.MinimumQuantityNeeded),
		StatusExplanation:        updated.StatusExplanation,
		Status:                   generated.GroceryListItemStatus(updated.Status),
		ID:                       updated.ID,
		MaximumQuantityNeeded:    database.NullStringFromFloat32Pointer(updated.MaximumQuantityNeeded),
		QuantityPurchased:        database.NullStringFromFloat32Pointer(updated.QuantityPurchased),
		PurchasedMeasurementUnit: database.NullStringFromStringPointer(purchasedMeasurementUnitID),
		PurchasedUpc:             database.NullStringFromStringPointer(updated.PurchasedUPC),
		PurchasePrice:            database.NullStringFromFloat32Pointer(updated.PurchasePrice),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan grocery list")
	}

	logger.Info("meal plan grocery list updated")

	return nil
}

// ArchiveMealPlanGroceryListItem archives a meal plan grocery list from the database by its ID.
func (q *Querier) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if mealPlanGroceryListItemID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
	tracing.AttachToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)

	if _, err := q.generatedQuerier.ArchiveMealPlanGroceryListItem(ctx, q.db, mealPlanGroceryListItemID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan grocery list")
	}

	logger.Info("meal plan grocery list archived")

	return nil
}
