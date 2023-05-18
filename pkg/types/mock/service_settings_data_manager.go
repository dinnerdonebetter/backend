package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ServiceSettingDataManager = (*ServiceSettingDataManager)(nil)

// ServiceSettingDataManager is a mocked types.ServiceSettingDataManager for testing.
type ServiceSettingDataManager struct {
	mock.Mock
}

// ServiceSettingExists is a mock function.
func (m *ServiceSettingDataManager) ServiceSettingExists(ctx context.Context, serviceSettingID string) (bool, error) {
	args := m.Called(ctx, serviceSettingID)
	return args.Bool(0), args.Error(1)
}

// GetServiceSetting is a mock function.
func (m *ServiceSettingDataManager) GetServiceSetting(ctx context.Context, serviceSettingID string) (*types.ServiceSetting, error) {
	args := m.Called(ctx, serviceSettingID)
	return args.Get(0).(*types.ServiceSetting), args.Error(1)
}

// GetRandomServiceSetting is a mock function.
func (m *ServiceSettingDataManager) GetRandomServiceSetting(ctx context.Context) (*types.ServiceSetting, error) {
	args := m.Called(ctx)
	return args.Get(0).(*types.ServiceSetting), args.Error(1)
}

// SearchForServiceSettings is a mock function.
func (m *ServiceSettingDataManager) SearchForServiceSettings(ctx context.Context, query string) ([]*types.ServiceSetting, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ServiceSetting), args.Error(1)
}

// GetServiceSettings is a mock function.
func (m *ServiceSettingDataManager) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSetting], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ServiceSetting]), args.Error(1)
}

// CreateServiceSetting is a mock function.
func (m *ServiceSettingDataManager) CreateServiceSetting(ctx context.Context, input *types.ServiceSettingDatabaseCreationInput) (*types.ServiceSetting, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ServiceSetting), args.Error(1)
}

// UpdateServiceSetting is a mock function.
func (m *ServiceSettingDataManager) UpdateServiceSetting(ctx context.Context, updated *types.ServiceSetting) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveServiceSetting is a mock function.
func (m *ServiceSettingDataManager) ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error {
	return m.Called(ctx, serviceSettingID).Error(0)
}
