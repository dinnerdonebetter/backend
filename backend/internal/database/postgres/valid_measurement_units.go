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
	_ types.ValidMeasurementUnitDataManager = (*Querier)(nil)
)

// ValidMeasurementUnitExists fetches whether a valid measurement unit exists from the database.
func (q *Querier) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	result, err := q.generatedQuerier.CheckValidMeasurementUnitExistence(ctx, q.db, validMeasurementUnitID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit existence check")
	}

	return result, nil
}

// GetValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *Querier) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	result, err := q.generatedQuerier.GetValidMeasurementUnit(ctx, q.db, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement unit")
	}

	validMeasurementUnit := &types.ValidMeasurementUnit{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		Name:          result.Name,
		IconPath:      result.IconPath,
		ID:            result.ID,
		Description:   result.Description,
		PluralName:    result.PluralName,
		Slug:          result.Slug,
		Volumetric:    database.BoolFromNullBool(result.Volumetric),
		Universal:     result.Universal,
		Metric:        result.Metric,
		Imperial:      result.Imperial,
	}

	return validMeasurementUnit, nil
}

// GetRandomValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *Querier) GetRandomValidMeasurementUnit(ctx context.Context) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidMeasurementUnit(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid measurement unit")
	}

	validMeasurementUnit := &types.ValidMeasurementUnit{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		Name:          result.Name,
		IconPath:      result.IconPath,
		ID:            result.ID,
		Description:   result.Description,
		PluralName:    result.PluralName,
		Slug:          result.Slug,
		Volumetric:    database.BoolFromNullBool(result.Volumetric),
		Universal:     result.Universal,
		Metric:        result.Metric,
		Imperial:      result.Imperial,
	}

	return validMeasurementUnit, nil
}

// SearchForValidMeasurementUnits fetches a valid measurement unit from the database.
func (q *Querier) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, query)

	results, err := q.generatedQuerier.SearchForValidMeasurementUnits(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	x := []*types.ValidMeasurementUnit{}
	for _, result := range results {
		x = append(x, &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		})
	}

	return x, nil
}

// ValidMeasurementUnitsForIngredientID fetches a valid measurement unit from the database.
func (q *Querier) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &types.QueryFilteredResult[types.ValidMeasurementUnit]{
		Pagination: filter.ToPagination(),
	}

	if validIngredientID == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	results, err := q.generatedQuerier.SearchValidMeasurementUnitsByIngredientID(ctx, q.db, &generated.SearchValidMeasurementUnitsByIngredientIDParams{
		CreatedBefore:     database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:      database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:     database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:      database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:       database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:        database.NullInt32FromUint8Pointer(filter.Limit),
		ValidIngredientID: validIngredientID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	for _, result := range results {
		validMeasurementUnit := &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		}

		x.Data = append(x.Data, validMeasurementUnit)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidMeasurementUnits fetches a list of valid measurement units from the database that meet a particular filter.
func (q *Querier) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidMeasurementUnit], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidMeasurementUnit]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidMeasurementUnits(ctx, q.db, &generated.GetValidMeasurementUnitsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	for _, result := range results {
		validMeasurementUnit := &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		}

		x.Data = append(x.Data, validMeasurementUnit)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidMeasurementUnitsWithIDs fetches a list of valid measurement unit from the database that meet a particular filter.
func (q *Querier) GetValidMeasurementUnitsWithIDs(ctx context.Context, ids []string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	results, err := q.generatedQuerier.GetValidMeasurementUnitsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement unit id list retrieval query")
	}

	x := []*types.ValidMeasurementUnit{}
	for _, result := range results {
		x = append(x, &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		})
	}

	return x, nil
}

// GetValidMeasurementUnitIDsThatNeedSearchIndexing fetches a list of valid measurement units from the database that meet a particular filter.
func (q *Querier) GetValidMeasurementUnitIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidMeasurementUnitsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid measurement units list retrieval query")
	}

	return results, nil
}

// CreateValidMeasurementUnit creates a valid measurement unit in the database.
func (q *Querier) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, input.ID)

	// create the valid measurement unit.
	if err := q.generatedQuerier.CreateValidMeasurementUnit(ctx, q.db, &generated.CreateValidMeasurementUnitParams{
		Name:        input.Name,
		Description: input.Description,
		IconPath:    input.IconPath,
		Slug:        input.Slug,
		PluralName:  input.PluralName,
		ID:          input.ID,
		Volumetric:  database.NullBoolFromBool(input.Volumetric),
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit creation query")
	}

	x := &types.ValidMeasurementUnit{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
		Slug:        input.Slug,
		PluralName:  input.PluralName,
		CreatedAt:   q.currentTime(),
	}

	logger.Info("valid measurement unit created")

	return x, nil
}

// UpdateValidMeasurementUnit updates a particular valid measurement unit.
func (q *Querier) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidMeasurementUnit(ctx, q.db, &generated.UpdateValidMeasurementUnitParams{
		Name:        updated.Name,
		Description: updated.Description,
		IconPath:    updated.IconPath,
		Slug:        updated.Slug,
		PluralName:  updated.PluralName,
		ID:          updated.ID,
		Volumetric:  database.NullBoolFromBool(updated.Volumetric),
		Universal:   updated.Universal,
		Metric:      updated.Metric,
		Imperial:    updated.Imperial,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit updated")

	return nil
}

// MarkValidMeasurementUnitAsIndexed updates a particular valid measurement unit's last_indexed_at value.
func (q *Querier) MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	if _, err := q.generatedQuerier.UpdateValidMeasurementUnitLastIndexedAt(ctx, q.db, validMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid measurement unit as indexed")
	}

	logger.Info("valid measurement unit marked as indexed")

	return nil
}

// ArchiveValidMeasurementUnit archives a valid measurement unit from the database by its ID.
func (q *Querier) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	if _, err := q.generatedQuerier.ArchiveValidMeasurementUnit(ctx, q.db, validMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement unit")
	}

	logger.Info("valid measurement unit archived")

	return nil
}
