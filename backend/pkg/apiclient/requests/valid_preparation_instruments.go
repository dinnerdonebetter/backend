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
	validPreparationInstrumentsBasePath = "valid_preparation_instruments"
)

// BuildGetValidPreparationInstrumentRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrumentID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidPreparationInstrumentsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsForPreparationRequest builds an HTTP request for fetching a list of valid preparation instruments.
func (b *Builder) BuildGetValidPreparationInstrumentsForPreparationRequest(ctx context.Context, validPreparationID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validPreparationID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
		"by_preparation",
		validPreparationID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidPreparationInstrumentsForInstrumentRequest builds an HTTP request for fetching a list of valid preparation instruments.
func (b *Builder) BuildGetValidPreparationInstrumentsForInstrumentRequest(ctx context.Context, validInstrumentID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validPreparationInstrumentsBasePath,
		"by_instrument",
		validInstrumentID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		logger.Error(err, "building request")
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidPreparationInstrumentRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidPreparationInstrumentRequest(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*http.Request, error) {
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
		validPreparationInstrumentsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidPreparationInstrumentRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrument *types.ValidPreparationInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrument == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrument.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentUpdateRequestInput(validPreparationInstrument)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidPreparationInstrumentRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidPreparationInstrumentRequest(ctx context.Context, validPreparationInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validPreparationInstrumentsBasePath,
		validPreparationInstrumentID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
