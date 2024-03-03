package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.ValidVesselDataManager = (*Querier)(nil)
)

// ValidVesselExists fetches whether a valid vessel exists from the database.
func (q *Querier) ValidVesselExists(ctx context.Context, validVesselID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	result, err := q.generatedQuerier.CheckValidVesselExistence(ctx, q.db, validVesselID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid vessel existence check")
	}

	return result, nil
}

// GetValidVessel fetches a valid vessel from the database.
func (q *Querier) GetValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	result, err := q.generatedQuerier.GetValidVessel(ctx, q.db, validVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid vessel existence check")
	}

	validVessel := &types.ValidVessel{
		CreatedAt:     result.CreatedAt,
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		CapacityUnit: &types.ValidMeasurementUnit{
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
		IconPath:                       result.IconPath,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Name:                           result.Name,
		Slug:                           result.Slug,
		Shape:                          string(result.Shape),
		ID:                             result.ID,
		WidthInMillimeters:             database.Float32FromNullString(result.WidthInMillimeters),
		LengthInMillimeters:            database.Float32FromNullString(result.LengthInMillimeters),
		HeightInMillimeters:            database.Float32FromNullString(result.HeightInMillimeters),
		Capacity:                       database.Float32FromString(result.Capacity),
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		UsableForStorage:               result.UsableForStorage,
	}

	return validVessel, nil
}

// GetRandomValidVessel fetches a valid vessel from the database.
func (q *Querier) GetRandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidVessel(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "querying for random valid vessel")
	}

	validVessel := &types.ValidVessel{
		CreatedAt:     result.CreatedAt,
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		CapacityUnit: &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
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
		IconPath:                       result.IconPath,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Name:                           result.Name,
		Slug:                           result.Slug,
		Shape:                          string(result.Shape),
		ID:                             result.ID,
		WidthInMillimeters:             database.Float32FromNullString(result.WidthInMillimeters),
		LengthInMillimeters:            database.Float32FromNullString(result.LengthInMillimeters),
		HeightInMillimeters:            database.Float32FromNullString(result.HeightInMillimeters),
		Capacity:                       database.Float32FromString(result.Capacity),
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		UsableForStorage:               result.UsableForStorage,
	}

	return validVessel, nil
}

// SearchForValidVessels fetches a valid vessel from the database.
func (q *Querier) SearchForValidVessels(ctx context.Context, query string) ([]*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, query)

	results, err := q.generatedQuerier.SearchForValidVessels(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareError(err, span, "querying for valid vessel")
	}

	validVessels := []*types.ValidVessel{}
	for _, result := range results {
		validVessel := &types.ValidVessel{
			CreatedAt:                      result.CreatedAt,
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			IconPath:                       result.IconPath,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Name:                           result.Name,
			Slug:                           result.Slug,
			Shape:                          string(result.Shape),
			ID:                             result.ID,
			WidthInMillimeters:             database.Float32FromNullString(result.WidthInMillimeters),
			LengthInMillimeters:            database.Float32FromNullString(result.LengthInMillimeters),
			HeightInMillimeters:            database.Float32FromNullString(result.HeightInMillimeters),
			Capacity:                       database.Float32FromString(result.Capacity),
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			UsableForStorage:               result.UsableForStorage,
		}

		if result.CapacityUnit.Valid && result.CapacityUnit.String != "" {
			validVessel.CapacityUnit, err = q.GetValidMeasurementUnit(ctx, result.CapacityUnit.String)
			if err != nil {
				return nil, observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
			}
		}

		validVessels = append(validVessels, validVessel)
	}

	if len(validVessels) == 0 {
		return nil, sql.ErrNoRows
	}

	return validVessels, nil
}

// GetValidVessels fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVessels(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidVessel]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidVessels(ctx, q.db, &generated.GetValidVesselsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareError(err, span, "querying for valid vessels")
	}

	for _, result := range results {
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
		validVessel := &types.ValidVessel{
			CreatedAt:                      result.CreatedAt,
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			IconPath:                       result.IconPath,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Name:                           result.Name,
			Slug:                           result.Slug,
			Shape:                          string(result.Shape),
			ID:                             result.ID,
			WidthInMillimeters:             database.Float32FromNullString(result.WidthInMillimeters),
			LengthInMillimeters:            database.Float32FromNullString(result.LengthInMillimeters),
			HeightInMillimeters:            database.Float32FromNullString(result.HeightInMillimeters),
			Capacity:                       database.Float32FromString(result.Capacity),
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			UsableForStorage:               result.UsableForStorage,
		}

		if result.CapacityUnit.Valid && result.CapacityUnit.String != "" {
			validVessel.CapacityUnit, err = q.GetValidMeasurementUnit(ctx, result.CapacityUnit.String)
			if err != nil {
				return nil, observability.PrepareAndLogError(err, logger, span, "getting valid measurement unit")
			}
		}

		x.Data = append(x.Data, validVessel)
	}

	return x, nil
}

