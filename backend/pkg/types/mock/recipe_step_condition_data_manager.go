package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.RecipeStepCompletionConditionDataManager = (*RecipeStepCompletionConditionDataManagerMock)(nil)

// RecipeStepCompletionConditionDataManagerMock is a mocked types.RecipeStepCompletionConditionDataManager for testing.
type RecipeStepCompletionConditionDataManagerMock struct {
	mock.Mock
}

// RecipeStepCompletionConditionExists is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*types.RecipeStepCompletionCondition, error) {
	args := m.Called(ctx, recipeID, recipeStepID, recipeStepIngredientID)
	return args.Get(0).(*types.RecipeStepCompletionCondition), args.Error(1)
}

// GetRecipeStepCompletionConditions is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeStepCompletionCondition], error) {
	args := m.Called(ctx, recipeID, recipeStepID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.RecipeStepCompletionCondition]), args.Error(1)
}

// CreateRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) CreateRecipeStepCompletionCondition(ctx context.Context, input *types.RecipeStepCompletionConditionDatabaseCreationInput) (*types.RecipeStepCompletionCondition, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.RecipeStepCompletionCondition), args.Error(1)
}

// UpdateRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) UpdateRecipeStepCompletionCondition(ctx context.Context, updated *types.RecipeStepCompletionCondition) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipeStepCompletionCondition is a mock function.
func (m *RecipeStepCompletionConditionDataManagerMock) ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error {
	return m.Called(ctx, recipeStepID, recipeStepIngredientID).Error(0)
}
