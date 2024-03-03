package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidInstrumentDataManager = (*ValidInstrumentDataManagerMock)(nil)

// ValidInstrumentDataManagerMock is a mocked types.ValidInstrumentDataManager for testing.
type ValidInstrumentDataManagerMock struct {
	mock.Mock
}

// ValidInstrumentExists is a mock function.
func (m *ValidInstrumentDataManagerMock) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (bool, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Bool(0), args.Error(1)
}

// GetValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	args := m.Called(ctx, validInstrumentID)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// GetRandomValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// SearchForValidInstruments is a mock function.
func (m *ValidInstrumentDataManagerMock) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidInstrument), args.Error(1)
}

// GetValidInstruments is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstruments(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidInstrument], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidInstrument]), args.Error(1)
}

// CreateValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidInstrument), args.Error(1)
}

// UpdateValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidInstrument is a mock function.
func (m *ValidInstrumentDataManagerMock) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// MarkValidInstrumentAsIndexed is a mock function.
func (m *ValidInstrumentDataManagerMock) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	return m.Called(ctx, validInstrumentID).Error(0)
}

// GetValidInstrumentIDsThatNeedSearchIndexing is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidInstrumentsWithIDs is a mock function.
func (m *ValidInstrumentDataManagerMock) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.ValidInstrument), returnValues.Error(1)
}
