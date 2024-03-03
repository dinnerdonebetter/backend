package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetMealPlan gets a meal plan.
func (c *Client) GetMealPlan(ctx context.Context, mealPlanID string) (*types.MealPlan, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	req, err := c.requestBuilder.BuildGetMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get meal plan request")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetMealPlans retrieves a list of meal plans.
func (c *Client) GetMealPlans(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlan], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetMealPlansRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building meal plans list request")
	}

	var apiResponse *types.APIResponse[[]*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving meal plans")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.MealPlan]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateMealPlan creates a meal plan.
func (c *Client) CreateMealPlan(ctx context.Context, input *types.MealPlanCreationRequestInput) (*types.MealPlan, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateMealPlanRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create meal plan request")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateMealPlan updates a meal plan.
func (c *Client) UpdateMealPlan(ctx context.Context, mealPlan *types.MealPlan) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlan == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlan.ID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlan.ID)

	req, err := c.requestBuilder.BuildUpdateMealPlanRequest(ctx, mealPlan)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update meal plan request")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating meal plan %s", mealPlan.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveMealPlan archives a meal plan.
func (c *Client) ArchiveMealPlan(ctx context.Context, mealPlanID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	req, err := c.requestBuilder.BuildArchiveMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive meal plan request")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving meal plan %s", mealPlanID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// FinalizeMealPlan gets a meal plan.
func (c *Client) FinalizeMealPlan(ctx context.Context, mealPlanID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealPlanID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealPlanIDKey, mealPlanID)
	tracing.AttachToSpan(span, keys.MealPlanIDKey, mealPlanID)

	req, err := c.requestBuilder.BuildFinalizeMealPlanRequest(ctx, mealPlanID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building get meal plan request")
	}

	var apiResponse *types.APIResponse[*types.MealPlan]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "retrieving meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
