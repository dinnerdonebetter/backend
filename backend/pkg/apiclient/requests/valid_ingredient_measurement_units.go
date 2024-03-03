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
	validIngredientMeasurementUnitsBasePath = "valid_ingredient_measurement_units"
)

// BuildGetValidIngredientMeasurementUnitRequest builds an HTTP request for fetching a valid ingredient measurement unit.
func (b *Builder) BuildGetValidIngredientMeasurementUnitRequest(ctx context.Context, validIngredientMeasurementUnitID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientMeasurementUnitsBasePath,
		validIngredientMeasurementUnitID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientMeasurementUnitsRequest builds an HTTP request for fetching a list of valid ingredient measurement units.
func (b *Builder) BuildGetValidIngredientMeasurementUnitsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientMeasurementUnitsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientMeasurementUnitsForIngredientRequest builds an HTTP request for fetching a list of valid ingredient measurement units.
func (b *Builder) BuildGetValidIngredientMeasurementUnitsForIngredientRequest(ctx context.Context, ingredientID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, ingredientID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientMeasurementUnitsBasePath,
		"by_ingredient",
		ingredientID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetValidIngredientMeasurementUnitsForMeasurementUnitRequest builds an HTTP request for fetching a list of valid ingredient measurement units.
func (b *Builder) BuildGetValidIngredientMeasurementUnitsForMeasurementUnitRequest(ctx context.Context, validMeasurementUnitID string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		validIngredientMeasurementUnitsBasePath,
		"by_measurement_unit",
		validMeasurementUnitID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateValidIngredientMeasurementUnitRequest builds an HTTP request for creating a valid ingredient measurement unit.
func (b *Builder) BuildCreateValidIngredientMeasurementUnitRequest(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*http.Request, error) {
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
		validIngredientMeasurementUnitsBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateValidIngredientMeasurementUnitRequest builds an HTTP request for updating a valid ingredient measurement unit.
func (b *Builder) BuildUpdateValidIngredientMeasurementUnitRequest(ctx context.Context, validIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientMeasurementUnit == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnit.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientMeasurementUnitsBasePath,
		validIngredientMeasurementUnit.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitUpdateRequestInput(validIngredientMeasurementUnit)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveValidIngredientMeasurementUnitRequest builds an HTTP request for archiving a valid ingredient measurement unit.
func (b *Builder) BuildArchiveValidIngredientMeasurementUnitRequest(ctx context.Context, validIngredientMeasurementUnitID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	uri := b.BuildURL(
		ctx,
		nil,
		validIngredientMeasurementUnitsBasePath,
		validIngredientMeasurementUnitID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
