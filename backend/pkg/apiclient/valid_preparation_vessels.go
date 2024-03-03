package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidPreparationVessel gets a valid ingredient vessel.
func (c *Client) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	req, err := c.requestBuilder.BuildGetValidPreparationVesselRequest(ctx, validPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid ingredient vessel request")
	}

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid ingredient vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetValidPreparationVessels retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVessels(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidPreparationVesselsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation vessels list request")
	}

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// GetValidPreparationVesselsForPreparation retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVesselsForPreparation(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, validPreparationID)
	tracing.AttachToSpan(span, keys.ValidPreparationIDKey, validPreparationID)

	req, err := c.requestBuilder.BuildGetValidPreparationVesselsForPreparationRequest(ctx, validPreparationID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation vessels list request")
	}

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// GetValidPreparationVesselsForVessel retrieves a list of valid preparation vessels.
func (c *Client) GetValidPreparationVesselsForVessel(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	req, err := c.requestBuilder.BuildGetValidPreparationVesselsForVesselRequest(ctx, validInstrumentID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid preparation vessels list request")
	}

	var apiResponse *types.APIResponse[[]*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid preparation vessels")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	validPreparationVessels := &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return validPreparationVessels, nil
}

// CreateValidPreparationVessel creates a valid ingredient vessel.
func (c *Client) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidPreparationVesselRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid ingredient vessel request")
	}

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid ingredient vessel")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateValidPreparationVessel updates a valid ingredient vessel.
func (c *Client) UpdateValidPreparationVessel(ctx context.Context, validPreparationVessel *types.ValidPreparationVessel) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVessel == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVessel.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVessel.ID)

	req, err := c.requestBuilder.BuildUpdateValidPreparationVesselRequest(ctx, validPreparationVessel)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid ingredient vessel request")
	}

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient vessel %s", validPreparationVessel.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveValidPreparationVessel archives a valid ingredient vessel.
func (c *Client) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validPreparationVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	req, err := c.requestBuilder.BuildArchiveValidPreparationVesselRequest(ctx, validPreparationVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid ingredient vessel request")
	}

	var apiResponse *types.APIResponse[*types.ValidPreparationVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient vessel %s", validPreparationVesselID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
