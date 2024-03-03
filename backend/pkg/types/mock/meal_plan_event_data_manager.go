package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.MealPlanEventDataManager = (*MealPlanEventDataManagerMock)(nil)

// MealPlanEventDataManagerMock is a mocked types.MealPlanEventDataManager for testing.
type MealPlanEventDataManagerMock struct {
	mock.Mock
}

// MealPlanEventIsEligibleForVoting is a mock function.
func (m *MealPlanEventDataManagerMock) MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID)
	return args.Bool(0), args.Error(1)
}

// MealPlanEventExists is a mock function.
func (m *MealPlanEventDataManagerMock) MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID)
	return args.Bool(0), args.Error(1)
}

// GetMealPlanEvent is a mock function.
func (m *MealPlanEventDataManagerMock) GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*types.MealPlanEvent, error) {
	args := m.Called(ctx, mealPlanID, mealPlanEventID)
	return args.Get(0).(*types.MealPlanEvent), args.Error(1)
}

// GetMealPlanEvents is a mock function.
func (m *MealPlanEventDataManagerMock) GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlanEvent], error) {
	args := m.Called(ctx, mealPlanID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.MealPlanEvent]), args.Error(1)
}

// CreateMealPlanEvent is a mock function.
func (m *MealPlanEventDataManagerMock) CreateMealPlanEvent(ctx context.Context, input *types.MealPlanEventDatabaseCreationInput) (*types.MealPlanEvent, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.MealPlanEvent), args.Error(1)
}

// UpdateMealPlanEvent is a mock function.
func (m *MealPlanEventDataManagerMock) UpdateMealPlanEvent(ctx context.Context, updated *types.MealPlanEvent) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveMealPlanEvent is a mock function.
func (m *MealPlanEventDataManagerMock) ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error {
	return m.Called(ctx, mealPlanID, mealPlanEventID).Error(0)
}

// AttemptToFinalizeMealPlanEvent is a mock function.
func (m *MealPlanEventDataManagerMock) AttemptToFinalizeMealPlanEvent(ctx context.Context, mealPlanEventID string) (changed bool, err error) {
	args := m.Called(ctx, mealPlanEventID)
	return args.Bool(0), args.Error(1)
}

// GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods is a mock function.
func (m *MealPlanEventDataManagerMock) GetUnfinalizedMealPlanEventsWithExpiredVotingPeriods(ctx context.Context) ([]*types.MealPlanEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*types.MealPlanEvent), args.Error(1)
}

// GetFinalizedMealPlanEventIDsForTheNextWeek is a mock function.
func (m *MealPlanEventDataManagerMock) GetFinalizedMealPlanEventIDsForTheNextWeek(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}
