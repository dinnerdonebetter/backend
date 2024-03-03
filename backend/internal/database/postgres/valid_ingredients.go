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
	_ types.ValidIngredientDataManager = (*Querier)(nil)
)

// ValidIngredientExists fetches whether a valid ingredient exists from the database.
func (q *Querier) ValidIngredientExists(ctx context.Context, validIngredientID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	result, err := q.generatedQuerier.CheckValidIngredientExistence(ctx, q.db, validIngredientID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient existence check")
	}

	return result, nil
}

// GetValidIngredient fetches a valid ingredient from the database.
func (q *Querier) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	result, err := q.generatedQuerier.GetValidIngredient(ctx, q.db, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	validIngredient := &types.ValidIngredient{
		CreatedAt:                               result.CreatedAt,
		LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
		MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
		IconPath:                                result.IconPath,
		Warning:                                 result.Warning,
		PluralName:                              result.PluralName,
		StorageInstructions:                     result.StorageInstructions,
		Name:                                    result.Name,
		ID:                                      result.ID,
		Description:                             result.Description,
		Slug:                                    result.Slug,
		ShoppingSuggestions:                     result.ShoppingSuggestions,
		ContainsShellfish:                       result.ContainsShellfish,
		IsMeasuredVolumetrically:                result.Volumetric,
		IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
		ContainsPeanut:                          result.ContainsPeanut,
		ContainsTreeNut:                         result.ContainsTreeNut,
		ContainsEgg:                             result.ContainsEgg,
		ContainsWheat:                           result.ContainsWheat,
		ContainsSoy:                             result.ContainsSoy,
		AnimalDerived:                           result.AnimalDerived,
		RestrictToPreparations:                  result.RestrictToPreparations,
		ContainsSesame:                          result.ContainsSesame,
		ContainsFish:                            result.ContainsFish,
		ContainsGluten:                          result.ContainsGluten,
		ContainsDairy:                           result.ContainsDairy,
		ContainsAlcohol:                         result.ContainsAlcohol,
		AnimalFlesh:                             result.AnimalFlesh,
		IsStarch:                                result.IsStarch,
		IsProtein:                               result.IsProtein,
		IsGrain:                                 result.IsGrain,
		IsFruit:                                 result.IsFruit,
		IsSalt:                                  result.IsSalt,
		IsFat:                                   result.IsFat,
		IsAcid:                                  result.IsAcid,
		IsHeat:                                  result.IsHeat,
	}

	return validIngredient, nil
}

// GetRandomValidIngredient fetches a valid ingredient from the database.
func (q *Querier) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidIngredient(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching random valid ingredient")
	}

	validIngredient := &types.ValidIngredient{
		CreatedAt:                               result.CreatedAt,
		LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
		MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
		IconPath:                                result.IconPath,
		Warning:                                 result.Warning,
		PluralName:                              result.PluralName,
		StorageInstructions:                     result.StorageInstructions,
		Name:                                    result.Name,
		ID:                                      result.ID,
		Description:                             result.Description,
		Slug:                                    result.Slug,
		ShoppingSuggestions:                     result.ShoppingSuggestions,
		ContainsShellfish:                       result.ContainsShellfish,
		IsMeasuredVolumetrically:                result.Volumetric,
		IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
		ContainsPeanut:                          result.ContainsPeanut,
		ContainsTreeNut:                         result.ContainsTreeNut,
		ContainsEgg:                             result.ContainsEgg,
		ContainsWheat:                           result.ContainsWheat,
		ContainsSoy:                             result.ContainsSoy,
		AnimalDerived:                           result.AnimalDerived,
		RestrictToPreparations:                  result.RestrictToPreparations,
		ContainsSesame:                          result.ContainsSesame,
		ContainsFish:                            result.ContainsFish,
		ContainsGluten:                          result.ContainsGluten,
		ContainsDairy:                           result.ContainsDairy,
		ContainsAlcohol:                         result.ContainsAlcohol,
		AnimalFlesh:                             result.AnimalFlesh,
		IsStarch:                                result.IsStarch,
		IsProtein:                               result.IsProtein,
		IsGrain:                                 result.IsGrain,
		IsFruit:                                 result.IsFruit,
		IsSalt:                                  result.IsSalt,
		IsFat:                                   result.IsFat,
		IsAcid:                                  result.IsAcid,
		IsHeat:                                  result.IsHeat,
	}

	return validIngredient, nil
}

