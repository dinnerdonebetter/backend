package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.UserIngredientPreferenceDataManager = (*UserIngredientPreferenceDataManagerMock)(nil)

// UserIngredientPreferenceDataManagerMock is a mocked types.UserIngredientPreferenceDataManager for testing.
type UserIngredientPreferenceDataManagerMock struct {
	mock.Mock
}

// UserIngredientPreferenceExists is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) UserIngredientPreferenceExists(ctx context.Context, userIngredientPreferenceID, userID string) (bool, error) {
	args := m.Called(ctx, userIngredientPreferenceID, userID)
	return args.Bool(0), args.Error(1)
}

// GetUserIngredientPreference is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) GetUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) (*types.UserIngredientPreference, error) {
	args := m.Called(ctx, userIngredientPreferenceID, userID)
	return args.Get(0).(*types.UserIngredientPreference), args.Error(1)
}

// GetUserIngredientPreferences is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) GetUserIngredientPreferences(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.UserIngredientPreference], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.UserIngredientPreference]), args.Error(1)
}

// CreateUserIngredientPreference is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) CreateUserIngredientPreference(ctx context.Context, input *types.UserIngredientPreferenceDatabaseCreationInput) ([]*types.UserIngredientPreference, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]*types.UserIngredientPreference), args.Error(1)
}

// UpdateUserIngredientPreference is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) UpdateUserIngredientPreference(ctx context.Context, updated *types.UserIngredientPreference) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveUserIngredientPreference is a mock function.
func (m *UserIngredientPreferenceDataManagerMock) ArchiveUserIngredientPreference(ctx context.Context, userIngredientPreferenceID, userID string) error {
	return m.Called(ctx, userIngredientPreferenceID, userID).Error(0)
}
