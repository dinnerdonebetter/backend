package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationDataManager = (*ValidPreparationDataManagerMock)(nil)

// ValidPreparationDataManagerMock is a mocked types.ValidPreparationDataManager for testing.
type ValidPreparationDataManagerMock struct {
	mock.Mock
}

// ValidPreparationExists is a mock function.
func (m *ValidPreparationDataManagerMock) ValidPreparationExists(ctx context.Context, validPreparationID string) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparation is a mock function.
func (m *ValidPreparationDataManagerMock) GetValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// GetRandomValidPreparation is a mock function.
func (m *ValidPreparationDataManagerMock) GetRandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// SearchForValidPreparations is a mock function.
func (m *ValidPreparationDataManagerMock) SearchForValidPreparations(ctx context.Context, query string) ([]*types.ValidPreparation, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidPreparation), args.Error(1)
}

// GetValidPreparations is a mock function.
func (m *ValidPreparationDataManagerMock) GetValidPreparations(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparation], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidPreparation]), args.Error(1)
}

// CreateValidPreparation is a mock function.
func (m *ValidPreparationDataManagerMock) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationDatabaseCreationInput) (*types.ValidPreparation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidPreparation), args.Error(1)
}

// UpdateValidPreparation is a mock function.
func (m *ValidPreparationDataManagerMock) UpdateValidPreparation(ctx context.Context, updated *types.ValidPreparation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparation is a mock function.
func (m *ValidPreparationDataManagerMock) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}

// MarkValidPreparationAsIndexed is a mock function.
func (m *ValidPreparationDataManagerMock) MarkValidPreparationAsIndexed(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}

// GetValidPreparationIDsThatNeedSearchIndexing is a mock function.
func (m *ValidPreparationDataManagerMock) GetValidPreparationIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidPreparationsWithIDs is a mock function.
func (m *ValidPreparationDataManagerMock) GetValidPreparationsWithIDs(ctx context.Context, ids []string) ([]*types.ValidPreparation, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*types.ValidPreparation), args.Error(1)
}