// SearchForValidIngredients fetches a valid ingredient from the database.
func (q *Querier) SearchForValidIngredients(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, query)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	x := &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.SearchForValidIngredients(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient")
	}

	for _, result := range results {
		validIngredient := &types.ValidIngredient{
			CreatedAt:                               result.CreatedAt,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.IconPath,
			Warning:                                 result.Warning,
			PluralName:                              result.PluralName,
			StorageInstructions:                     result.StorageInstructions,
			Name:                                    result.Name,
			ID:                                      result.ID,
			Description:                             result.Description,
			Slug:                                    result.Slug,
			ShoppingSuggestions:                     result.ShoppingSuggestions,
			ContainsShellfish:                       result.ContainsShellfish,
			IsMeasuredVolumetrically:                result.Volumetric,
			IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:                          result.ContainsPeanut,
			ContainsTreeNut:                         result.ContainsTreeNut,
			ContainsEgg:                             result.ContainsEgg,
			ContainsWheat:                           result.ContainsWheat,
			ContainsSoy:                             result.ContainsSoy,
			AnimalDerived:                           result.AnimalDerived,
			RestrictToPreparations:                  result.RestrictToPreparations,
			ContainsSesame:                          result.ContainsSesame,
			ContainsFish:                            result.ContainsFish,
			ContainsGluten:                          result.ContainsGluten,
			ContainsDairy:                           result.ContainsDairy,
			ContainsAlcohol:                         result.ContainsAlcohol,
			AnimalFlesh:                             result.AnimalFlesh,
			IsStarch:                                result.IsStarch,
			IsProtein:                               result.IsProtein,
			IsGrain:                                 result.IsGrain,
			IsFruit:                                 result.IsFruit,
			IsSalt:                                  result.IsSalt,
			IsFat:                                   result.IsFat,
			IsAcid:                                  result.IsAcid,
			IsHeat:                                  result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
	}

	return x, nil
}

// SearchForValidIngredientsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, preparationID)

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	x = &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.SearchValidIngredientsByPreparationAndIngredientName(ctx, q.db, &generated.SearchValidIngredientsByPreparationAndIngredientNameParams{
		ValidPreparationID: preparationID,
		NameQuery:          query,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredients for preparation")
	}

	for _, result := range results {
		validIngredient := &types.ValidIngredient{
			CreatedAt:                               result.CreatedAt,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.IconPath,
			Warning:                                 result.Warning,
			PluralName:                              result.PluralName,
			StorageInstructions:                     result.StorageInstructions,
			Name:                                    result.Name,
			ID:                                      result.ID,
			Description:                             result.Description,
			Slug:                                    result.Slug,
			ShoppingSuggestions:                     result.ShoppingSuggestions,
			ContainsShellfish:                       result.ContainsShellfish,
			IsMeasuredVolumetrically:                result.Volumetric,
			IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:                          result.ContainsPeanut,
			ContainsTreeNut:                         result.ContainsTreeNut,
			ContainsEgg:                             result.ContainsEgg,
			ContainsWheat:                           result.ContainsWheat,
			ContainsSoy:                             result.ContainsSoy,
			AnimalDerived:                           result.AnimalDerived,
			RestrictToPreparations:                  result.RestrictToPreparations,
			ContainsSesame:                          result.ContainsSesame,
			ContainsFish:                            result.ContainsFish,
			ContainsGluten:                          result.ContainsGluten,
			ContainsDairy:                           result.ContainsDairy,
			ContainsAlcohol:                         result.ContainsAlcohol,
			AnimalFlesh:                             result.AnimalFlesh,
			IsStarch:                                result.IsStarch,
			IsProtein:                               result.IsProtein,
			IsGrain:                                 result.IsGrain,
			IsFruit:                                 result.IsFruit,
			IsSalt:                                  result.IsSalt,
			IsFat:                                   result.IsFat,
			IsAcid:                                  result.IsAcid,
			IsHeat:                                  result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
	}

	return x, nil
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredient], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredients(ctx, q.db, &generated.GetValidIngredientsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	for _, result := range results {
		validIngredient := &types.ValidIngredient{
			CreatedAt:                               result.CreatedAt,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.IconPath,
			Warning:                                 result.Warning,
			PluralName:                              result.PluralName,
			StorageInstructions:                     result.StorageInstructions,
			Name:                                    result.Name,
			ID:                                      result.ID,
			Description:                             result.Description,
			Slug:                                    result.Slug,
			ShoppingSuggestions:                     result.ShoppingSuggestions,
			ContainsShellfish:                       result.ContainsShellfish,
			IsMeasuredVolumetrically:                result.Volumetric,
			IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:                          result.ContainsPeanut,
			ContainsTreeNut:                         result.ContainsTreeNut,
			ContainsEgg:                             result.ContainsEgg,
			ContainsWheat:                           result.ContainsWheat,
			ContainsSoy:                             result.ContainsSoy,
			AnimalDerived:                           result.AnimalDerived,
			RestrictToPreparations:                  result.RestrictToPreparations,
			ContainsSesame:                          result.ContainsSesame,
			ContainsFish:                            result.ContainsFish,
			ContainsGluten:                          result.ContainsGluten,
			ContainsDairy:                           result.ContainsDairy,
			ContainsAlcohol:                         result.ContainsAlcohol,
			AnimalFlesh:                             result.AnimalFlesh,
			IsStarch:                                result.IsStarch,
			IsProtein:                               result.IsProtein,
			IsGrain:                                 result.IsGrain,
			IsFruit:                                 result.IsFruit,
			IsSalt:                                  result.IsSalt,
			IsFat:                                   result.IsFat,
			IsAcid:                                  result.IsAcid,
			IsHeat:                                  result.IsHeat,
		}

		x.Data = append(x.Data, validIngredient)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientsWithIDs fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ids == nil {
		return nil, ErrEmptyInputProvided
	}

	results, err := q.generatedQuerier.GetValidIngredientsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients id list retrieval query")
	}

	var ingredients []*types.ValidIngredient
	for _, result := range results {
		validIngredient := &types.ValidIngredient{
			CreatedAt:                               result.CreatedAt,
			LastUpdatedAt:                           database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                              database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MaximumIdealStorageTemperatureInCelsius),
			MinimumIdealStorageTemperatureInCelsius: database.Float32PointerFromNullString(result.MinimumIdealStorageTemperatureInCelsius),
			IconPath:                                result.IconPath,
			Warning:                                 result.Warning,
			PluralName:                              result.PluralName,
			StorageInstructions:                     result.StorageInstructions,
			Name:                                    result.Name,
			ID:                                      result.ID,
			Description:                             result.Description,
			Slug:                                    result.Slug,
			ShoppingSuggestions:                     result.ShoppingSuggestions,
			ContainsShellfish:                       result.ContainsShellfish,
			IsMeasuredVolumetrically:                result.Volumetric,
			IsLiquid:                                database.BoolFromNullBool(result.IsLiquid),
			ContainsPeanut:                          result.ContainsPeanut,
			ContainsTreeNut:                         result.ContainsTreeNut,
			ContainsEgg:                             result.ContainsEgg,
			ContainsWheat:                           result.ContainsWheat,
			ContainsSoy:                             result.ContainsSoy,
			AnimalDerived:                           result.AnimalDerived,
			RestrictToPreparations:                  result.RestrictToPreparations,
			ContainsSesame:                          result.ContainsSesame,
			ContainsFish:                            result.ContainsFish,
			ContainsGluten:                          result.ContainsGluten,
			ContainsDairy:                           result.ContainsDairy,
			ContainsAlcohol:                         result.ContainsAlcohol,
			AnimalFlesh:                             result.AnimalFlesh,
			IsStarch:                                result.IsStarch,
			IsProtein:                               result.IsProtein,
			IsGrain:                                 result.IsGrain,
			IsFruit:                                 result.IsFruit,
			IsSalt:                                  result.IsSalt,
			IsFat:                                   result.IsFat,
			IsAcid:                                  result.IsAcid,
			IsHeat:                                  result.IsHeat,
		}

		ingredients = append(ingredients, validIngredient)
	}

	return ingredients, nil
}

// GetValidIngredientIDsThatNeedSearchIndexing fetches a list of valid ingredients from the database that meet a particular filter.
func (q *Querier) GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidIngredientsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredients list retrieval query")
	}

	return results, err
}

