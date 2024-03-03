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
)

var (
	_ types.ValidPreparationDataManager = (*Querier)(nil)
)

// ValidPreparationExists fetches whether a valid preparation exists from the database.
func (q *Querier) ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	exists, err := q.generatedQuerier.CheckValidPreparationExistence(ctx, q.db, validPreparationID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "checking valid preparation existence")
	}

	return exists, nil
}

// GetValidPreparation fetches a valid preparation from the database.
func (q *Querier) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	result, err := q.generatedQuerier.GetValidPreparation(ctx, q.db, validPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting random valid preparation")
	}

	validPreparation := &types.ValidPreparation{
		CreatedAt:                   result.CreatedAt,
		MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.MaximumInstrumentCount),
		ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
		MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.MaximumIngredientCount),
		LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
		MaximumVesselCount:          database.Int32PointerFromNullInt32(result.MaximumVesselCount),
		IconPath:                    result.IconPath,
		PastTense:                   result.PastTense,
		ID:                          result.ID,
		Name:                        result.Name,
		Description:                 result.Description,
		Slug:                        result.Slug,
		MinimumIngredientCount:      result.MinimumIngredientCount,
		MinimumInstrumentCount:      result.MinimumInstrumentCount,
		MinimumVesselCount:          result.MinimumVesselCount,
		RestrictToIngredients:       result.RestrictToIngredients,
		TemperatureRequired:         result.TemperatureRequired,
		TimeEstimateRequired:        result.TimeEstimateRequired,
		ConditionExpressionRequired: result.ConditionExpressionRequired,
		ConsumesVessel:              result.ConsumesVessel,
		OnlyForVessels:              result.OnlyForVessels,
		YieldsNothing:               result.YieldsNothing,
	}

	return validPreparation, nil
}

// GetRandomValidPreparation fetches a valid preparation from the database.
func (q *Querier) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidPreparation(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid preparation")
	}

	validPreparation := &types.ValidPreparation{
		CreatedAt:                   result.CreatedAt,
		MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.MaximumInstrumentCount),
		ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
		MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.MaximumIngredientCount),
		LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
		MaximumVesselCount:          database.Int32PointerFromNullInt32(result.MaximumVesselCount),
		IconPath:                    result.IconPath,
		PastTense:                   result.PastTense,
		ID:                          result.ID,
		Name:                        result.Name,
		Description:                 result.Description,
		Slug:                        result.Slug,
		MinimumIngredientCount:      result.MinimumIngredientCount,
		MinimumInstrumentCount:      result.MinimumInstrumentCount,
		MinimumVesselCount:          result.MinimumVesselCount,
		RestrictToIngredients:       result.RestrictToIngredients,
		TemperatureRequired:         result.TemperatureRequired,
		TimeEstimateRequired:        result.TimeEstimateRequired,
		ConditionExpressionRequired: result.ConditionExpressionRequired,
		ConsumesVessel:              result.ConsumesVessel,
		OnlyForVessels:              result.OnlyForVessels,
		YieldsNothing:               result.YieldsNothing,
	}

	return validPreparation, nil
}

// SearchForValidPreparations fetches a valid preparation from the database.
func (q *Querier) SearchForValidPreparations(ctx context.Context, query string) ([]*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	results, err := q.generatedQuerier.SearchForValidPreparations(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparations search")
	}

	x := []*types.ValidPreparation{}
	for _, result := range results {
		x = append(x, &types.ValidPreparation{
			CreatedAt:                   result.CreatedAt,
			MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.MaximumInstrumentCount),
			ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.MaximumIngredientCount),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
			MaximumVesselCount:          database.Int32PointerFromNullInt32(result.MaximumVesselCount),
			IconPath:                    result.IconPath,
			PastTense:                   result.PastTense,
			ID:                          result.ID,
			Name:                        result.Name,
			Description:                 result.Description,
			Slug:                        result.Slug,
			MinimumIngredientCount:      result.MinimumIngredientCount,
			MinimumInstrumentCount:      result.MinimumInstrumentCount,
			MinimumVesselCount:          result.MinimumVesselCount,
			RestrictToIngredients:       result.RestrictToIngredients,
			TemperatureRequired:         result.TemperatureRequired,
			TimeEstimateRequired:        result.TimeEstimateRequired,
			ConditionExpressionRequired: result.ConditionExpressionRequired,
			ConsumesVessel:              result.ConsumesVessel,
			OnlyForVessels:              result.OnlyForVessels,
			YieldsNothing:               result.YieldsNothing,
		})
	}

	return x, nil
}

