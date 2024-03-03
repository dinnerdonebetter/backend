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
	_ types.ValidPreparationInstrumentDataManager = (*Querier)(nil)
)

// ValidPreparationInstrumentExists fetches whether a valid preparation instrument exists from the database.
func (q *Querier) ValidPreparationInstrumentExists(ctx context.Context, validPreparationInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	result, err := q.generatedQuerier.CheckValidPreparationInstrumentExistence(ctx, q.db, validPreparationInstrumentID)
	if err != nil {
		return false, observability.PrepareError(err, span, "checking valid preparation instrument existence")
	}

	return result, nil
}

// GetValidPreparationInstrument fetches a valid preparation instrument from the database.
func (q *Querier) GetValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	result, err := q.generatedQuerier.GetValidPreparationInstrument(ctx, q.db, validPreparationInstrumentID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting valid preparation instrument")
	}

	x := &types.ValidPreparationInstrument{
		CreatedAt:     result.ValidPreparationInstrumentCreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationInstrumentLastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationInstrumentArchivedAt),
		ID:            result.ValidPreparationInstrumentID,
		Notes:         result.ValidPreparationInstrumentNotes,
		Instrument: types.ValidInstrument{
			CreatedAt:                      result.ValidInstrumentCreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
			IconPath:                       result.ValidInstrumentIconPath,
			ID:                             result.ValidInstrumentID,
			Name:                           result.ValidInstrumentName,
			PluralName:                     result.ValidInstrumentPluralName,
			Description:                    result.ValidInstrumentDescription,
			Slug:                           result.ValidInstrumentSlug,
			DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
			UsableForStorage:               result.ValidInstrumentUsableForStorage,
		},
		Preparation: types.ValidPreparation{
			CreatedAt:                   result.ValidPreparationCreatedAt,
			MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
			ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
			MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
			MaximumVesselCount:          database.Int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
			IconPath:                    result.ValidPreparationIconPath,
			PastTense:                   result.ValidPreparationPastTense,
			ID:                          result.ValidPreparationID,
			Name:                        result.ValidPreparationName,
			Description:                 result.ValidPreparationDescription,
			Slug:                        result.ValidPreparationSlug,
			MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
			MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
			MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
			RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
			TemperatureRequired:         result.ValidPreparationTemperatureRequired,
			TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
			ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
			ConsumesVessel:              result.ValidPreparationConsumesVessel,
			OnlyForVessels:              result.ValidPreparationOnlyForVessels,
			YieldsNothing:               result.ValidPreparationYieldsNothing,
		},
	}

	return x, nil
}

// GetValidPreparationInstruments fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidPreparationInstruments(ctx, q.db, &generated.GetValidPreparationInstrumentsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	for _, result := range results {
		validPreparationInstrument := &types.ValidPreparationInstrument{
			CreatedAt:     result.ValidPreparationInstrumentCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationInstrumentLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationInstrumentArchivedAt),
			ID:            result.ValidPreparationInstrumentID,
			Notes:         result.ValidPreparationInstrumentNotes,
			Instrument: types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath,
				ID:                             result.ValidInstrumentID,
				Name:                           result.ValidInstrumentName,
				PluralName:                     result.ValidInstrumentPluralName,
				Description:                    result.ValidInstrumentDescription,
				Slug:                           result.ValidInstrumentSlug,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
				UsableForStorage:               result.ValidInstrumentUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          database.Int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		}

		x.Data = append(x.Data, validPreparationInstrument)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidPreparationInstrumentsForPreparation fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, preparationID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidPreparationInstrumentsForPreparation(ctx, q.db, &generated.GetValidPreparationInstrumentsForPreparationParams{
		ID:            preparationID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	for _, result := range results {
		validPreparationInstrument := &types.ValidPreparationInstrument{
			CreatedAt:     result.ValidPreparationInstrumentCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationInstrumentLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationInstrumentArchivedAt),
			ID:            result.ValidPreparationInstrumentID,
			Notes:         result.ValidPreparationInstrumentNotes,
			Instrument: types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath,
				ID:                             result.ValidInstrumentID,
				Name:                           result.ValidInstrumentName,
				PluralName:                     result.ValidInstrumentPluralName,
				Description:                    result.ValidInstrumentDescription,
				Slug:                           result.ValidInstrumentSlug,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
				UsableForStorage:               result.ValidInstrumentUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          database.Int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		}

		x.Data = append(x.Data, validPreparationInstrument)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidPreparationInstrumentsForInstrument fetches a list of valid preparation instruments from the database that meet a particular filter.
func (q *Querier) GetValidPreparationInstrumentsForInstrument(ctx context.Context, instrumentID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if instrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, instrumentID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationInstrument]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationInstrumentsForInstrument(ctx, q.db, &generated.GetValidPreparationInstrumentsForInstrumentParams{
		ID:            instrumentID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation instruments list retrieval query")
	}

	for _, result := range results {
		validPreparationInstrument := &types.ValidPreparationInstrument{
			CreatedAt:     result.ValidPreparationInstrumentCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationInstrumentLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationInstrumentArchivedAt),
			ID:            result.ValidPreparationInstrumentID,
			Notes:         result.ValidPreparationInstrumentNotes,
			Instrument: types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath,
				ID:                             result.ValidInstrumentID,
				Name:                           result.ValidInstrumentName,
				PluralName:                     result.ValidInstrumentPluralName,
				Description:                    result.ValidInstrumentDescription,
				Slug:                           result.ValidInstrumentSlug,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
				UsableForStorage:               result.ValidInstrumentUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          database.Int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		}

		x.Data = append(x.Data, validPreparationInstrument)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// CreateValidPreparationInstrument creates a valid preparation instrument in the database.
func (q *Querier) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentDatabaseCreationInput) (*types.ValidPreparationInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, input.ID)

	// create the valid preparation instrument.
	if err := q.generatedQuerier.CreateValidPreparationInstrument(ctx, q.db, &generated.CreateValidPreparationInstrumentParams{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidInstrumentID:  input.ValidInstrumentID,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation instrument creation query")
	}

	x := &types.ValidPreparationInstrument{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Instrument:  types.ValidInstrument{ID: input.ValidInstrumentID},
		CreatedAt:   q.currentTime(),
	}

	logger.Info("valid preparation instrument created")

	return x, nil
}

// UpdateValidPreparationInstrument updates a particular valid preparation instrument.
func (q *Querier) UpdateValidPreparationInstrument(ctx context.Context, updated *types.ValidPreparationInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidPreparationInstrument(ctx, q.db, &generated.UpdateValidPreparationInstrumentParams{
		Notes:              updated.Notes,
		ValidPreparationID: updated.Preparation.ID,
		ValidInstrumentID:  updated.Instrument.ID,
		ID:                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument updated")

	return nil
}

// ArchiveValidPreparationInstrument archives a valid preparation instrument from the database by its ID.
func (q *Querier) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	if _, err := q.generatedQuerier.ArchiveValidPreparationInstrument(ctx, q.db, validPreparationInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation instrument")
	}

	logger.Info("valid preparation instrument archived")

	return nil
}
