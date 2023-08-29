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
	_ types.RecipeStepProductDataManager = (*Querier)(nil)
)

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *Querier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	result, err := q.generatedQuerier.CheckRecipeStepProductExistence(ctx, q.db, &generated.CheckRecipeStepProductExistenceParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

// GetRecipeStepProduct fetches a recipe step product from the database.
func (q *Querier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	result, err := q.generatedQuerier.GetRecipeStepProduct(ctx, q.db, &generated.GetRecipeStepProductParams{
		RecipeStepID:        recipeStepID,
		RecipeStepProductID: recipeStepProductID,
		RecipeID:            recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step product")
	}

	recipeStepProduct := &types.RecipeStepProduct{
		CreatedAt:                          result.CreatedAt,
		MaximumStorageTemperatureInCelsius: float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
		MaximumStorageDurationInSeconds:    uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
		MinimumStorageTemperatureInCelsius: float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
		ArchivedAt:                         timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:                      timePointerFromNullTime(result.LastUpdatedAt),
		MinimumQuantity:                    float32PointerFromNullString(result.MinimumQuantityValue),
		MeasurementUnit:                    nil,
		MaximumQuantity:                    float32PointerFromNullString(result.MaximumQuantityValue),
		ContainedInVesselIndex:             uint16PointerFromNullInt32(result.ContainedInVesselIndex),
		Name:                               result.Name,
		BelongsToRecipeStep:                result.BelongsToRecipeStep,
		Type:                               string(result.Type),
		ID:                                 result.ID,
		StorageInstructions:                result.StorageInstructions,
		QuantityNotes:                      result.QuantityNotes,
		Index:                              uint16(result.Index),
		IsWaste:                            result.IsWaste,
		IsLiquid:                           result.IsLiquid,
		Compostable:                        result.Compostable,
	}

	if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
		recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
			LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName.String,
			IconPath:      result.ValidMeasurementUnitIconPath.String,
			ID:            result.ValidMeasurementUnitID.String,
			Description:   result.ValidMeasurementUnitDescription.String,
			PluralName:    result.ValidMeasurementUnitPluralName.String,
			Slug:          result.ValidMeasurementUnitSlug.String,
			Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal.Bool,
			Metric:        result.ValidMeasurementUnitMetric.Bool,
			Imperial:      result.ValidMeasurementUnitImperial.Bool,
		}
	}

	return recipeStepProduct, nil
}

