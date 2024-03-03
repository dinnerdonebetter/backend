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
	_ types.RecipeStepVesselDataManager = (*Querier)(nil)
)

// RecipeStepVesselExists fetches whether a recipe step vessel exists from the database.
func (q *Querier) RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (exists bool, err error) {
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

	if recipeStepVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	result, err := q.generatedQuerier.CheckRecipeStepVesselExistence(ctx, q.db, &generated.CheckRecipeStepVesselExistenceParams{
		RecipeStepID:       recipeStepID,
		RecipeStepVesselID: recipeStepVesselID,
		RecipeID:           recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step vessel existence check")
	}

	return result, nil
}

// GetRecipeStepVessel fetches a recipe step vessel from the database.
func (q *Querier) GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepVesselID string) (*types.RecipeStepVessel, error) {
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

	if recipeStepVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	result, err := q.generatedQuerier.GetRecipeStepVessel(ctx, q.db, &generated.GetRecipeStepVesselParams{
		RecipeStepID:       recipeStepID,
		RecipeStepVesselID: recipeStepVesselID,
		RecipeID:           recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step vessel")
	}

	recipeStepVessel := &types.RecipeStepVessel{
		CreatedAt:            result.CreatedAt,
		MaximumQuantity:      database.Uint32PointerFromNullInt32(result.MaximumQuantity),
		LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
		RecipeStepProductID:  database.StringPointerFromNullString(result.RecipeStepProductID),
		Vessel:               nil,
		ID:                   result.ID,
		Notes:                result.Notes,
		BelongsToRecipeStep:  result.BelongsToRecipeStep,
		VesselPreposition:    result.VesselPredicate,
		Name:                 result.Name,
		MinimumQuantity:      uint32(result.MinimumQuantity),
		UnavailableAfterStep: result.UnavailableAfterStep,
	}

	if result.ValidVesselID.Valid {
		recipeStepVessel.Vessel = &types.ValidVessel{
			CreatedAt:                      result.ValidVesselCreatedAt.Time,
			ArchivedAt:                     database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
			CapacityUnit:                   nil,
			IconPath:                       result.ValidVesselIconPath.String,
			PluralName:                     result.ValidVesselPluralName.String,
			Description:                    result.ValidVesselDescription.String,
			Name:                           result.ValidVesselName.String,
			Slug:                           result.ValidVesselSlug.String,
			Shape:                          string(result.ValidVesselShape.VesselShape),
			ID:                             result.ValidVesselID.String,
			WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
			LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
			HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
			Capacity:                       database.Float32FromNullString(result.ValidVesselCapacity),
			IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions.Bool,
			DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists.Bool,
			UsableForStorage:               result.ValidVesselUsableForStorage.Bool,
		}

		if result.ValidMeasurementUnitID.Valid {
			recipeStepVessel.Vessel.CapacityUnit = &types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			}
		}
	}

	return recipeStepVessel, nil
}

