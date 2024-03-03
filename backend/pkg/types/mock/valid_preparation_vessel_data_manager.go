package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationVesselDataManager = (*ValidPreparationVesselDataManagerMock)(nil)

// ValidPreparationVesselDataManagerMock is a mocked types.ValidPreparationVesselDataManager for testing.
type ValidPreparationVesselDataManagerMock struct {
	mock.Mock
}

// ValidPreparationVesselExists is a mock function.
func (m *ValidPreparationVesselDataManagerMock) ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (bool, error) {
	args := m.Called(ctx, validPreparationVesselID)
	return args.Bool(0), args.Error(1)
}

// GetValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	args := m.Called(ctx, validPreparationVesselID)
	return args.Get(0).(*types.ValidPreparationVessel), args.Error(1)
}

// GetValidPreparationVessels is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVessels(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidPreparationVessel]), args.Error(1)
}

// GetValidPreparationVesselsForPreparation is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	args := m.Called(ctx, preparationID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidPreparationVessel]), args.Error(1)
}

// GetValidPreparationVesselsForVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidPreparationVessel], error) {
	args := m.Called(ctx, vesselID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidPreparationVessel]), args.Error(1)
}

// CreateValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselDatabaseCreationInput) (*types.ValidPreparationVessel, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidPreparationVessel), args.Error(1)
}

// UpdateValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) UpdateValidPreparationVessel(ctx context.Context, updated *types.ValidPreparationVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	return m.Called(ctx, validPreparationVesselID).Error(0)
}
