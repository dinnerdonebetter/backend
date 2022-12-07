package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

const (
	validIngredientStateIngredientsBasePath = "valid_ingredient_state_ingredients"
)

// BuildGetValidIngredientStateIngredientRequest builds an HTTP request for fetching a valid ingredient preparation.
func (b *Builder) BuildGetValidIngredientStateIngredientRequest(ctx context.Context, validIngredientStateIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientStateIngredientsBasePath,
		validIngredientStateIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientStateIngredientsRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientStateIngredientsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientStateIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientStateIngredientsForIngredientRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientStateIngredientsForIngredientRequest(ctx context.Context, ingredientID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientIDToSpan(span, ingredientID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientStateIngredientsBasePath,
		"by_ingredient",
		ingredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientStateIngredientsForPreparationRequest builds an HTTP request for fetching a list of valid ingredient preparations.
func (b *Builder) BuildGetValidIngredientStateIngredientsForPreparationRequest(ctx context.Context, ingredientState string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if ingredientState == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientStateIDToSpan(span, ingredientState)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientStateIngredientsBasePath,
		"by_ingredient_state",
		ingredientState,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidIngredientStateIngredientRequest builds an HTTP request for creating a valid ingredient preparation.
func (b *Builder) BuildCreateValidIngredientStateIngredientRequest(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*http.Request, error) {
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
		validIngredientStateIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientStateIngredientRequest builds an HTTP request for updating a valid ingredient preparation.
func (b *Builder) BuildUpdateValidIngredientStateIngredientRequest(ctx context.Context, validIngredientStateIngredient *types.ValidIngredientStateIngredient) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientStateIngredient == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredient.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientStateIngredientsBasePath,
		validIngredientStateIngredient.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	input := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientUpdateRequestInput(validIngredientStateIngredient)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientStateIngredientRequest builds an HTTP request for archiving a valid ingredient preparation.
func (b *Builder) BuildArchiveValidIngredientStateIngredientRequest(ctx context.Context, validIngredientStateIngredientID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientStateIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidIngredientStateIngredientIDToSpan(span, validIngredientStateIngredientID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientStateIngredientsBasePath,
		validIngredientStateIngredientID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}