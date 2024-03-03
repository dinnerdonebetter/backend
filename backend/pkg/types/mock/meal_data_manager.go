package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealDataManager = (*MealDataManagerMock)(nil)

// MealDataManagerMock is a mocked types.MealDataManager for testing.
type MealDataManagerMock struct {
	mock.Mock
}

// MealExists is a mock function.
func (m *MealDataManagerMock) MealExists(ctx context.Context, recipeID string) (bool, error) {
	args := m.Called(ctx, recipeID)
	return args.Bool(0), args.Error(1)
}

// GetMeal is a mock function.
func (m *MealDataManagerMock) GetMeal(ctx context.Context, recipeID string) (*types.Meal, error) {
	args := m.Called(ctx, recipeID)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// GetMealByIDAndUser is a mock function.
func (m *MealDataManagerMock) GetMealByIDAndUser(ctx context.Context, recipeID, userID string) (*types.Meal, error) {
	args := m.Called(ctx, recipeID, userID)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// GetMeals is a mock function.
func (m *MealDataManagerMock) GetMeals(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Meal], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.Meal]), args.Error(1)
}

// SearchForMeals is a mock function.
func (m *MealDataManagerMock) SearchForMeals(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Meal], error) {
	args := m.Called(ctx, query, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.Meal]), args.Error(1)
}

// CreateMeal is a mock function.
func (m *MealDataManagerMock) CreateMeal(ctx context.Context, input *types.MealDatabaseCreationInput) (*types.Meal, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Meal), args.Error(1)
}

// UpdateMeal is a mock function.
func (m *MealDataManagerMock) UpdateMeal(ctx context.Context, updated *types.Meal) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMeal is a mock function.
func (m *MealDataManagerMock) ArchiveMeal(ctx context.Context, recipeID, householdID string) error {
	return m.Called(ctx, recipeID, householdID).Error(0)
}

// MarkMealAsIndexed is a mock function.
func (m *MealDataManagerMock) MarkMealAsIndexed(ctx context.Context, mealID string) error {
	return m.Called(ctx, mealID).Error(0)
}

// GetMealIDsThatNeedSearchIndexing is a mock function.
func (m *MealDataManagerMock) GetMealIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetMealsWithIDs is a mock function.
func (m *MealDataManagerMock) GetMealsWithIDs(ctx context.Context, ids []string) ([]*types.Meal, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.Meal), returnValues.Error(1)
}