// GetValidPreparations fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparation], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidPreparations(ctx, q.db, &generated.GetValidPreparationsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparations list retrieval query")
	}

	for _, result := range results {
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
		x.Data = append(x.Data, &types.ValidPreparation{
			CreatedAt:                   result.CreatedAt,
			MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.MaximumInstrumentCount),
			ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.MaximumIngredientCount),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
			MaximumVesselCount:          database.Int32PointerFromNullInt32(result.MaximumVesselCount),
			IconPath:                    result.IconPath,
			PastTense:                   result.PastTense,
			ID:                          result.ID,
			Name:                        result.Name,
			Description:                 result.Description,
			Slug:                        result.Slug,
			MinimumIngredientCount:      result.MinimumIngredientCount,
			MinimumInstrumentCount:      result.MinimumInstrumentCount,
			MinimumVesselCount:          result.MinimumVesselCount,
			RestrictToIngredients:       result.RestrictToIngredients,
			TemperatureRequired:         result.TemperatureRequired,
			TimeEstimateRequired:        result.TimeEstimateRequired,
			ConditionExpressionRequired: result.ConditionExpressionRequired,
			ConsumesVessel:              result.ConsumesVessel,
			OnlyForVessels:              result.OnlyForVessels,
			YieldsNothing:               result.YieldsNothing,
		})
	}

	return x, nil
}

// GetValidPreparationsWithIDs fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if len(ids) == 0 {
		return nil, sql.ErrNoRows
	}
	logger := q.logger.WithValue("ids_count", len(ids))

	results, err := q.generatedQuerier.GetValidPreparationsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid preparations by ID")
	}

	preparations := []*types.ValidPreparation{}
	for _, result := range results {
		preparations = append(preparations, &types.ValidPreparation{
			CreatedAt:                   result.CreatedAt,
			MaximumInstrumentCount:      database.Int32PointerFromNullInt32(result.MaximumInstrumentCount),
			ArchivedAt:                  database.TimePointerFromNullTime(result.ArchivedAt),
			MaximumIngredientCount:      database.Int32PointerFromNullInt32(result.MaximumIngredientCount),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.LastUpdatedAt),
			MaximumVesselCount:          database.Int32PointerFromNullInt32(result.MaximumVesselCount),
			IconPath:                    result.IconPath,
			PastTense:                   result.PastTense,
			ID:                          result.ID,
			Name:                        result.Name,
			Description:                 result.Description,
			Slug:                        result.Slug,
			MinimumIngredientCount:      result.MinimumIngredientCount,
			MinimumInstrumentCount:      result.MinimumInstrumentCount,
			MinimumVesselCount:          result.MinimumVesselCount,
			RestrictToIngredients:       result.RestrictToIngredients,
			TemperatureRequired:         result.TemperatureRequired,
			TimeEstimateRequired:        result.TimeEstimateRequired,
			ConditionExpressionRequired: result.ConditionExpressionRequired,
			ConsumesVessel:              result.ConsumesVessel,
			OnlyForVessels:              result.OnlyForVessels,
			YieldsNothing:               result.YieldsNothing,
		})
	}

	return preparations, nil
}

// GetValidPreparationIDsThatNeedSearchIndexing fetches a list of valid preparations from the database that meet a particular filter.
func (q *Querier) GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidPreparationsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid preparations list retrieval query")
	}

	return results, nil
}