// GetRecipeStepVessels fetches a list of recipe step vessels from the database that meet a particular filter.
func (q *Querier) GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepVessel], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepVessel]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepVessels(ctx, q.db, &generated.GetRecipeStepVesselsParams{
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
		return nil, observability.PrepareAndLogError(err, logger, span, "executing recipe step vessels list retrieval query")
	}

	for _, result := range results {
		recipeStepVessel := &types.RecipeStepVessel{
			CreatedAt:            result.CreatedAt,
			MaximumQuantity:      database.Uint32PointerFromNullInt32(result.MaximumQuantity),
			LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
			RecipeStepProductID:  database.StringPointerFromNullString(result.RecipeStepProductID),
			Vessel:               nil,
			ID:                   result.ID,
			Notes:                result.Notes,
			BelongsToRecipeStep:  result.BelongsToRecipeStep,
			VesselPreposition:    result.VesselPredicate,
			Name:                 result.Name,
			MinimumQuantity:      uint32(result.MinimumQuantity),
			UnavailableAfterStep: result.UnavailableAfterStep,
		}

		if result.ValidVesselID.Valid {
			recipeStepVessel.Vessel = &types.ValidVessel{
				CreatedAt:                      result.ValidVesselCreatedAt.Time,
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit:                   nil,
				IconPath:                       result.ValidVesselIconPath.String,
				PluralName:                     result.ValidVesselPluralName.String,
				Description:                    result.ValidVesselDescription.String,
				Name:                           result.ValidVesselName.String,
				Slug:                           result.ValidVesselSlug.String,
				Shape:                          string(result.ValidVesselShape.VesselShape),
				ID:                             result.ValidVesselID.String,
				WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
				Capacity:                       database.Float32FromNullString(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions.Bool,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists.Bool,
				UsableForStorage:               result.ValidVesselUsableForStorage.Bool,
			}

			if result.ValidMeasurementUnitID.Valid {
				recipeStepVessel.Vessel.CapacityUnit = &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
					Name:          result.ValidMeasurementUnitName.String,
					IconPath:      result.ValidMeasurementUnitIconPath.String,
					ID:            result.ValidMeasurementUnitID.String,
					Description:   result.ValidMeasurementUnitDescription.String,
					PluralName:    result.ValidMeasurementUnitPluralName.String,
					Slug:          result.ValidMeasurementUnitSlug.String,
					Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
					Universal:     result.ValidMeasurementUnitUniversal.Bool,
					Metric:        result.ValidMeasurementUnitMetric.Bool,
					Imperial:      result.ValidMeasurementUnitImperial.Bool,
				}
			}
		}

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
		x.Data = append(x.Data, recipeStepVessel)
	}

	return x, nil
}

// getRecipeStepVesselsForRecipe fetches a list of recipe step vessels from the database that meet a particular filter.
func (q *Querier) getRecipeStepVesselsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.GetRecipeStepVesselsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting recipe step vessels for a recipe")
	}

	recipeStepVessels := []*types.RecipeStepVessel{}
	for _, result := range results {
		recipeStepVessel := &types.RecipeStepVessel{
			CreatedAt:            result.CreatedAt,
			MaximumQuantity:      database.Uint32PointerFromNullInt32(result.MaximumQuantity),
			LastUpdatedAt:        database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:           database.TimePointerFromNullTime(result.ArchivedAt),
			RecipeStepProductID:  database.StringPointerFromNullString(result.RecipeStepProductID),
			Vessel:               nil,
			ID:                   result.ID,
			Notes:                result.Notes,
			BelongsToRecipeStep:  result.BelongsToRecipeStep,
			VesselPreposition:    result.VesselPredicate,
			Name:                 result.Name,
			MinimumQuantity:      uint32(result.MinimumQuantity),
			UnavailableAfterStep: result.UnavailableAfterStep,
		}

		if result.ValidVesselID.Valid {
			recipeStepVessel.Vessel = &types.ValidVessel{
				CreatedAt:                      result.ValidVesselCreatedAt.Time,
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit:                   nil,
				IconPath:                       result.ValidVesselIconPath.String,
				PluralName:                     result.ValidVesselPluralName.String,
				Description:                    result.ValidVesselDescription.String,
				Name:                           result.ValidVesselName.String,
				Slug:                           result.ValidVesselSlug.String,
				Shape:                          string(result.ValidVesselShape.VesselShape),
				ID:                             result.ValidVesselID.String,
				WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
				Capacity:                       database.Float32FromNullString(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions.Bool,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists.Bool,
				UsableForStorage:               result.ValidVesselUsableForStorage.Bool,
			}

			if result.ValidMeasurementUnitID.Valid {
				recipeStepVessel.Vessel.CapacityUnit = &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
					Name:          result.ValidMeasurementUnitName.String,
					IconPath:      result.ValidMeasurementUnitIconPath.String,
					ID:            result.ValidMeasurementUnitID.String,
					Description:   result.ValidMeasurementUnitDescription.String,
					PluralName:    result.ValidMeasurementUnitPluralName.String,
					Slug:          result.ValidMeasurementUnitSlug.String,
					Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
					Universal:     result.ValidMeasurementUnitUniversal.Bool,
					Metric:        result.ValidMeasurementUnitMetric.Bool,
					Imperial:      result.ValidMeasurementUnitImperial.Bool,
				}
			}
		}

		recipeStepVessels = append(recipeStepVessels, recipeStepVessel)
	}

	return recipeStepVessels, nil
}

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) createRecipeStepVessel(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, input.ID)

	// create the recipe step vessel.
	if err := q.generatedQuerier.CreateRecipeStepVessel(ctx, querier, &generated.CreateRecipeStepVesselParams{
		ID:                   input.ID,
		Name:                 input.Name,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		VesselPredicate:      input.VesselPreposition,
		RecipeStepProductID:  database.NullStringFromStringPointer(input.RecipeStepProductID),
		ValidVesselID:        database.NullStringFromStringPointer(input.VesselID),
		MaximumQuantity:      database.NullInt32FromUint32Pointer(input.MaximumQuantity),
		MinimumQuantity:      int32(input.MinimumQuantity),
		UnavailableAfterStep: input.UnavailableAfterStep,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step vessel creation query")
	}

	x := &types.RecipeStepVessel{
		ID:                   input.ID,
		RecipeStepProductID:  input.RecipeStepProductID,
		Name:                 input.Name,
		Notes:                input.Notes,
		BelongsToRecipeStep:  input.BelongsToRecipeStep,
		MinimumQuantity:      input.MinimumQuantity,
		MaximumQuantity:      input.MaximumQuantity,
		VesselPreposition:    input.VesselPreposition,
		UnavailableAfterStep: input.UnavailableAfterStep,
		CreatedAt:            q.currentTime(),
	}

	if input.VesselID != nil {
		x.Vessel = &types.ValidVessel{ID: *input.VesselID}
	}

	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, x.ID)
	logger.Info("recipe step vessel created")

	return x, nil
}

