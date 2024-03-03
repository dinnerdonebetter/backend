package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipeStep gets a recipe step.
func (c *Client) GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, buildInvalidIDError("recipe step")
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	req, err := c.requestBuilder.BuildGetRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe step request")
	}

	var apiResponse *types.APIResponse[*types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe step")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// GetRecipeSteps retrieves a list of recipe steps.
func (c *Client) GetRecipeSteps(ctx context.Context, recipeID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStep], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	req, err := c.requestBuilder.BuildGetRecipeStepsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe steps list request")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe steps")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	response := &types.QueryFilteredResult[types.RecipeStep]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateRecipeStep creates a recipe step.
func (c *Client) CreateRecipeStep(ctx context.Context, recipeID string, input *types.RecipeStepCreationRequestInput) (*types.RecipeStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeID == "" {
		return nil, buildInvalidIDError("recipe")
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeStepRequest(ctx, recipeID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe step request")
	}

	var apiResponse *types.APIResponse[*types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe step")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// UpdateRecipeStep updates a recipe step.
func (c *Client) UpdateRecipeStep(ctx context.Context, recipeStep *types.RecipeStep) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeStep == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStep.ID)
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStep.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeStepRequest(ctx, recipeStep)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe step request")
	}

	var apiResponse *types.APIResponse[*types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe step %s", recipeStep.ID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// ArchiveRecipeStep archives a recipe step.
func (c *Client) ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	req, err := c.requestBuilder.BuildArchiveRecipeStepRequest(ctx, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive recipe step request")
	}

	var apiResponse *types.APIResponse[*types.RecipeStep]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "archiving recipe step %s", recipeStepID)
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// UploadRecipeMediaForStep uploads a new avatar.
func (c *Client) UploadRecipeMediaForStep(ctx context.Context, files map[string][]byte, recipeID, recipeStepID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	if files == nil {
		return ErrNilInputProvided
	}

	req, err := c.requestBuilder.BuildMultipleRecipeMediaUploadRequestForRecipeStep(ctx, files, recipeID, recipeStepID)
	if err != nil {
		return observability.PrepareError(err, span, "building media upload request")
	}

	var apiResponse *types.APIResponse[[]*types.RecipeMedia]
	if err = c.fetchAndUnmarshal(ctx, req, &apiResponse); err != nil {
		return observability.PrepareError(err, span, "uploading media")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
