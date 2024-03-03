package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealPlanEvent gets a meal plan event.
func (c *Client) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanEventID)

	if mealPlanEventID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	req, err := c.requestBuilder.BuildGetMealPlanEventRequest(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building get meal plan event request")
	}

	var apiResponse *types.APIResponse[*types.MealPlanEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan event")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetMealPlanEvents retrieves a list of meal plan events.
func (c *Client) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlanEvent], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanEventsRequest(ctx, mealPlanID, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building meal plan events list request")
	}

	var apiResponse *types.APIResponse[[]*types.MealPlanEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving meal plan events")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.MealPlanEvent]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateMealPlanEvent creates a meal plan event.
func (c *Client) CreateMealPlanEvent(ctx context.Context, mealPlanID string, input *types.MealPlanEventCreationRequestInput) (*types.MealPlanEvent, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanEventRequest(ctx, mealPlanID, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create meal plan event request")
	}

	var apiResponse *types.APIResponse[*types.MealPlanEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareError(err, span, "creating meal plan event")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateMealPlanEvent updates a meal plan event.
func (c *Client) UpdateMealPlanEvent(ctx context.Context, mealPlanEvent *types.MealPlanEvent) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanEvent == nil {
		return ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEvent.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanEventRequest(ctx, mealPlanEvent)
	if err != nil {
		return observability.PrepareError(err, span, "building update meal plan event request")
	}

	var apiResponse *types.APIResponse[*types.MealPlanEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "updating meal plan event %s", mealPlanEvent.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveMealPlanEvent archives a meal plan event.
func (c *Client) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanEventID)

	if mealPlanEventID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)

	req, err := c.requestBuilder.BuildArchiveMealPlanEventRequest(ctx, mealPlanID, mealPlanEventID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive meal plan event request")
	}

	var apiResponse *types.APIResponse[*types.MealPlanEvent]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "archiving meal plan event %s", mealPlanEventID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
