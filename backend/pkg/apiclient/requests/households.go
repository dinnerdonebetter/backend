package requests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	householdsBasePath = "households"
)

// BuildSwitchActiveHouseholdRequest builds an HTTP request for switching active households.
func (b *Builder) BuildSwitchActiveHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	uri := b.buildAPIV1URL(ctx, nil, usersBasePath, "household", "select").String()

	input := &types.ChangeActiveHouseholdInput{
		HouseholdID: householdID,
	}

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildGetCurrentHouseholdRequest builds an HTTP request for fetching a household.
func (b *Builder) BuildGetCurrentHouseholdRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		"current",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetHouseholdRequest builds an HTTP request for fetching a household.
func (b *Builder) BuildGetHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetHouseholdsRequest builds an HTTP request for fetching a list of households.
func (b *Builder) BuildGetHouseholdsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(ctx, filter.ToValues(), householdsBasePath)

	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateHouseholdRequest builds an HTTP request for creating a household.
func (b *Builder) BuildCreateHouseholdRequest(ctx context.Context, input *types.HouseholdCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildUpdateHouseholdRequest builds an HTTP request for updating a household.
func (b *Builder) BuildUpdateHouseholdRequest(ctx context.Context, household *types.Household) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if household == nil {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		household.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertHouseholdToHouseholdUpdateRequestInput(household)

	return b.buildDataRequest(ctx, http.MethodPut, uri, input)
}

// BuildArchiveHouseholdRequest builds an HTTP request for archiving a household.
func (b *Builder) BuildArchiveHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildInviteUserToHouseholdRequest builds a request that adds a user to a household.
func (b *Builder) BuildInviteUserToHouseholdRequest(ctx context.Context, destinationHouseholdID string, input *types.HouseholdInvitationCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if destinationHouseholdID == "" {
		return nil, ErrInvalidIDProvided
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// we don't validate here because it needs to have the user ID

	uri := b.BuildURL(ctx, nil, householdsBasePath, destinationHouseholdID, "invite")
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildMarkAsDefaultRequest builds a request that marks a given household as the default for a given user.
func (b *Builder) BuildMarkAsDefaultRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "default")
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildRemoveUserRequest builds a request that removes a user from a household.
func (b *Builder) BuildRemoveUserRequest(ctx context.Context, householdID, userID, reason string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	u := b.buildAPIV1URL(ctx, nil, householdsBasePath, householdID, "members", userID)

	if reason != "" {
		q := u.Query()
		q.Set("reason", reason)
		u.RawQuery = q.Encode()
	}

	tracing.AttachToSpan(span, keys.RequestURIKey, u.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildModifyMemberPermissionsRequest builds a request that modifies a given user's permissions for a given household.
func (b *Builder) BuildModifyMemberPermissionsRequest(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "members", userID, "permissions")
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	return b.buildDataRequest(ctx, http.MethodPatch, uri, input)
}

// BuildTransferHouseholdOwnershipRequest builds a request that transfers ownership of a household to a given user.
func (b *Builder) BuildTransferHouseholdOwnershipRequest(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, fmt.Errorf("householdID: %w", ErrInvalidIDProvided)
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "transfer")
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}