// CreateRecipeStepVessel creates a recipe step vessel in the database.
func (q *Querier) CreateRecipeStepVessel(ctx context.Context, input *types.RecipeStepVesselDatabaseCreationInput) (*types.RecipeStepVessel, error) {
	return q.createRecipeStepVessel(ctx, q.db, input)
}

// UpdateRecipeStepVessel updates a particular recipe step vessel.
func (q *Querier) UpdateRecipeStepVessel(ctx context.Context, updated *types.RecipeStepVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepVesselIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, updated.ID)

	var vesselID *string
	if updated.Vessel != nil {
		vesselID = &updated.Vessel.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepVessel(ctx, q.db, &generated.UpdateRecipeStepVesselParams{
		Name:                 updated.Name,
		Notes:                updated.Notes,
		BelongsToRecipeStep:  updated.BelongsToRecipeStep,
		VesselPredicate:      updated.VesselPreposition,
		ID:                   updated.ID,
		RecipeStepProductID:  database.NullStringFromStringPointer(updated.RecipeStepProductID),
		ValidVesselID:        database.NullStringFromStringPointer(vesselID),
		MaximumQuantity:      database.NullInt32FromUint32Pointer(updated.MaximumQuantity),
		MinimumQuantity:      int32(updated.MinimumQuantity),
		UnavailableAfterStep: updated.UnavailableAfterStep,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step vessel")
	}

	logger.Info("recipe step vessel updated")

	return nil
}

// ArchiveRecipeStepVessel archives a recipe step vessel from the database by its ID.
func (q *Querier) ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepVesselIDKey, recipeStepVesselID)
	tracing.AttachToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepVessel(ctx, q.db, &generated.ArchiveRecipeStepVesselParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepVesselID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step vessel")
	}

	logger.Info("recipe step vessel archived")

	return nil
}
