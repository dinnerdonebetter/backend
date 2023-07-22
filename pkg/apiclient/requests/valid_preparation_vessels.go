package requests

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	validPreparationVesselsBasePath = "valid_preparation_vessels"
)

// BuildGetValidPreparationVesselRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidPreparationVesselRequest(ctx context.Context, validPreparationVesselID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationVesselsBasePath,
		validPreparationVesselID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationVesselsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidPreparationVesselsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationVesselsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationVesselsForPreparationRequest builds an HTTP request for fetching a list of valid preparation vessels.
func (b *Builder) BuildGetValidPreparationVesselsForPreparationRequest(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, validPreparationID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationVesselsBasePath,
		"by_preparation",
		validPreparationID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationVesselsForVesselRequest builds an HTTP request for fetching a list of valid preparation vessels.
func (b *Builder) BuildGetValidPreparationVesselsForVesselRequest(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationVesselsBasePath,
		"by_vessel",
		validInstrumentID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidPreparationVesselRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidPreparationVesselRequest(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*http.Request, error) {
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
		validPreparationVesselsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidPreparationVesselRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidPreparationVesselRequest(ctx context.Context, validPreparationVessel *types.ValidPreparationVessel) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVessel == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVessel.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationVesselsBasePath,
		validPreparationVessel.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidPreparationVesselToValidPreparationVesselUpdateRequestInput(validPreparationVessel)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationVesselRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidPreparationVesselRequest(ctx context.Context, validPreparationVesselID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationVesselsBasePath,
		validPreparationVesselID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
