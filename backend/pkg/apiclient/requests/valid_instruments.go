package requests

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	validInstrumentsBasePath = "valid_instruments"
)

// BuildGetValidInstrumentRequest builds an HTTP request for fetching a valid instrument.
func (b *Builder) BuildGetValidInstrumentRequest(ctx context.Context, validInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrumentID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRandomValidInstrumentRequest builds an HTTP request for fetching a valid instrument.
func (b *Builder) BuildGetRandomValidInstrumentRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		randomBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchValidInstrumentsRequest builds an HTTP request for querying valid instruments.
func (b *Builder) BuildSearchValidInstrumentsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		validInstrumentsBasePath,
		"search",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidInstrumentsRequest builds an HTTP request for fetching a list of valid instruments.
func (b *Builder) BuildGetValidInstrumentsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validInstrumentsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidInstrumentRequest builds an HTTP request for creating a valid instrument.
func (b *Builder) BuildCreateValidInstrumentRequest(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*http.Request, error) {
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
		validInstrumentsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidInstrumentRequest builds an HTTP request for updating a valid instrument.
func (b *Builder) BuildUpdateValidInstrumentRequest(ctx context.Context, validInstrument *types.ValidInstrument) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validInstrument == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrument.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrument.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertValidInstrumentToValidInstrumentUpdateRequestInput(validInstrument)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidInstrumentRequest builds an HTTP request for archiving a valid instrument.
func (b *Builder) BuildArchiveValidInstrumentRequest(ctx context.Context, validInstrumentID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validInstrumentID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	uri := b.BuildURL(
		ctx,
		nil,
		validInstrumentsBasePath,
		validInstrumentID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
