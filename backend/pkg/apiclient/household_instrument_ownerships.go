package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetHouseholdInstrumentOwnership gets a household instrument ownership.
func (c *Client) GetHouseholdInstrumentOwnership(ctx context.Context, validInstrumentID string) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, validInstrumentID)

	req, err := c.requestBuilder.BuildGetHouseholdInstrumentOwnershipRequest(ctx, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get household instrument ownership request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving household instrument ownership")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetHouseholdInstrumentOwnerships retrieves a list of household instrument ownerships.
func (c *Client) GetHouseholdInstrumentOwnerships(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInstrumentOwnership], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetHouseholdInstrumentOwnershipsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building household instrument ownerships list request")
	}

	var apiResponse *types.APIResponse[[]*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving household instrument ownerships")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.HouseholdInstrumentOwnership]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateHouseholdInstrumentOwnership creates a household instrument ownership.
func (c *Client) CreateHouseholdInstrumentOwnership(ctx context.Context, input *types.HouseholdInstrumentOwnershipCreationRequestInput) (*types.HouseholdInstrumentOwnership, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateHouseholdInstrumentOwnershipRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create household instrument ownership request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating household instrument ownership")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateHouseholdInstrumentOwnership updates a household instrument ownership.
func (c *Client) UpdateHouseholdInstrumentOwnership(ctx context.Context, instrumentOwnership *types.HouseholdInstrumentOwnership) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if instrumentOwnership == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, instrumentOwnership.ID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, instrumentOwnership.ID)

	req, err := c.requestBuilder.BuildUpdateHouseholdInstrumentOwnershipRequest(ctx, instrumentOwnership)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update household instrument ownership request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating household instrument ownership %s", instrumentOwnership.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveHouseholdInstrumentOwnership archives a household instrument ownership.
func (c *Client) ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInstrumentOwnershipID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	req, err := c.requestBuilder.BuildArchiveHouseholdInstrumentOwnershipRequest(ctx, householdInstrumentOwnershipID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive household instrument ownership request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInstrumentOwnership]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving household instrument ownership %s", householdInstrumentOwnershipID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
