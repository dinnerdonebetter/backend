package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMeal gets a meal.
func (c *Client) GetMeal(ctx context.Context, mealID string) (*types.Meal, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	req, err := c.requestBuilder.BuildGetMealRequest(ctx, mealID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal request")
	}

	var apiResponse *types.APIResponse[*types.Meal]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal")
	}

	return apiResponse.Data, nil
}

// GetMeals retrieves a list of meals.
func (c *Client) GetMeals(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Meal], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building meals list request")
	}

	var apiResponse *types.APIResponse[[]*types.Meal]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meals")
	}

	response := &types.QueryFilteredResult[types.Meal]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// SearchForMeals retrieves a list of meals.
func (c *Client) SearchForMeals(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Meal], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(c.logger.Clone())

	tracing.AttachToSpan(span, keys.SearchQueryKey, query)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildSearchForMealsRequest(ctx, query, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building meals list request")
	}

	var apiResponse *types.APIResponse[[]*types.Meal]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meals")
	}

	response := &types.QueryFilteredResult[types.Meal]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateMeal creates a meal.
func (c *Client) CreateMeal(ctx context.Context, input *types.MealCreationRequestInput) (*types.Meal, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create meal request")
	}

	var apiResponse *types.APIResponse[*types.Meal]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal")
	}

	return apiResponse.Data, nil
}

// ArchiveMeal archives a meal.
func (c *Client) ArchiveMeal(ctx context.Context, mealID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	req, err := c.requestBuilder.BuildArchiveMealRequest(ctx, mealID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive meal request")
	}

	var apiResponse *types.APIResponse[*types.Meal]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal %s", mealID)
	}

	return nil
}
