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
	_ types.RecipeStepInstrumentDataManager = (*Querier)(nil)
)

// RecipeStepInstrumentExists fetches whether a recipe step instrument exists from the database.
func (q *Querier) RecipeStepInstrumentExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (exists bool, err error) {
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

	if recipeStepInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	result, err := q.generatedQuerier.CheckRecipeStepInstrumentExistence(ctx, q.db, &generated.CheckRecipeStepInstrumentExistenceParams{
		RecipeStepID:           recipeStepID,
		RecipeStepInstrumentID: recipeStepInstrumentID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing recipe step instrument existence check")
	}

	return result, nil
}

// GetRecipeStepInstrument fetches a recipe step instrument from the database.
func (q *Querier) GetRecipeStepInstrument(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*types.RecipeStepInstrument, error) {
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

	if recipeStepInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	result, err := q.generatedQuerier.GetRecipeStepInstrument(ctx, q.db, &generated.GetRecipeStepInstrumentParams{
		RecipeStepID:           recipeStepID,
		RecipeStepInstrumentID: recipeStepInstrumentID,
		RecipeID:               recipeID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step instrument get")
	}

	recipeStepInstrument := &types.RecipeStepInstrument{
		CreatedAt:           result.CreatedAt,
		Instrument:          nil,
		LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
		RecipeStepProductID: database.StringPointerFromNullString(result.RecipeStepProductID),
		ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
		MaximumQuantity:     database.Uint32PointerFromNullInt32(result.MaximumQuantity),
		Notes:               result.Notes,
		Name:                result.Name,
		BelongsToRecipeStep: result.BelongsToRecipeStep,
		ID:                  result.ID,
		MinimumQuantity:     uint32(result.MinimumQuantity),
		OptionIndex:         uint16(result.OptionIndex),
		PreferenceRank:      uint8(result.PreferenceRank),
		Optional:            result.Optional,
	}

	if result.ValidInstrumentID.Valid {
		recipeStepInstrument.Instrument = &types.ValidInstrument{
			CreatedAt:                      result.ValidInstrumentCreatedAt.Time,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
			IconPath:                       result.ValidInstrumentIconPath.String,
			ID:                             result.ValidInstrumentID.String,
			Name:                           result.ValidInstrumentName.String,
			PluralName:                     result.ValidInstrumentPluralName.String,
			Description:                    result.ValidInstrumentDescription.String,
			Slug:                           result.ValidInstrumentSlug.String,
			DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists.Bool,
			IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions.Bool,
			UsableForStorage:               result.ValidInstrumentUsableForStorage.Bool,
		}
	}

	return recipeStepInstrument, nil
}

// GetRecipeStepInstruments fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *Querier) GetRecipeStepInstruments(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.RecipeStepInstrument], err error) {
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

	x = &types.QueryFilteredResult[types.RecipeStepInstrument]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetRecipeStepInstruments(ctx, q.db, &generated.GetRecipeStepInstrumentsParams{
		RecipeStepID:  recipeStepID,
		RecipeID:      recipeID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step instruments list retrieval")
	}

	for _, result := range results {
		recipeStepInstrument := &types.RecipeStepInstrument{
			CreatedAt:           result.CreatedAt,
			Instrument:          nil,
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			RecipeStepProductID: database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumQuantity:     database.Uint32PointerFromNullInt32(result.MaximumQuantity),
			Notes:               result.Notes,
			Name:                result.Name,
			BelongsToRecipeStep: result.BelongsToRecipeStep,
			ID:                  result.ID,
			MinimumQuantity:     uint32(result.MinimumQuantity),
			OptionIndex:         uint16(result.OptionIndex),
			PreferenceRank:      uint8(result.PreferenceRank),
			Optional:            result.Optional,
		}

		if result.ValidInstrumentID.Valid {
			recipeStepInstrument.Instrument = &types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt.Time,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath.String,
				ID:                             result.ValidInstrumentID.String,
				Name:                           result.ValidInstrumentName.String,
				PluralName:                     result.ValidInstrumentPluralName.String,
				Description:                    result.ValidInstrumentDescription.String,
				Slug:                           result.ValidInstrumentSlug.String,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists.Bool,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions.Bool,
				UsableForStorage:               result.ValidInstrumentUsableForStorage.Bool,
			}
		}

		x.Data = append(x.Data, recipeStepInstrument)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// getRecipeStepInstrumentsForRecipe fetches a list of recipe step instruments from the database that meet a particular filter.
func (q *Querier) getRecipeStepInstrumentsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	results, err := q.generatedQuerier.GetRecipeStepInstrumentsForRecipe(ctx, q.db, recipeID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step instruments list retrieval")
	}

	recipeStepInstruments := []*types.RecipeStepInstrument{}
	for _, result := range results {
		recipeStepInstrument := &types.RecipeStepInstrument{
			CreatedAt:           result.CreatedAt,
			Instrument:          nil,
			LastUpdatedAt:       database.TimePointerFromNullTime(result.LastUpdatedAt),
			RecipeStepProductID: database.StringPointerFromNullString(result.RecipeStepProductID),
			ArchivedAt:          database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumQuantity:     database.Uint32PointerFromNullInt32(result.MaximumQuantity),
			Notes:               result.Notes,
			Name:                result.Name,
			BelongsToRecipeStep: result.BelongsToRecipeStep,
			ID:                  result.ID,
			MinimumQuantity:     uint32(result.MinimumQuantity),
			OptionIndex:         uint16(result.OptionIndex),
			PreferenceRank:      uint8(result.PreferenceRank),
			Optional:            result.Optional,
		}

		if result.ValidInstrumentID.Valid {
			recipeStepInstrument.Instrument = &types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt.Time,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath.String,
				ID:                             result.ValidInstrumentID.String,
				Name:                           result.ValidInstrumentName.String,
				PluralName:                     result.ValidInstrumentPluralName.String,
				Description:                    result.ValidInstrumentDescription.String,
				Slug:                           result.ValidInstrumentSlug.String,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists.Bool,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions.Bool,
				UsableForStorage:               result.ValidInstrumentUsableForStorage.Bool,
			}
		}

		recipeStepInstruments = append(recipeStepInstruments, recipeStepInstrument)
	}

	return recipeStepInstruments, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *Querier) createRecipeStepInstrument(ctx context.Context, querier database.SQLQueryExecutor, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, input.ID)
	logger := q.logger.WithValue(keys.RecipeStepInstrumentIDKey, input.ID)

	// create the recipe step instrument.
	if err := q.generatedQuerier.CreateRecipeStepInstrument(ctx, querier, &generated.CreateRecipeStepInstrumentParams{
		ID:                  input.ID,
		Name:                input.Name,
		Notes:               input.Notes,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		InstrumentID:        database.NullStringFromStringPointer(input.InstrumentID),
		RecipeStepProductID: database.NullStringFromStringPointer(input.RecipeStepProductID),
		MaximumQuantity:     database.NullInt32FromUint32Pointer(input.MaximumQuantity),
		PreferenceRank:      int32(input.PreferenceRank),
		OptionIndex:         int32(input.OptionIndex),
		MinimumQuantity:     int32(input.MinimumQuantity),
		Optional:            input.Optional,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing recipe step instrument creation query")
	}

	x := &types.RecipeStepInstrument{
		ID:                  input.ID,
		Instrument:          nil,
		RecipeStepProductID: input.RecipeStepProductID,
		Name:                input.Name,
		Notes:               input.Notes,
		PreferenceRank:      input.PreferenceRank,
		BelongsToRecipeStep: input.BelongsToRecipeStep,
		Optional:            input.Optional,
		OptionIndex:         input.OptionIndex,
		MinimumQuantity:     input.MinimumQuantity,
		MaximumQuantity:     input.MaximumQuantity,
		CreatedAt:           q.currentTime(),
	}

	if input.InstrumentID != nil {
		x.Instrument = &types.ValidInstrument{ID: *input.InstrumentID}
	}

	logger.Info("recipe step instrument created")

	return x, nil
}

// CreateRecipeStepInstrument creates a recipe step instrument in the database.
func (q *Querier) CreateRecipeStepInstrument(ctx context.Context, input *types.RecipeStepInstrumentDatabaseCreationInput) (*types.RecipeStepInstrument, error) {
	return q.createRecipeStepInstrument(ctx, q.db, input)
}

// UpdateRecipeStepInstrument updates a particular recipe step instrument.
func (q *Querier) UpdateRecipeStepInstrument(ctx context.Context, updated *types.RecipeStepInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.RecipeStepInstrumentIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, updated.ID)

	var instrumentID *string
	if updated.Instrument != nil {
		instrumentID = &updated.Instrument.ID
	}

	if _, err := q.generatedQuerier.UpdateRecipeStepInstrument(ctx, q.db, &generated.UpdateRecipeStepInstrumentParams{
		InstrumentID:        database.NullStringFromStringPointer(instrumentID),
		RecipeStepProductID: database.NullStringFromStringPointer(updated.RecipeStepProductID),
		Name:                updated.Name,
		Notes:               updated.Notes,
		PreferenceRank:      int32(updated.PreferenceRank),
		Optional:            updated.Optional,
		OptionIndex:         int32(updated.OptionIndex),
		MinimumQuantity:     int32(updated.MinimumQuantity),
		MaximumQuantity:     database.NullInt32FromUint32Pointer(updated.MaximumQuantity),
		BelongsToRecipeStep: updated.BelongsToRecipeStep,
		ID:                  updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step instrument")
	}

	logger.Info("recipe step instrument updated")

	return nil
}

// ArchiveRecipeStepInstrument archives a recipe step instrument from the database by its ID.
func (q *Querier) ArchiveRecipeStepInstrument(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if recipeStepInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
	tracing.AttachToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)

	if _, err := q.generatedQuerier.ArchiveRecipeStepInstrument(ctx, q.db, &generated.ArchiveRecipeStepInstrumentParams{
		BelongsToRecipeStep: recipeStepID,
		ID:                  recipeStepInstrumentID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe step instrument")
	}

	logger.Info("recipe step instrument archived")

	return nil
}
