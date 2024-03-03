package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetHouseholdInvitation retrieves a household invitation.
func (c *Client) GetHouseholdInvitation(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	if householdInvitationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	req, err := c.requestBuilder.BuildGetHouseholdInvitationRequest(ctx, householdID, householdInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get invitation request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetPendingHouseholdInvitationsFromUser retrieves household invitations sent by the user.
func (c *Client) GetPendingHouseholdInvitationsFromUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	filter.AttachToLogger(logger)

	req, err := c.requestBuilder.BuildGetPendingHouseholdInvitationsFromUserRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var apiResponse *types.APIResponse[[]*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// GetPendingHouseholdInvitationsForUser retrieves household invitations received by the user.
func (c *Client) GetPendingHouseholdInvitationsForUser(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	filter.AttachToLogger(logger)

	req, err := c.requestBuilder.BuildGetPendingHouseholdInvitationsForUserRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var apiResponse *types.APIResponse[[]*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.HouseholdInvitation]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// AcceptHouseholdInvitation accepts a given household invitation.
func (c *Client) AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildAcceptHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// CancelHouseholdInvitation cancels a given household invitation.
func (c *Client) CancelHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildCancelHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RejectHouseholdInvitation rejects a given household invitation.
func (c *Client) RejectHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if householdInvitationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdInvitationIDKey, householdInvitationID)

	if token == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildRejectHouseholdInvitationRequest(ctx, householdInvitationID, token, note)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building reject invitation request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "rejecting invitation")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