// CreateValidIngredient creates a valid ingredient in the database.
func (q *Querier) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientIDKey, input.ID)

	// create the valid ingredient.
	if err := q.generatedQuerier.CreateValidIngredient(ctx, q.db, &generated.CreateValidIngredientParams{
		ID:                                      input.ID,
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		Volumetric:                              input.IsMeasuredVolumetrically,
		IsLiquid:                                database.NullBoolFromBool(input.IsLiquid),
		IconPath:                                input.IconPath,
		AnimalDerived:                           input.AnimalDerived,
		PluralName:                              input.PluralName,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MaximumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(input.MinimumIdealStorageTemperatureInCelsius),
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
		IsStarch:                                input.IsStarch,
		IsProtein:                               input.IsProtein,
		IsGrain:                                 input.IsGrain,
		IsFruit:                                 input.IsFruit,
		IsSalt:                                  input.IsSalt,
		IsFat:                                   input.IsFat,
		IsAcid:                                  input.IsAcid,
		IsHeat:                                  input.IsHeat,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient creation query")
	}

	x := &types.ValidIngredient{
		ID:                                      input.ID,
		Name:                                    input.Name,
		Description:                             input.Description,
		Warning:                                 input.Warning,
		ContainsEgg:                             input.ContainsEgg,
		ContainsDairy:                           input.ContainsDairy,
		ContainsPeanut:                          input.ContainsPeanut,
		ContainsTreeNut:                         input.ContainsTreeNut,
		ContainsSoy:                             input.ContainsSoy,
		ContainsWheat:                           input.ContainsWheat,
		ContainsShellfish:                       input.ContainsShellfish,
		ContainsSesame:                          input.ContainsSesame,
		ContainsFish:                            input.ContainsFish,
		ContainsGluten:                          input.ContainsGluten,
		AnimalFlesh:                             input.AnimalFlesh,
		IsMeasuredVolumetrically:                input.IsMeasuredVolumetrically,
		IsLiquid:                                input.IsLiquid,
		IconPath:                                input.IconPath,
		AnimalDerived:                           input.AnimalDerived,
		PluralName:                              input.PluralName,
		IsStarch:                                input.IsStarch,
		IsProtein:                               input.IsProtein,
		IsGrain:                                 input.IsGrain,
		IsFruit:                                 input.IsFruit,
		IsSalt:                                  input.IsSalt,
		IsFat:                                   input.IsFat,
		IsAcid:                                  input.IsAcid,
		IsHeat:                                  input.IsHeat,
		RestrictToPreparations:                  input.RestrictToPreparations,
		MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
		MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
		StorageInstructions:                     input.StorageInstructions,
		Slug:                                    input.Slug,
		ContainsAlcohol:                         input.ContainsAlcohol,
		ShoppingSuggestions:                     input.ShoppingSuggestions,
		CreatedAt:                               q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, x.ID)
	logger.Info("valid ingredient created")

	return x, nil
}

// UpdateValidIngredient updates a particular valid ingredient.
func (q *Querier) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredient(ctx, q.db, &generated.UpdateValidIngredientParams{
		Description:                             updated.Description,
		Warning:                                 updated.Warning,
		ID:                                      updated.ID,
		ShoppingSuggestions:                     updated.ShoppingSuggestions,
		Slug:                                    updated.Slug,
		StorageInstructions:                     updated.StorageInstructions,
		Name:                                    updated.Name,
		PluralName:                              updated.PluralName,
		IconPath:                                updated.IconPath,
		MaximumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.MaximumIdealStorageTemperatureInCelsius),
		MinimumIdealStorageTemperatureInCelsius: database.NullStringFromFloat32Pointer(updated.MinimumIdealStorageTemperatureInCelsius),
		IsLiquid:                                database.NullBoolFromBool(updated.IsLiquid),
		ContainsWheat:                           updated.ContainsWheat,
		ContainsPeanut:                          updated.ContainsPeanut,
		Volumetric:                              updated.IsMeasuredVolumetrically,
		ContainsGluten:                          updated.ContainsGluten,
		ContainsFish:                            updated.ContainsFish,
		AnimalDerived:                           updated.AnimalDerived,
		ContainsSesame:                          updated.ContainsSesame,
		RestrictToPreparations:                  updated.RestrictToPreparations,
		ContainsShellfish:                       updated.ContainsShellfish,
		ContainsSoy:                             updated.ContainsSoy,
		ContainsTreeNut:                         updated.ContainsTreeNut,
		AnimalFlesh:                             updated.AnimalFlesh,
		ContainsAlcohol:                         updated.ContainsAlcohol,
		ContainsDairy:                           updated.ContainsDairy,
		IsStarch:                                updated.IsStarch,
		IsProtein:                               updated.IsProtein,
		IsGrain:                                 updated.IsGrain,
		IsFruit:                                 updated.IsFruit,
		IsSalt:                                  updated.IsSalt,
		IsFat:                                   updated.IsFat,
		IsAcid:                                  updated.IsAcid,
		IsHeat:                                  updated.IsHeat,
		ContainsEgg:                             updated.ContainsEgg,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient")
	}

	logger.Info("valid ingredient updated")

	return nil
}

// MarkValidIngredientAsIndexed updates a particular valid ingredient's last_indexed_at value.
func (q *Querier) MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	if _, err := q.generatedQuerier.UpdateValidIngredientLastIndexedAt(ctx, q.db, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid ingredient as indexed")
	}

	logger.Info("valid ingredient marked as indexed")

	return nil
}

// ArchiveValidIngredient archives a valid ingredient from the database by its ID.
func (q *Querier) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	if _, err := q.generatedQuerier.ArchiveValidIngredient(ctx, q.db, validIngredientID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient")
	}

	logger.Info("valid ingredient archived")

	return nil
}
