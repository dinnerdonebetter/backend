package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	householdInstrumentOwnershipsBasePath = "instruments"
)

// BuildGetHouseholdInstrumentOwnershipRequest builds an HTTP request for fetching a household instrument ownership.
func (b *Builder) BuildGetHouseholdInstrumentOwnershipRequest(ctx context.Context, householdInstrumentOwnershipID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdInstrumentOwnershipID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdInstrumentOwnershipsBasePath,
		householdInstrumentOwnershipID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetHouseholdInstrumentOwnershipsRequest builds an HTTP request for fetching a list of household instrument ownerships.
func (b *Builder) BuildGetHouseholdInstrumentOwnershipsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		householdsBasePath,
		householdInstrumentOwnershipsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateHouseholdInstrumentOwnershipRequest builds an HTTP request for creating a household instrument ownership.
func (b *Builder) BuildCreateHouseholdInstrumentOwnershipRequest(ctx context.Context, input *types.HouseholdInstrumentOwnershipCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdInstrumentOwnershipsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateHouseholdInstrumentOwnershipRequest builds an HTTP request for updating a household instrument ownership.
func (b *Builder) BuildUpdateHouseholdInstrumentOwnershipRequest(ctx context.Context, householdInstrumentOwnership *types.HouseholdInstrumentOwnership) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdInstrumentOwnership == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnership.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdInstrumentOwnershipsBasePath,
		householdInstrumentOwnership.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput(householdInstrumentOwnership)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveHouseholdInstrumentOwnershipRequest builds an HTTP request for archiving a household instrument ownership.
func (b *Builder) BuildArchiveHouseholdInstrumentOwnershipRequest(ctx context.Context, householdInstrumentOwnershipID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdInstrumentOwnershipID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdInstrumentOwnershipIDKey, householdInstrumentOwnershipID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdInstrumentOwnershipsBasePath,
		householdInstrumentOwnershipID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
