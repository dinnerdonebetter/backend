package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientGroupDataManager = (*ValidIngredientGroupDataManagerMock)(nil)

// ValidIngredientGroupDataManagerMock is a mocked types.ValidIngredientGroupDataManager for testing.
type ValidIngredientGroupDataManagerMock struct {
	mock.Mock
}

// ValidIngredientGroupExists is a mock method.
func (m *ValidIngredientGroupDataManagerMock) ValidIngredientGroupExists(ctx context.Context, validIngredientID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManagerMock) GetValidIngredientGroup(ctx context.Context, validIngredientID string) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

// GetValidIngredientGroups is a mock method.
func (m *ValidIngredientGroupDataManagerMock) GetValidIngredientGroups(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientGroup], error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).(*types.QueryFilteredResult[types.ValidIngredientGroup]), returnValues.Error(1)
}

// SearchForValidIngredientGroups is a mock method.
func (m *ValidIngredientGroupDataManagerMock) SearchForValidIngredientGroups(ctx context.Context, query string, filter *types.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, query, filter)

	return returnValues.Get(0).([]*types.ValidIngredientGroup), returnValues.Error(1)
}

// CreateValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManagerMock) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupDatabaseCreationInput) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

// UpdateValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManagerMock) UpdateValidIngredientGroup(ctx context.Context, updated *types.ValidIngredientGroup) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientGroup is a mock method.
func (m *ValidIngredientGroupDataManagerMock) ArchiveValidIngredientGroup(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}
