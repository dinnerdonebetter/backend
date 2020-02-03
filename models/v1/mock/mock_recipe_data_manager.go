package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeDataManager = (*RecipeDataManager)(nil)

// RecipeDataManager is a mocked models.RecipeDataManager for testing
type RecipeDataManager struct {
	mock.Mock
}

// GetRecipe is a mock function
func (m *RecipeDataManager) GetRecipe(ctx context.Context, recipeID, userID uint64) (*models.Recipe, error) {
	args := m.Called(ctx, recipeID, userID)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

// GetRecipeCount is a mock function
func (m *RecipeDataManager) GetRecipeCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllRecipesCount is a mock function
func (m *RecipeDataManager) GetAllRecipesCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetRecipes is a mock function
func (m *RecipeDataManager) GetRecipes(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.RecipeList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.RecipeList), args.Error(1)
}

// GetAllRecipesForUser is a mock function
func (m *RecipeDataManager) GetAllRecipesForUser(ctx context.Context, userID uint64) ([]models.Recipe, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Recipe), args.Error(1)
}

// CreateRecipe is a mock function
func (m *RecipeDataManager) CreateRecipe(ctx context.Context, input *models.RecipeCreationInput) (*models.Recipe, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

// UpdateRecipe is a mock function
func (m *RecipeDataManager) UpdateRecipe(ctx context.Context, updated *models.Recipe) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveRecipe is a mock function
func (m *RecipeDataManager) ArchiveRecipe(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
