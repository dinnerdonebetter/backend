package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// RunFinalizeMealPlansWorker runs a worker.
func (c *Client) RunFinalizeMealPlansWorker(ctx context.Context, input *types.FinalizeMealPlansRequest) (*types.FinalizeMealPlansResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	req, err := c.requestBuilder.BuildRunFinalizeMealPlansWorkerRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building worker execution request")
	}

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// RunMealPlanGroceryListInitializationWorker runs a worker.
func (c *Client) RunMealPlanGroceryListInitializationWorker(ctx context.Context) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	req, err := c.requestBuilder.BuildRunMealPlanGroceryListInitializationWorkerRequest(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building worker execution request")
	}

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "initializing meal plan grocery lists")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RunMealPlanTaskCreationWorker runs a worker.
func (c *Client) RunMealPlanTaskCreationWorker(ctx context.Context) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	req, err := c.requestBuilder.BuildRunMealPlanTaskCreationWorkerRequest(ctx)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building worker execution request")
	}

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating meal plan tasks")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
