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
	_ types.ValidIngredientStateDataManager = (*Querier)(nil)
)

// ValidIngredientStateExists fetches whether a valid ingredient state exists from the database.
func (q *Querier) ValidIngredientStateExists(ctx context.Context, validIngredientStateID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	result, err := q.generatedQuerier.CheckValidIngredientStateExistence(ctx, q.db, validIngredientStateID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state existence check")
	}

	return result, nil
}

// GetValidIngredientState fetches a valid ingredient state from the database.
func (q *Querier) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	result, err := q.generatedQuerier.GetValidIngredientState(ctx, q.db, validIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient state retrieval query")
	}

	validIngredientState := &types.ValidIngredientState{
		CreatedAt:     result.CreatedAt,
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		PastTense:     result.PastTense,
		Description:   result.Description,
		IconPath:      result.IconPath,
		ID:            result.ID,
		Name:          result.Name,
		AttributeType: string(result.AttributeType),
		Slug:          result.Slug,
	}

	return validIngredientState, nil
}

// SearchForValidIngredientStates fetches a valid ingredient state from the database.
func (q *Querier) SearchForValidIngredientStates(ctx context.Context, query string) ([]*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, query)

	results, err := q.generatedQuerier.SearchForValidIngredientStates(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient states list retrieval query")
	}

	x := []*types.ValidIngredientState{}
	for _, result := range results {
		x = append(x, &types.ValidIngredientState{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			PastTense:     result.PastTense,
			Description:   result.Description,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Name:          result.Name,
			AttributeType: string(result.AttributeType),
			Slug:          result.Slug,
		})
	}

	return x, nil
}

// GetValidIngredientStates fetches a list of valid ingredient states from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStates(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidIngredientState], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidIngredientState]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientStates(ctx, q.db, &generated.GetValidIngredientStatesParams{
		CreatedBefore: database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    database.NullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient states list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientState{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			PastTense:     result.PastTense,
			Description:   result.Description,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Name:          result.Name,
			AttributeType: string(result.AttributeType),
			Slug:          result.Slug,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientStatesWithIDs fetches a list of valid ingredientStates from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStatesWithIDs(ctx context.Context, ids []string) ([]*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	results, err := q.generatedQuerier.GetValidIngredientStatesWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient states id list retrieval query")
	}

	ingredientStates := []*types.ValidIngredientState{}
	for _, result := range results {
		ingredientStates = append(ingredientStates, &types.ValidIngredientState{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			PastTense:     result.PastTense,
			Description:   result.Description,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Name:          result.Name,
			AttributeType: string(result.AttributeType),
			Slug:          result.Slug,
		})
	}

	return ingredientStates, nil
}

// GetValidIngredientStateIDsThatNeedSearchIndexing fetches a list of valid ingredientStates from the database that meet a particular filter.
func (q *Querier) GetValidIngredientStateIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidIngredientStatesNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid ingredient states list retrieval query")
	}

	return results, nil
}

// CreateValidIngredientState creates a valid ingredient state in the database.
func (q *Querier) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateDatabaseCreationInput) (*types.ValidIngredientState, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientStateIDKey, input.ID)

	// create the valid ingredient state.
	if err := q.generatedQuerier.CreateValidIngredientState(ctx, q.db, &generated.CreateValidIngredientStateParams{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		IconPath:      input.IconPath,
		PastTense:     input.PastTense,
		Slug:          input.Slug,
		AttributeType: generated.IngredientAttributeType(input.AttributeType),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient state creation query")
	}

	x := &types.ValidIngredientState{
		ID:            input.ID,
		Name:          input.Name,
		Description:   input.Description,
		IconPath:      input.IconPath,
		Slug:          input.Slug,
		PastTense:     input.PastTense,
		AttributeType: input.AttributeType,
		CreatedAt:     q.currentTime(),
	}

	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, x.ID)
	logger.Info("valid ingredient state created")

	return x, nil
}

// UpdateValidIngredientState updates a particular valid ingredient state.
func (q *Querier) UpdateValidIngredientState(ctx context.Context, updated *types.ValidIngredientState) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientStateIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredientState(ctx, q.db, &generated.UpdateValidIngredientStateParams{
		Name:          updated.Name,
		Description:   updated.Description,
		IconPath:      updated.IconPath,
		Slug:          updated.Slug,
		PastTense:     updated.PastTense,
		AttributeType: generated.IngredientAttributeType(updated.AttributeType),
		ID:            updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state")
	}

	logger.Info("valid ingredient state updated")

	return nil
}

// MarkValidIngredientStateAsIndexed updates a particular valid ingredient state's last_indexed_at value.
func (q *Querier) MarkValidIngredientStateAsIndexed(ctx context.Context, validIngredientStateID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	if _, err := q.generatedQuerier.UpdateValidIngredientStateLastIndexedAt(ctx, q.db, validIngredientStateID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid ingredient state as indexed")
	}

	logger.Info("valid ingredient state marked as indexed")

	return nil
}

// ArchiveValidIngredientState archives a valid ingredient state from the database by its ID.
func (q *Querier) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	if _, err := q.generatedQuerier.ArchiveValidIngredientState(ctx, q.db, validIngredientStateID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state")
	}

	logger.Info("valid ingredient state archived")

	return nil
}
