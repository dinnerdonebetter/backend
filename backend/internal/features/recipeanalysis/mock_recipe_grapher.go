package recipeanalysis

import (
	"context"
	"image"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
	"gonum.org/v1/gonum/graph/simple"
)

var _ RecipeAnalyzer = (*MockRecipeAnalyzer)(nil)

// MockRecipeAnalyzer is a mock RecipeAnalyzer.
type MockRecipeAnalyzer struct {
	mock.Mock
}

// MakeGraphForRecipe implements our interface.
func (m *MockRecipeAnalyzer) MakeGraphForRecipe(ctx context.Context, recipe *types.Recipe) (*simple.DirectedGraph, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(*simple.DirectedGraph), returnArgs.Error(1)
}

// GenerateMealPlanTasksForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateMealPlanTasksForRecipe(ctx context.Context, mealPlanOptionID string, recipe *types.Recipe) ([]*types.MealPlanTaskDatabaseCreationInput, error) {
	returnArgs := m.Called(ctx, mealPlanOptionID, recipe)

	return returnArgs.Get(0).([]*types.MealPlanTaskDatabaseCreationInput), returnArgs.Error(1)
}

// FindStepsEligibleForMealPlanTasks implements our interface.
func (m *MockRecipeAnalyzer) FindStepsEligibleForMealPlanTasks(ctx context.Context, recipe *types.Recipe) ([]*types.RecipeStep, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).([]*types.RecipeStep), returnArgs.Error(1)
}

// GenerateDAGDiagramForRecipe implements our interface.
func (m *MockRecipeAnalyzer) GenerateDAGDiagramForRecipe(ctx context.Context, recipe *types.Recipe) (image.Image, error) {
	returnArgs := m.Called(ctx, recipe)

	return returnArgs.Get(0).(image.Image), returnArgs.Error(1)
}

// RenderMermaidDiagramForRecipe implements our interface.
func (m *MockRecipeAnalyzer) RenderMermaidDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string {
	return m.Called(ctx, recipe).String(0)
}

// RenderGraphvizDiagramForRecipe implements our interface.
func (m *MockRecipeAnalyzer) RenderGraphvizDiagramForRecipe(ctx context.Context, recipe *types.Recipe) string {
	return m.Called(ctx, recipe).String(0)
}