// CreateValidPreparation creates a valid preparation in the database.
func (q *Querier) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationDatabaseCreationInput) (*types.ValidPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationIDKey, input.ID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, input.ID)

	// create the valid preparation.
	if err := q.generatedQuerier.CreateValidPreparation(ctx, q.db, &generated.CreateValidPreparationParams{
		ID:                          input.ID,
		Name:                        input.Name,
		Description:                 input.Description,
		IconPath:                    input.IconPath,
		YieldsNothing:               input.YieldsNothing,
		RestrictToIngredients:       input.RestrictToIngredients,
		MinimumIngredientCount:      input.MinimumIngredientCount,
		MaximumIngredientCount:      database.NullInt32FromInt32Pointer(input.MaximumIngredientCount),
		MinimumInstrumentCount:      input.MinimumInstrumentCount,
		MaximumInstrumentCount:      database.NullInt32FromInt32Pointer(input.MaximumInstrumentCount),
		TemperatureRequired:         input.TemperatureRequired,
		TimeEstimateRequired:        input.TimeEstimateRequired,
		ConditionExpressionRequired: input.ConditionExpressionRequired,
		ConsumesVessel:              input.ConsumesVessel,
		OnlyForVessels:              input.OnlyForVessels,
		MinimumVesselCount:          input.MinimumVesselCount,
		MaximumVesselCount:          database.NullInt32FromInt32Pointer(input.MaximumVesselCount),
		PastTense:                   input.PastTense,
		Slug:                        input.Slug,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation creation query")
	}

	x := &types.ValidPreparation{
		ID:                          input.ID,
		Name:                        input.Name,
		Description:                 input.Description,
		IconPath:                    input.IconPath,
		YieldsNothing:               input.YieldsNothing,
		RestrictToIngredients:       input.RestrictToIngredients,
		Slug:                        input.Slug,
		PastTense:                   input.PastTense,
		MinimumIngredientCount:      input.MinimumIngredientCount,
		MaximumIngredientCount:      input.MaximumIngredientCount,
		MinimumInstrumentCount:      input.MinimumInstrumentCount,
		MaximumInstrumentCount:      input.MaximumInstrumentCount,
		TemperatureRequired:         input.TemperatureRequired,
		TimeEstimateRequired:        input.TimeEstimateRequired,
		ConditionExpressionRequired: input.ConditionExpressionRequired,
		ConsumesVessel:              input.ConsumesVessel,
		OnlyForVessels:              input.OnlyForVessels,
		MinimumVesselCount:          input.MinimumVesselCount,
		MaximumVesselCount:          input.MaximumVesselCount,
		CreatedAt:                   q.currentTime(),
	}

	logger.Info("valid preparation created")

	return x, nil
}

// UpdateValidPreparation updates a particular valid preparation.
func (q *Querier) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidPreparation(ctx, q.db, &generated.UpdateValidPreparationParams{
		Description:                 updated.Description,
		IconPath:                    updated.IconPath,
		ID:                          updated.ID,
		Name:                        updated.Name,
		PastTense:                   updated.PastTense,
		Slug:                        updated.Slug,
		MaximumIngredientCount:      database.NullInt32FromInt32Pointer(updated.MaximumIngredientCount),
		MaximumInstrumentCount:      database.NullInt32FromInt32Pointer(updated.MaximumInstrumentCount),
		MaximumVesselCount:          database.NullInt32FromInt32Pointer(updated.MaximumVesselCount),
		MinimumVesselCount:          updated.MinimumVesselCount,
		MinimumIngredientCount:      updated.MinimumIngredientCount,
		MinimumInstrumentCount:      updated.MinimumInstrumentCount,
		RestrictToIngredients:       updated.RestrictToIngredients,
		OnlyForVessels:              updated.OnlyForVessels,
		ConsumesVessel:              updated.ConsumesVessel,
		ConditionExpressionRequired: updated.ConditionExpressionRequired,
		TimeEstimateRequired:        updated.TimeEstimateRequired,
		TemperatureRequired:         updated.TemperatureRequired,
		YieldsNothing:               updated.YieldsNothing,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation updated")

	return nil
}

// MarkValidPreparationAsIndexed updates a particular valid preparation's last_indexed_at value.
func (q *Querier) MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	if _, err := q.generatedQuerier.UpdateValidPreparationLastIndexedAt(ctx, q.db, validPreparationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid preparation as indexed")
	}

	logger.Info("valid preparation marked as indexed")

	return nil
}

// ArchiveValidPreparation archives a valid preparation from the database by its ID.
func (q *Querier) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	if _, err := q.generatedQuerier.ArchiveValidPreparation(ctx, q.db, validPreparationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation")
	}

	logger.Info("valid preparation archived")

	return nil
}
