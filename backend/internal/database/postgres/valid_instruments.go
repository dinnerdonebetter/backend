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
	_ types.ValidInstrumentDataManager = (*Querier)(nil)
)

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (q *Querier) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	result, err := q.generatedQuerier.CheckValidInstrumentExistence(ctx, q.db, validInstrumentID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

// GetValidInstrument fetches a valid instrument from the database.
func (q *Querier) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	result, err := q.generatedQuerier.GetValidInstrument(ctx, q.db, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
	}

	validInstrument := &types.ValidInstrument{
		CreatedAt:                      result.CreatedAt,
		LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
		IconPath:                       result.IconPath,
		ID:                             result.ID,
		Name:                           result.Name,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Slug:                           result.Slug,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		UsableForStorage:               result.UsableForStorage,
	}

	return validInstrument, nil
}

// GetRandomValidInstrument fetches a valid instrument from the database.
func (q *Querier) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidInstrument(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validInstrument")
	}

	validInstrument := &types.ValidInstrument{
		CreatedAt:                      result.CreatedAt,
		LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
		IconPath:                       result.IconPath,
		ID:                             result.ID,
		Name:                           result.Name,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Slug:                           result.Slug,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		UsableForStorage:               result.UsableForStorage,
	}

	return validInstrument, nil
}

// SearchForValidInstruments fetches a valid instrument from the database.
func (q *Querier) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, query)

	results, err := q.generatedQuerier.SearchForValidInstruments(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredients list retrieval query")
	}

	validInstruments := []*types.ValidInstrument{}
	for _, result := range results {
		validInstrument := &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		}
		validInstruments = append(validInstruments, validInstrument)
	}

	return validInstruments, nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidInstrument], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidInstrument]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidInstruments(ctx, q.db, &generated.GetValidInstrumentsParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidInstrumentsWithIDs fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	results, err := q.generatedQuerier.GetValidInstrumentsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments id list retrieval query")
	}

	instruments := []*types.ValidInstrument{}
	for _, result := range results {
		instruments = append(instruments, &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		})
	}

	return instruments, nil
}

// GetValidInstrumentIDsThatNeedSearchIndexing fetches a list of valid instruments from the database that meet a particular filter.
func (q *Querier) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidInstrumentsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid instruments list retrieval query")
	}

	return results, nil
}

// CreateValidInstrument creates a valid instrument in the database.
func (q *Querier) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	// create the valid instrument.
	if err := q.generatedQuerier.CreateValidInstrument(ctx, q.db, &generated.CreateValidInstrumentParams{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Slug:                           input.Slug,
		UsableForStorage:               input.UsableForStorage,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid instrument creation query")
	}

	x := &types.ValidInstrument{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		UsableForStorage:               input.UsableForStorage,
		Slug:                           input.Slug,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		CreatedAt:                      q.currentTime(),
	}

	logger.Info("valid instrument created")

	return x, nil
}

// UpdateValidInstrument updates a particular valid instrument.
func (q *Querier) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidInstrument(ctx, q.db, &generated.UpdateValidInstrumentParams{
		Name:                           updated.Name,
		PluralName:                     updated.PluralName,
		Description:                    updated.Description,
		IconPath:                       updated.IconPath,
		Slug:                           updated.Slug,
		ID:                             updated.ID,
		UsableForStorage:               updated.UsableForStorage,
		DisplayInSummaryLists:          updated.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: updated.IncludeInGeneratedInstructions,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument updated")

	return nil
}

// MarkValidInstrumentAsIndexed updates a particular valid instrument's last_indexed_at value.
func (q *Querier) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	if _, err := q.generatedQuerier.UpdateValidInstrumentLastIndexedAt(ctx, q.db, validInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid instrument as indexed")
	}

	logger.Info("valid instrument marked as indexed")

	return nil
}

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (q *Querier) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validInstrumentID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	if _, err := q.generatedQuerier.ArchiveValidInstrument(ctx, q.db, validInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid instrument")
	}

	logger.Info("valid instrument archived")

	return nil
}
