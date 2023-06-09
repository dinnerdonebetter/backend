package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipeRating gets a recipe rating.
func (c *Client) GetRecipeRating(ctx context.Context, mealID, recipeRatingID string) (*types.RecipeRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if recipeRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRatingID)

	req, err := c.requestBuilder.BuildGetRecipeRatingRequest(ctx, mealID, recipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get recipe rating request")
	}

	var recipeRating *types.RecipeRating
	if err = c.fetchAndUnmarshal(ctx, req, &recipeRating); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe rating")
	}

	return recipeRating, nil
}

// GetRecipeRatings retrieves a list of recipe ratings.
func (c *Client) GetRecipeRatings(ctx context.Context, mealID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeRating], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	req, err := c.requestBuilder.BuildGetRecipeRatingsRequest(ctx, mealID, filter)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building recipe ratings list request")
	}

	var recipeRatings *types.QueryFilteredResult[types.RecipeRating]
	if err = c.fetchAndUnmarshal(ctx, req, &recipeRatings); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe ratings")
	}

	return recipeRatings, nil
}

// CreateRecipeRating creates a recipe rating.
func (c *Client) CreateRecipeRating(ctx context.Context, mealID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateRecipeRatingRequest(ctx, mealID, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create recipe rating request")
	}

	var recipeRating *types.RecipeRating
	if err = c.fetchAndUnmarshal(ctx, req, &recipeRating); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe rating")
	}

	return recipeRating, nil
}

// UpdateRecipeRating updates a recipe rating.
func (c *Client) UpdateRecipeRating(ctx context.Context, recipeRating *types.RecipeRating) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeRating == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRating.ID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRating.ID)

	req, err := c.requestBuilder.BuildUpdateRecipeRatingRequest(ctx, recipeRating)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building update recipe rating request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &recipeRating); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe rating %s", recipeRating.ID)
	}

	return nil
}

// ArchiveRecipeRating archives a recipe rating.
func (c *Client) ArchiveRecipeRating(ctx context.Context, mealID, recipeRatingID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachMealIDToSpan(span, mealID)

	if recipeRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachRecipeRatingIDToSpan(span, recipeRatingID)

	req, err := c.requestBuilder.BuildArchiveRecipeRatingRequest(ctx, mealID, recipeRatingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "building archive recipe rating request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe rating %s", recipeRatingID)
	}

	return nil
}