// GetValidVesselsWithIDs fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVesselsWithIDs(ctx context.Context, ids []string) ([]*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if len(ids) == 0 {
		return nil, sql.ErrNoRows
	}

	results, err := q.generatedQuerier.GetValidVesselsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid vessels id list retrieval query")
	}

	validVessels := []*types.ValidVessel{}
	for _, result := range results {
		validVessels = append(validVessels, &types.ValidVessel{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			CapacityUnit: &types.ValidMeasurementUnit{
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
			IconPath:                       result.IconPath,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Name:                           result.Name,
			Slug:                           result.Slug,
			Shape:                          string(result.Shape),
			ID:                             result.ID,
			WidthInMillimeters:             database.Float32FromNullString(result.WidthInMillimeters),
			LengthInMillimeters:            database.Float32FromNullString(result.LengthInMillimeters),
			HeightInMillimeters:            database.Float32FromNullString(result.HeightInMillimeters),
			Capacity:                       database.Float32FromString(result.Capacity),
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			UsableForStorage:               result.UsableForStorage,
		})
	}

	return validVessels, nil
}

// GetValidVesselIDsThatNeedSearchIndexing fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVesselIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidVesselIDsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid vessels list retrieval query")
	}

	return results, nil
}

// CreateValidVessel creates a valid vessel in the database.
func (q *Querier) CreateValidVessel(ctx context.Context, input *types.ValidVesselDatabaseCreationInput) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidVesselIDKey, input.ID)

	// create the valid vessel.
	if err := q.generatedQuerier.CreateValidVessel(ctx, q.db, &generated.CreateValidVesselParams{
		Slug:                           input.Slug,
		ID:                             input.ID,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Shape:                          generated.VesselShape(input.Shape),
		Name:                           input.Name,
		Capacity:                       database.StringFromFloat32(input.Capacity),
		CapacityUnit:                   database.NullStringFromStringPointer(input.CapacityUnitID),
		WidthInMillimeters:             database.NullStringFromFloat32(input.WidthInMillimeters),
		LengthInMillimeters:            database.NullStringFromFloat32(input.LengthInMillimeters),
		HeightInMillimeters:            database.NullStringFromFloat32(input.HeightInMillimeters),
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		UsableForStorage:               input.UsableForStorage,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid vessel creation query")
	}

	x := &types.ValidVessel{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Slug:                           input.Slug,
		Shape:                          input.Shape,
		WidthInMillimeters:             input.WidthInMillimeters,
		LengthInMillimeters:            input.LengthInMillimeters,
		HeightInMillimeters:            input.HeightInMillimeters,
		Capacity:                       input.Capacity,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		UsableForStorage:               input.UsableForStorage,
		CreatedAt:                      q.currentTime(),
	}

	if input.CapacityUnitID != nil {
		x.CapacityUnit = &types.ValidMeasurementUnit{ID: *input.CapacityUnitID}
	}

	tracing.AttachToSpan(span, keys.ValidVesselIDKey, x.ID)
	logger.Info("valid vessel created")

	return x, nil
}

// UpdateValidVessel updates a particular valid vessel.
func (q *Querier) UpdateValidVessel(ctx context.Context, updated *types.ValidVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidVesselIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, updated.ID)

	if updated.CapacityUnit == nil {
		return fmt.Errorf("capacity unit: %w", ErrNilInputProvided)
	}

	if _, err := q.generatedQuerier.UpdateValidVessel(ctx, q.db, &generated.UpdateValidVesselParams{
		Name:                           updated.Name,
		PluralName:                     updated.PluralName,
		Description:                    updated.Description,
		IconPath:                       updated.IconPath,
		UsableForStorage:               updated.UsableForStorage,
		Slug:                           updated.Slug,
		DisplayInSummaryLists:          updated.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: updated.IncludeInGeneratedInstructions,
		Capacity:                       database.StringFromFloat32(updated.Capacity),
		CapacityUnit:                   database.NullStringFromString(updated.CapacityUnit.ID),
		WidthInMillimeters:             database.NullStringFromFloat32(updated.WidthInMillimeters),
		LengthInMillimeters:            database.NullStringFromFloat32(updated.LengthInMillimeters),
		HeightInMillimeters:            database.NullStringFromFloat32(updated.HeightInMillimeters),
		Shape:                          generated.VesselShape(updated.Shape),
		ID:                             updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid vessel")
	}

	logger.Info("valid vessel updated")

	return nil
}

// MarkValidVesselAsIndexed updates a particular valid vessel's last_indexed_at value.
func (q *Querier) MarkValidVesselAsIndexed(ctx context.Context, validVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	if _, err := q.generatedQuerier.UpdateValidVesselLastIndexedAt(ctx, q.db, validVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid vessel as indexed")
	}

	logger.Info("valid vessel marked as indexed")

	return nil
}

// ArchiveValidVessel archives a valid vessel from the database by its ID.
func (q *Querier) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validVesselID)

	if _, err := q.generatedQuerier.ArchiveValidVessel(ctx, q.db, validVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid vessel")
	}

	logger.Info("valid vessel archived")

	return nil
}
