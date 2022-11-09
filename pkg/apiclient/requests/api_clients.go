package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
)

const (
	apiClientsBasePath = "api_clients"
)

// BuildGetAPIClientRequest builds an HTTP request for fetching an API client.
func (b *Builder) BuildGetAPIClientRequest(ctx context.Context, clientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return nil, ErrInvalidIDProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		apiClientsBasePath,
		clientID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetAPIClientsRequest builds an HTTP request for fetching a list of API clients.
func (b *Builder) BuildGetAPIClientsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	uri := b.BuildURL(ctx, filter.ToValues(), apiClientsBasePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateAPIClientRequest builds an HTTP request for creating an API client.
func (b *Builder) BuildCreateAPIClientRequest(ctx context.Context, cookie *http.Cookie, input *types.APIClientCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if cookie == nil {
		return nil, ErrCookieRequired
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(ctx, nil, apiClientsBasePath)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	req.AddCookie(cookie)

	return req, nil
}

// BuildArchiveAPIClientRequest builds an HTTP request for archiving an API client.
func (b *Builder) BuildArchiveAPIClientRequest(ctx context.Context, clientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return nil, ErrInvalidIDProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		apiClientsBasePath,
		clientID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
