package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetValidVessel gets a valid vessel.
func (c *Client) GetValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	req, err := c.requestBuilder.BuildGetValidVesselRequest(ctx, validVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid vessel request")
	}

	var validVessel *types.ValidVessel
	if err = c.fetchAndUnmarshal(ctx, req, &validVessel); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}

	return validVessel, nil
}

// GetRandomValidVessel gets a valid vessel.
func (c *Client) GetRandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	req, err := c.requestBuilder.BuildGetRandomValidVesselRequest(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get valid vessel request")
	}

	var validVessel *types.ValidVessel
	if err = c.fetchAndUnmarshal(ctx, req, &validVessel); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessel")
	}

	return validVessel, nil
}

// SearchValidVessels searches through a list of valid vessels.
func (c *Client) SearchValidVessels(ctx context.Context, query string, limit uint8) ([]*types.ValidVessel, error) {
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

	req, err := c.requestBuilder.BuildSearchValidVesselsRequest(ctx, query, limit)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building search for valid vessels request")
	}

	var validVessels []*types.ValidVessel
	if err = c.fetchAndUnmarshal(ctx, req, &validVessels); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessels")
	}

	return validVessels, nil
}

// GetValidVessels retrieves a list of valid vessels.
func (c *Client) GetValidVessels(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidVessel], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetValidVesselsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building valid vessels list request")
	}

	var validVessels *types.QueryFilteredResult[types.ValidVessel]
	if err = c.fetchAndUnmarshal(ctx, req, &validVessels); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving valid vessels")
	}

	return validVessels, nil
}

// CreateValidVessel creates a valid vessel.
func (c *Client) CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateValidVesselRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create valid vessel request")
	}

	var validVessel *types.ValidVessel
	if err = c.fetchAndUnmarshal(ctx, req, &validVessel); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating valid vessel")
	}

	return validVessel, nil
}

// UpdateValidVessel updates a valid vessel.
func (c *Client) UpdateValidVessel(ctx context.Context, validVessel *types.ValidVessel) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVessel == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVessel.ID)
	tracing.AttachValidVesselIDToSpan(span, validVessel.ID)

	req, err := c.requestBuilder.BuildUpdateValidVesselRequest(ctx, validVessel)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update valid vessel request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &validVessel); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid vessel %s", validVessel.ID)
	}

	return nil
}

// ArchiveValidVessel archives a valid vessel.
func (c *Client) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	req, err := c.requestBuilder.BuildArchiveValidVesselRequest(ctx, validVesselID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive valid vessel request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid vessel %s", validVesselID)
	}

	return nil
}
