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
	validIngredientGroupsBasePath = "valid_ingredient_groups"
)

// BuildGetValidIngredientGroupRequest builds an HTTP request for fetching a valid ingredient.
func (b *Builder) BuildGetValidIngredientGroupRequest(ctx context.Context, validIngredientGroupID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientGroupID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientGroupsBasePath,
		validIngredientGroupID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchValidIngredientGroupsRequest builds an HTTP request for querying valid ingredients.
func (b *Builder) BuildSearchValidIngredientGroupsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	params := url.Values{}
	params.Set(types.SearchQueryKey, query)
	params.Set(types.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := b.BuildURL(
		ctx,
		params,
		validIngredientGroupsBasePath,
		"search",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientGroupsRequest builds an HTTP request for fetching a list of valid ingredients.
func (b *Builder) BuildGetValidIngredientGroupsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientGroupsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidIngredientGroupRequest builds an HTTP request for creating a valid ingredient.
func (b *Builder) BuildCreateValidIngredientGroupRequest(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*http.Request, error) {
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
		validIngredientGroupsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientGroupRequest builds an HTTP request for updating a valid ingredient.
func (b *Builder) BuildUpdateValidIngredientGroupRequest(ctx context.Context, validIngredientGroup *types.ValidIngredientGroup) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientGroup == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroup.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientGroupsBasePath,
		validIngredientGroup.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(validIngredientGroup)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientGroupRequest builds an HTTP request for archiving a valid ingredient.
func (b *Builder) BuildArchiveValidIngredientGroupRequest(ctx context.Context, validIngredientGroupID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientGroupID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientGroupIDKey, validIngredientGroupID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientGroupsBasePath,
		validIngredientGroupID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
