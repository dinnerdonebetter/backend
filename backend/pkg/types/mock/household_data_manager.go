package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdDataManager = (*HouseholdDataManagerMock)(nil)

// HouseholdDataManagerMock is a mocked types.HouseholdDataManager for testing.
type HouseholdDataManagerMock struct {
	mock.Mock
}

// HouseholdExists is a mock function.
func (m *HouseholdDataManagerMock) HouseholdExists(ctx context.Context, householdID, userID string) (bool, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Bool(0), args.Error(1)
}

// GetHousehold is a mock function.
func (m *HouseholdDataManagerMock) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	args := m.Called(ctx, householdID)
	return args.Get(0).(*types.Household), args.Error(1)
}

// GetAllHouseholds is a mock function.
func (m *HouseholdDataManagerMock) GetAllHouseholds(ctx context.Context, results chan []*types.Household, bucketSize uint16) error {
	return m.Called(ctx, results, bucketSize).Error(0)
}

// GetHouseholds is a mock function.
func (m *HouseholdDataManagerMock) GetHouseholds(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.Household]), args.Error(1)
}

// CreateHousehold is a mock function.
func (m *HouseholdDataManagerMock) CreateHousehold(ctx context.Context, input *types.HouseholdDatabaseCreationInput) (*types.Household, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Household), args.Error(1)
}

// UpdateHousehold is a mock function.
func (m *HouseholdDataManagerMock) UpdateHousehold(ctx context.Context, updated *types.Household) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveHousehold is a mock function.
func (m *HouseholdDataManagerMock) ArchiveHousehold(ctx context.Context, householdID, userID string) error {
	return m.Called(ctx, householdID, userID).Error(0)
}
