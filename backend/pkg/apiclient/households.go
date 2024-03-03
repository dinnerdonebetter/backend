package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// SwitchActiveHousehold will switch the household on whose behalf requests are made.
func (c *Client) SwitchActiveHousehold(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if c.authMethod == cookieAuthMethod {
		req, err := c.requestBuilder.BuildSwitchActiveHouseholdRequest(ctx, householdID)
		if err != nil {
			return observability.PrepareError(err, span, "building household switch request")
		}

		var apiResponse *types.APIResponse[*types.Household]
		if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
			return observability.PrepareError(err, span, "executing household switch request")
		}

		if err = apiResponse.Error.AsError(); err != nil {
			return err
		}
	}

	c.householdID = householdID

	return nil
}

// GetCurrentHousehold retrieves a household.
func (c *Client) GetCurrentHousehold(ctx context.Context) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	req, err := c.requestBuilder.BuildGetCurrentHouseholdRequest(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building household retrieval request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetHousehold retrieves a household.
func (c *Client) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	req, err := c.requestBuilder.BuildGetHouseholdRequest(ctx, householdID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building household retrieval request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetHouseholds retrieves a list of households.
func (c *Client) GetHouseholds(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetHouseholdsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building household list request")
	}

	var apiResponse *types.APIResponse[[]*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving households")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	result := &types.QueryFilteredResult[types.Household]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return result, nil
}

// CreateHousehold creates a household.
func (c *Client) CreateHousehold(ctx context.Context, input *types.HouseholdCreationRequestInput) (*types.Household, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateHouseholdRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building household creation request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateHousehold updates a household.
func (c *Client) UpdateHousehold(ctx context.Context, household *types.Household) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if household == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, household.ID)

	req, err := c.requestBuilder.BuildUpdateHouseholdRequest(ctx, household)
	if err != nil {
		return observability.PrepareError(err, span, "building household update request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "updating household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveHousehold archives a household.
func (c *Client) ArchiveHousehold(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	req, err := c.requestBuilder.BuildArchiveHouseholdRequest(ctx, householdID)
	if err != nil {
		return observability.PrepareError(err, span, "building household archive request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "archiving household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// InviteUserToHousehold adds a user to a household.
func (c *Client) InviteUserToHousehold(ctx context.Context, destinationHouseholdID string, input *types.HouseholdInvitationCreationRequestInput) (*types.HouseholdInvitation, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, destinationHouseholdID)

	// we don't validate here because it needs to have the user ID

	req, err := c.requestBuilder.BuildInviteUserToHouseholdRequest(ctx, destinationHouseholdID, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building add user to household request")
	}

	var apiResponse *types.APIResponse[*types.HouseholdInvitation]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "adding user to household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// MarkAsDefault marks a given household as the default for a given user.
func (c *Client) MarkAsDefault(ctx context.Context, householdID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	req, err := c.requestBuilder.BuildMarkAsDefaultRequest(ctx, householdID)
	if err != nil {
		return observability.PrepareError(err, span, "building mark household as default request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "marking household as default")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RemoveUserFromHousehold removes a user from a household.
func (c *Client) RemoveUserFromHousehold(ctx context.Context, householdID, userID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	req, err := c.requestBuilder.BuildRemoveUserRequest(ctx, householdID, userID, "")
	if err != nil {
		return observability.PrepareError(err, span, "building remove user from household request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "removing user from household")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ModifyMemberPermissions modifies a given user's permissions for a given household.
func (c *Client) ModifyMemberPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if userID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildModifyMemberPermissionsRequest(ctx, householdID, userID, input)
	if err != nil {
		return observability.PrepareError(err, span, "building modify household member permissions request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "modifying user household permissions")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// TransferHouseholdOwnership transfers ownership of a household to a given user.
func (c *Client) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return ErrInvalidIDProvided
	}

	if input == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, "old_owner", input.CurrentOwner)
	tracing.AttachToSpan(span, "new_owner", input.NewOwner)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildTransferHouseholdOwnershipRequest(ctx, householdID, input)
	if err != nil {
		return observability.PrepareError(err, span, "building transfer household ownership request")
	}

	var apiResponse *types.APIResponse[*types.Household]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "transferring household to user")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
