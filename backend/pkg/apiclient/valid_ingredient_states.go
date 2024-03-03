package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidIngredientState gets a valid ingredient state.
func (c *Client) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	req, err := c.requestBuilder.BuildGetValidIngredientStateRequest(ctx, validIngredientStateID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient state request")
	}

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient state")
	}

	return apiResponse.Data, nil
}

// SearchValidIngredientStates searches through a list of valid ingredient states.
func (c *Client) SearchValidIngredientStates(ctx context.Context, query string, limit uint8) ([]*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if query == "" {
		return nil, ErrEmptyQueryProvided
	}

	if limit == 0 {
		limit = types.DefaultLimit
	}

	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue(keys.FilterLimitKey, limit)

	req, err := c.requestBuilder.BuildSearchValidIngredientStatesRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for valid ingredient states request")
	}

	var apiResponse *types.APIResponse[[]*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient states")
	}

	return apiResponse.Data, nil
}

// GetValidIngredientStates retrieves a list of valid ingredient states.
func (c *Client) GetValidIngredientStates(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientState], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidIngredientStatesRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid ingredient states list request")
	}

	var apiResponse *types.APIResponse[[]*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient states")
	}

	response := &types.QueryFilteredResult[types.ValidIngredientState]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateValidIngredientState creates a valid ingredient state.
func (c *Client) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidIngredientStateRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient state request")
	}

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient state")
	}

	return apiResponse.Data, nil
}

// UpdateValidIngredientState updates a valid ingredient state.
func (c *Client) UpdateValidIngredientState(ctx context.Context, validIngredientState *types.ValidIngredientState) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientState == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientState.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientState.ID)

	req, err := c.requestBuilder.BuildUpdateValidIngredientStateRequest(ctx, validIngredientState)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient state request")
	}

	var apiResponse *types.APIResponse[*types.ValidIngredientState]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient state %s", validIngredientState.ID)
	}

	return nil
}

// ArchiveValidIngredientState archives a valid ingredient state.
func (c *Client) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientStateID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientStateIDKey, validIngredientStateID)
	tracing.AttachToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)

	req, err := c.requestBuilder.BuildArchiveValidIngredientStateRequest(ctx, validIngredientStateID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient state request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient state %s", validIngredientStateID)
	}

	return nil
}
