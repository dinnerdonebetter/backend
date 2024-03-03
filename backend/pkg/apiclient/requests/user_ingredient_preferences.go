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
	userIngredientPreferencesBasePath = "user_ingredient_preferences"
)

// BuildGetUserIngredientPreferencesRequest builds an HTTP request for fetching a list of valid preparations.
func (b *Builder) BuildGetUserIngredientPreferencesRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		userIngredientPreferencesBasePath,
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

// BuildCreateUserIngredientPreferenceRequest builds an HTTP request for creating a valid preparation.
func (b *Builder) BuildCreateUserIngredientPreferenceRequest(ctx context.Context, input *types.UserIngredientPreferenceCreationRequestInput) (*http.Request, error) {
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
		userIngredientPreferencesBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateUserIngredientPreferenceRequest builds an HTTP request for updating a valid preparation.
func (b *Builder) BuildUpdateUserIngredientPreferenceRequest(ctx context.Context, userIngredientPreference *types.UserIngredientPreference) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if userIngredientPreference == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, userIngredientPreference.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		userIngredientPreferencesBasePath,
		userIngredientPreference.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(userIngredientPreference)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveUserIngredientPreferenceRequest builds an HTTP request for archiving a valid preparation.
func (b *Builder) BuildArchiveUserIngredientPreferenceRequest(ctx context.Context, userIngredientPreferenceID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if userIngredientPreferenceID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIngredientPreferenceIDKey, userIngredientPreferenceID)

	uri := b.BuildURL(
		ctx,
		nil,
		userIngredientPreferencesBasePath,
		userIngredientPreferenceID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
