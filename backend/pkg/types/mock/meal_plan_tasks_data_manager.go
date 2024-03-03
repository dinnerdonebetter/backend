package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanTaskDataManager = (*MealPlanTaskDataManagerMock)(nil)

// MealPlanTaskDataManagerMock is a mocked types.MealPlanTaskDataManager for testing.
type MealPlanTaskDataManagerMock struct {
	mock.Mock
}

// MarkMealPlanAsHavingTasksCreated is a mock function.
func (m *MealPlanTaskDataManagerMock) MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error {
	return m.Called(ctx, mealPlanID).Error(0)
}

// MealPlanTaskExists is a mock function.
func (m *MealPlanTaskDataManagerMock) MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanTaskID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanTask is a mock function.
func (m *MealPlanTaskDataManagerMock) GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanTaskID)
	return args.Get(0).(*types.MealPlanTask), args.Error(1)
}

// CreateMealPlanTask is a mock function.
func (m *MealPlanTaskDataManagerMock) CreateMealPlanTask(ctx context.Context, input *types.MealPlanTaskDatabaseCreationInput) (*types.MealPlanTask, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlanTask), args.Error(1)
}

// GetMealPlanTasksForMealPlan is a mock function.
func (m *MealPlanTaskDataManagerMock) GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*types.MealPlanTask, error) {
	args := m.Called(ctx, mealPlanID)
	return args.Get(0).([]*types.MealPlanTask), args.Error(1)
}

// CreateMealPlanTasksForMealPlanOption is a mock function.
func (m *MealPlanTaskDataManagerMock) CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*types.MealPlanTaskDatabaseCreationInput) ([]*types.MealPlanTask, error) {
	args := m.Called(ctx, inputs)
	return args.Get(0).([]*types.MealPlanTask), args.Error(1)
}

// ChangeMealPlanTaskStatus is a mock function.
func (m *MealPlanTaskDataManagerMock) ChangeMealPlanTaskStatus(ctx context.Context, input *types.MealPlanTaskStatusChangeRequestInput) error {
	return m.Called(ctx, input).Error(0)
}
