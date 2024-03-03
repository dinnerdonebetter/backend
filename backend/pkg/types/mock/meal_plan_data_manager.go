package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanDataManager = (*MealPlanDataManagerMock)(nil)

// MealPlanDataManagerMock is a mocked types.MealPlanDataManager for testing.
type MealPlanDataManagerMock struct {
	mock.Mock
}

// MealPlanExists is a mock function.
func (m *MealPlanDataManagerMock) MealPlanExists(ctx context.Context, mealPlanID, householdID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlan is a mock function.
func (m *MealPlanDataManagerMock) GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*types.MealPlan, error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Get(0).(*types.MealPlan), args.Error(1)
}

// GetMealPlans is a mock function.
func (m *MealPlanDataManagerMock) GetMealPlans(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlan], error) {
	args := m.Called(ctx, householdID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.MealPlan]), args.Error(1)
}

// CreateMealPlan is a mock function.
func (m *MealPlanDataManagerMock) CreateMealPlan(ctx context.Context, input *types.MealPlanDatabaseCreationInput) (*types.MealPlan, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlan), args.Error(1)
}

// UpdateMealPlan is a mock function.
func (m *MealPlanDataManagerMock) UpdateMealPlan(ctx context.Context, updated *types.MealPlan) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlan is a mock function.
func (m *MealPlanDataManagerMock) ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error {
	return m.Called(ctx, mealPlanID, householdID).Error(0)
}

// AttemptToFinalizeMealPlan is a mock function.
func (m *MealPlanDataManagerMock) AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (changed bool, err error) {
	args := m.Called(ctx, mealPlanID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetUnfinalizedMealPlansWithExpiredVotingPeriods is a mock function.
func (m *MealPlanDataManagerMock) GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlan), args.Error(1)
}

// GetFinalizedMealPlanIDsForTheNextWeek is a mock function.
func (m *MealPlanDataManagerMock) GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*types.FinalizedMealPlanDatabaseResult, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.FinalizedMealPlanDatabaseResult), args.Error(1)
}

// GetFinalizedMealPlansWithUninitializedGroceryLists is a mock function.
func (m *MealPlanDataManagerMock) GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*types.MealPlan, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlan), args.Error(1)
}
