package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientDataManager = (*ValidIngredientDataManagerMock)(nil)

// ValidIngredientDataManagerMock is a mocked types.ValidIngredientDataManager for testing.
type ValidIngredientDataManagerMock struct {
	mock.Mock
}

// ValidIngredientExists is a mock function.
func (m *ValidIngredientDataManagerMock) ValidIngredientExists(ctx context.Context, validIngredientID string) (bool, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredient is a mock function.
func (m *ValidIngredientDataManagerMock) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	args := m.Called(ctx, validIngredientID)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// GetRandomValidIngredient is a mock function.
func (m *ValidIngredientDataManagerMock) GetRandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// SearchForValidIngredients is a mock function.
func (m *ValidIngredientDataManagerMock) SearchForValidIngredients(ctx context.Context, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	args := m.Called(ctx, query, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredient]), args.Error(1)
}

// SearchForValidIngredientsForPreparation is a mock function.
func (m *ValidIngredientDataManagerMock) SearchForValidIngredientsForPreparation(ctx context.Context, preparationID, query string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	args := m.Called(ctx, preparationID, query, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredient]), args.Error(1)
}

// GetValidIngredients is a mock function.
func (m *ValidIngredientDataManagerMock) GetValidIngredients(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidIngredient], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidIngredient]), args.Error(1)
}

// CreateValidIngredient is a mock function.
func (m *ValidIngredientDataManagerMock) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientDatabaseCreationInput) (*types.ValidIngredient, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredient), args.Error(1)
}

// UpdateValidIngredient is a mock function.
func (m *ValidIngredientDataManagerMock) UpdateValidIngredient(ctx context.Context, updated *types.ValidIngredient) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredient is a mock function.
func (m *ValidIngredientDataManagerMock) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}

// MarkValidIngredientAsIndexed is a mock function.
func (m *ValidIngredientDataManagerMock) MarkValidIngredientAsIndexed(ctx context.Context, validIngredientID string) error {
	return m.Called(ctx, validIngredientID).Error(0)
}

// GetValidIngredientIDsThatNeedSearchIndexing is a mock function.
func (m *ValidIngredientDataManagerMock) GetValidIngredientIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidIngredientsWithIDs is a mock function.
func (m *ValidIngredientDataManagerMock) GetValidIngredientsWithIDs(ctx context.Context, ids []string) ([]*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.ValidIngredient), returnValues.Error(1)
}
