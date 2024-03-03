package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientStateDataManager = (*ValidIngredientStateDataManagerMock)(nil)

// ValidIngredientStateDataManagerMock is a mocked types.ValidIngredientStateDataManager for testing.
type ValidIngredientStateDataManagerMock struct {
	mock.Mock
}

// ValidIngredientStateExists is a mock function.
func (m *ValidIngredientStateDataManagerMock) ValidIngredientStateExists(ctx context.Context, validIngredientStateID string) (bool, error) {
	args := m.Called(ctx, validIngredientStateID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientState is a mock function.
func (m *ValidIngredientStateDataManagerMock) GetValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	args := m.Called(ctx, validIngredientStateID)
	return args.Get(0).(*types.ValidIngredientState), args.Error(1)
}

// SearchForValidIngredientStates is a mock function.
func (m *ValidIngredientStateDataManagerMock) SearchForValidIngredientStates(ctx context.Context, query string) ([]*types.ValidIngredientState, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidIngredientState), args.Error(1)
}

// GetValidIngredientStates is a mock function.
func (m *ValidIngredientStateDataManagerMock) GetValidIngredientStates(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientState], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientState]), args.Error(1)
}

// CreateValidIngredientState is a mock function.
func (m *ValidIngredientStateDataManagerMock) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateDatabaseCreationInput) (*types.ValidIngredientState, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredientState), args.Error(1)
}

// UpdateValidIngredientState is a mock function.
func (m *ValidIngredientStateDataManagerMock) UpdateValidIngredientState(ctx context.Context, updated *types.ValidIngredientState) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientState is a mock function.
func (m *ValidIngredientStateDataManagerMock) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	return m.Called(ctx, validIngredientStateID).Error(0)
}

// MarkValidIngredientStateAsIndexed is a mock function.
func (m *ValidIngredientStateDataManagerMock) MarkValidIngredientStateAsIndexed(ctx context.Context, validIngredientStateID string) error {
	return m.Called(ctx, validIngredientStateID).Error(0)
}

// GetValidIngredientStatesWithIDs is a mock function.
func (m *ValidIngredientStateDataManagerMock) GetValidIngredientStatesWithIDs(ctx context.Context, ids []string) ([]*types.ValidIngredientState, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*types.ValidIngredientState), args.Error(1)
}
