package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanGroceryListItemDataManager = (*MealPlanGroceryListItemDataManagerMock)(nil)

// MealPlanGroceryListItemDataManagerMock is a mocked types.MealPlanGroceryListItemDataManager for testing.
type MealPlanGroceryListItemDataManagerMock struct {
	mock.Mock
}

// MealPlanGroceryListItemExists is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) MealPlanGroceryListItemExists(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (bool, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) GetMealPlanGroceryListItem(ctx context.Context, mealPlanID, mealPlanGroceryListItemID string) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID, mealPlanGroceryListItemID)
	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// GetMealPlanGroceryListItemsForMealPlan is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, mealPlanID)
	return returnValues.Get(0).([]*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// CreateMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) CreateMealPlanGroceryListItem(ctx context.Context, input *types.MealPlanGroceryListItemDatabaseCreationInput) (*types.MealPlanGroceryListItem, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.MealPlanGroceryListItem), returnValues.Error(1)
}

// UpdateMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) UpdateMealPlanGroceryListItem(ctx context.Context, updated *types.MealPlanGroceryListItem) error {
	returnValues := m.Called(ctx, updated)
	return returnValues.Error(0)
}

// ArchiveMealPlanGroceryListItem is a mock function.
func (m *MealPlanGroceryListItemDataManagerMock) ArchiveMealPlanGroceryListItem(ctx context.Context, mealPlanGroceryListItemID string) error {
	returnValues := m.Called(ctx, mealPlanGroceryListItemID)
	return returnValues.Error(0)
}
