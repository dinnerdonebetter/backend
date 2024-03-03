package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdUserMembershipDataManager = (*HouseholdUserMembershipDataManagerMock)(nil)

// HouseholdUserMembershipDataManagerMock is a mocked types.HouseholdUserMembershipDataManager for testing.
type HouseholdUserMembershipDataManagerMock struct {
	mock.Mock
}

// BuildSessionContextDataForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManagerMock) BuildSessionContextDataForUser(ctx context.Context, userID string) (*types.SessionContextData, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(*types.SessionContextData), returnValues.Error(1)
}

// GetDefaultHouseholdIDForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManagerMock) GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(string), returnValues.Error(1)
}

// MarkHouseholdAsUserDefault implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID string) error {
	return m.Called(ctx, userID, householdID).Error(0)
}

// UserIsMemberOfHousehold implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) UserIsMemberOfHousehold(ctx context.Context, userID, householdID string) (bool, error) {
	returnValues := m.Called(ctx, userID, householdID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// AddUserToHousehold implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) AddUserToHousehold(ctx context.Context, input *types.HouseholdUserMembershipDatabaseCreationInput) error {
	return m.Called(ctx, input).Error(0)
}

// RemoveUserFromHousehold implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) RemoveUserFromHousehold(ctx context.Context, userID, householdID string) error {
	return m.Called(ctx, userID, householdID).Error(0)
}

// ModifyUserPermissions implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) ModifyUserPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	return m.Called(ctx, userID, householdID, input).Error(0)
}

// TransferHouseholdOwnership implements the interface.
func (m *HouseholdUserMembershipDataManagerMock) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	return m.Called(ctx, householdID, input).Error(0)
}
