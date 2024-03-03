package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientStateIngredientDataManager = (*ValidIngredientStateIngredientDataManagerMock)(nil)

// ValidIngredientStateIngredientDataManagerMock is a mocked types.ValidIngredientStateIngredientDataManager for testing.
type ValidIngredientStateIngredientDataManagerMock struct {
	mock.Mock
}

// ValidIngredientStateIngredientExists is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) ValidIngredientStateIngredientExists(ctx context.Context, validIngredientStateIngredientID string) (bool, error) {
	args := m.Called(ctx, validIngredientStateIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) GetValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	args := m.Called(ctx, validIngredientStateIngredientID)
	return args.Get(0).(*types.ValidIngredientStateIngredient), args.Error(1)
}

// GetValidIngredientStateIngredients is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) GetValidIngredientStateIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// GetValidIngredientStateIngredientsForIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) GetValidIngredientStateIngredientsForIngredient(ctx context.Context, ingredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, ingredientID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// GetValidIngredientStateIngredientsForIngredientState is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) GetValidIngredientStateIngredientsForIngredientState(ctx context.Context, ingredientStateID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredientStateIngredient], error) {
	args := m.Called(ctx, ingredientStateID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredientStateIngredient]), args.Error(1)
}

// CreateValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientDatabaseCreationInput) (*types.ValidIngredientStateIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredientStateIngredient), args.Error(1)
}

// UpdateValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) UpdateValidIngredientStateIngredient(ctx context.Context, updated *types.ValidIngredientStateIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientStateIngredient is a mock function.
func (m *ValidIngredientStateIngredientDataManagerMock) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	return m.Called(ctx, validIngredientStateIngredientID).Error(0)
}

// GetValidIngredientStateIDsThatNeedSearchIndexing is a mock function.
func (m *ValidIngredientStateDataManagerMock) GetValidIngredientStateIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}
