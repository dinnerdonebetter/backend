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
	_ types.HouseholdInstrumentOwnershipDataManager = (*Querier)(nil)
)

// HouseholdInstrumentOwnershipExists fetches whether a household instrument ownership exists from the database.
func (q *Querier) HouseholdInstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	if householdID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	result, err := q.generatedQuerier.CheckHouseholdInstrumentOwnershipExistence(ctx, q.db, &generated.CheckHouseholdInstrumentOwnershipExistenceParams{
		ID:                 householdInstrumentOwnershipID,
		BelongsToHousehold: householdID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing household instrument ownership existence check")
	}

	return result, nil
}

// GetHouseholdInstrumentOwnership fetches a household instrument ownership from the database.
func (q *Querier) GetHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	result, err := q.generatedQuerier.GetHouseholdInstrumentOwnership(ctx, q.db, &generated.GetHouseholdInstrumentOwnershipParams{
		ID:                 householdInstrumentOwnershipID,
		BelongsToHousehold: householdID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household instrument ownership")
	}

	householdInstrumentOwnership := &types.HouseholdInstrumentOwnership{
		CreatedAt:          result.CreatedAt,
		ArchivedAt:         database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:      database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:                 result.ID,
		Notes:              result.Notes,
		BelongsToHousehold: result.BelongsToHousehold,
		Quantity:           uint16(result.Quantity),
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
	}

	return householdInstrumentOwnership, nil
}

// GetHouseholdInstrumentOwnerships fetches a list of household instrument ownerships from the database that meet a particular filter.
func (q *Querier) GetHouseholdInstrumentOwnerships(ctx context.Context, householdID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.HouseholdInstrumentOwnership], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{
		Pagination: filter.ToPagination(),
	}

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	results, err := q.generatedQuerier.GetHouseholdInstrumentOwnerships(ctx, q.db, &generated.GetHouseholdInstrumentOwnershipsParams{
		HouseholdID:   householdID,
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching household instrument ownerships")
	}

	for _, result := range results {
		householdInstrumentOwnership := &types.HouseholdInstrumentOwnership{
			CreatedAt:          result.CreatedAt,
			ArchivedAt:         database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:      database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:                 result.ID,
			Notes:              result.Notes,
			BelongsToHousehold: result.BelongsToHousehold,
			Quantity:           uint16(result.Quantity),
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
		}

		x.Data = append(x.Data, householdInstrumentOwnership)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateHouseholdInstrumentOwnership creates a household instrument ownership in the database.
func (q *Querier) CreateHouseholdInstrumentOwnership(ctx context.Context, input *types.HouseholdInstrumentOwnershipDatabaseCreationInput) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, input.ID)
	logger := q.logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, input.ID)

	// create the household instrument ownership.
	if err := q.generatedQuerier.CreateHouseholdInstrumentOwnership(ctx, q.db, &generated.CreateHouseholdInstrumentOwnershipParams{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidInstrumentID:  input.ValidInstrumentID,
		BelongsToHousehold: input.BelongsToHousehold,
		Quantity:           int32(input.Quantity),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household instrument ownership creation query")
	}

	x := &types.HouseholdInstrumentOwnership{
		ID:                 input.ID,
		Notes:              input.Notes,
		Quantity:           input.Quantity,
		Instrument:         types.ValidInstrument{ID: input.ValidInstrumentID},
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	logger.Info("household instrument ownership created")

	return x, nil
}

// UpdateHouseholdInstrumentOwnership updates a particular household instrument ownership.
func (q *Querier) UpdateHouseholdInstrumentOwnership(ctx context.Context, updated *types.HouseholdInstrumentOwnership) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateHouseholdInstrumentOwnership(ctx, q.db, &generated.UpdateHouseholdInstrumentOwnershipParams{
		Notes:              updated.Notes,
		ValidInstrumentID:  updated.Instrument.ID,
		ID:                 updated.ID,
		BelongsToHousehold: updated.BelongsToHousehold,
		Quantity:           int32(updated.Quantity),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating household instrument ownership")
	}

	logger.Info("household instrument ownership updated")

	return nil
}

// ArchiveHouseholdInstrumentOwnership archives a household instrument ownership from the database by its ID.
func (q *Querier) ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	if householdID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if _, err := q.generatedQuerier.ArchiveHouseholdInstrumentOwnership(ctx, q.db, &generated.ArchiveHouseholdInstrumentOwnershipParams{
		ID:                 householdInstrumentOwnershipID,
		BelongsToHousehold: householdID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving household instrument ownership")
	}

	logger.Info("household instrument ownership archived")

	return nil
}