// getRecipeStepProductsForRecipe fetches a list of recipe step products from the database that meet a particular filter.
func (q *Querier) getRecipeStepProductsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	results, err := q.generatedQuerier.GetRecipeStepProductsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step products for recipe")
	}

	recipeStepProducts := []*types.RecipeStepProduct{}
	for _, result := range results {
		recipeStepProduct := &types.RecipeStepProduct{
			CreatedAt:                          result.CreatedAt,
			MaximumStorageTemperatureInCelsius: float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
			MaximumStorageDurationInSeconds:    uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
			MinimumStorageTemperatureInCelsius: float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
			ArchivedAt:                         timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:                      timePointerFromNullTime(result.LastUpdatedAt),
			MinimumQuantity:                    float32PointerFromNullString(result.MinimumQuantityValue),
			MeasurementUnit:                    nil,
			MaximumQuantity:                    float32PointerFromNullString(result.MaximumQuantityValue),
			ContainedInVesselIndex:             uint16PointerFromNullInt32(result.ContainedInVesselIndex),
			Name:                               result.Name,
			BelongsToRecipeStep:                result.BelongsToRecipeStep,
			Type:                               string(result.Type),
			ID:                                 result.ID,
			StorageInstructions:                result.StorageInstructions,
			QuantityNotes:                      result.QuantityNotes,
			Index:                              uint16(result.Index),
			IsWaste:                            result.IsWaste,
			IsLiquid:                           result.IsLiquid,
			Compostable:                        result.Compostable,
		}

		if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
			recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			}
		}

		recipeStepProducts = append(recipeStepProducts, recipeStepProduct)
	}

	return recipeStepProducts, nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (q *Querier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepProduct], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.RecipeStepProduct]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepProducts(ctx, q.db, &generated.GetRecipeStepProductsParams{
		RecipeStepID:  recipeStepID,
		RecipeID:      recipeID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step products")
	}

	for _, result := range results {
		recipeStepProduct := &types.RecipeStepProduct{
			CreatedAt:                          result.CreatedAt,
			MaximumStorageTemperatureInCelsius: float32PointerFromNullString(result.MaximumStorageTemperatureInCelsius),
			MaximumStorageDurationInSeconds:    uint32PointerFromNullInt32(result.MaximumStorageDurationInSeconds),
			MinimumStorageTemperatureInCelsius: float32PointerFromNullString(result.MinimumStorageTemperatureInCelsius),
			ArchivedAt:                         timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:                      timePointerFromNullTime(result.LastUpdatedAt),
			MinimumQuantity:                    float32PointerFromNullString(result.MinimumQuantityValue),
			MeasurementUnit:                    nil,
			MaximumQuantity:                    float32PointerFromNullString(result.MaximumQuantityValue),
			ContainedInVesselIndex:             uint16PointerFromNullInt32(result.ContainedInVesselIndex),
			Name:                               result.Name,
			BelongsToRecipeStep:                result.BelongsToRecipeStep,
			Type:                               string(result.Type),
			ID:                                 result.ID,
			StorageInstructions:                result.StorageInstructions,
			QuantityNotes:                      result.QuantityNotes,
			Index:                              uint16(result.Index),
			IsWaste:                            result.IsWaste,
			IsLiquid:                           result.IsLiquid,
			Compostable:                        result.Compostable,
		}

		if result.ValidMeasurementUnitID.Valid && result.ValidMeasurementUnitID.String != "" {
			recipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    boolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepProduct)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// create the recipe step product.
	if err := q.generatedQuerier.CreateRecipeStepProduct(ctx, db, &generated.CreateRecipeStepProductParams{
		QuantityNotes:                      input.QuantityNotes,
		Name:                               input.Name,
		Type:                               generated.RecipeStepProductType(input.Type),
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		ID:                                 input.ID,
		StorageInstructions:                input.StorageInstructions,
		MinimumQuantityValue:               nullStringFromFloat32Pointer(input.MinimumQuantity),
		MinimumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(input.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(input.MaximumStorageTemperatureInCelsius),
		MaximumQuantityValue:               nullStringFromFloat32Pointer(input.MaximumQuantity),
		MeasurementUnit:                    nullStringFromStringPointer(input.MeasurementUnitID),
		MaximumStorageDurationInSeconds:    nullInt32FromUint32Pointer(input.MaximumStorageDurationInSeconds),
		ContainedInVesselIndex:             nullInt32FromUint16Pointer(input.ContainedInVesselIndex),
		Index:                              int32(input.Index),
		Compostable:                        input.Compostable,
		IsLiquid:                           input.IsLiquid,
		IsWaste:                            input.IsWaste,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "performing recipe step product creation query")
	}

	x := &types.RecipeStepProduct{
		ID:                                 input.ID,
		Name:                               input.Name,
		Type:                               input.Type,
		MinimumQuantity:                    input.MinimumQuantity,
		MaximumQuantity:                    input.MaximumQuantity,
		QuantityNotes:                      input.QuantityNotes,
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		StorageInstructions:                input.StorageInstructions,
		IsLiquid:                           input.IsLiquid,
		IsWaste:                            input.IsWaste,
		Index:                              input.Index,
		ContainedInVesselIndex:             input.ContainedInVesselIndex,
		CreatedAt:                          q.currentTime(),
	}

	if input.MeasurementUnitID != nil {
		x.MeasurementUnit = &types.ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	tracing.AttachRecipeStepProductIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *Querier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	return q.createRecipeStepProduct(ctx, q.db, input)
}

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *Querier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)

	var measurementUnitID *string
	if updated.MeasurementUnit != nil {
		measurementUnitID = &updated.MeasurementUnit.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepProduct(ctx, q.db, &generated.UpdateRecipeStepProductParams{
		Name:                               updated.Name,
		Type:                               generated.RecipeStepProductType(updated.Type),
		MeasurementUnit:                    nullStringFromStringPointer(measurementUnitID),
		MinimumQuantityValue:               nullStringFromFloat32Pointer(updated.MinimumQuantity),
		MaximumQuantityValue:               nullStringFromFloat32Pointer(updated.MaximumQuantity),
		QuantityNotes:                      updated.QuantityNotes,
		Compostable:                        updated.Compostable,
		MaximumStorageDurationInSeconds:    nullInt32FromUint32Pointer(updated.MaximumStorageDurationInSeconds),
		MinimumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(updated.MinimumStorageTemperatureInCelsius),
		MaximumStorageTemperatureInCelsius: nullStringFromFloat32Pointer(updated.MaximumStorageTemperatureInCelsius),
		StorageInstructions:                updated.StorageInstructions,
		IsLiquid:                           updated.IsLiquid,
		IsWaste:                            updated.IsWaste,
		Index:                              int32(updated.Index),
		ContainedInVesselIndex:             nullInt32FromUint16Pointer(updated.ContainedInVesselIndex),
		BelongsToRecipeStep:                updated.BelongsToRecipeStep,
		ID:                                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (q *Querier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepProduct(ctx, q.db, &generated.ArchiveRecipeStepProductParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepProductID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}
