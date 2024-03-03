package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

type RecipePrepTaskDataManagerMock struct {
	mock.Mock
}

// RecipePrepTaskExists implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error) {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Bool(0), retVals.Error(1)
}

// GetRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

// GetRecipePrepTasksForRecipe implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, recipeID)

	return retVals.Get(0).([]*types.RecipePrepTask), retVals.Error(1)
}

// CreateRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) CreateRecipePrepTask(ctx context.Context, input *types.RecipePrepTaskDatabaseCreationInput) (*types.RecipePrepTask, error) {
	retVals := m.Called(ctx, input)

	return retVals.Get(0).(*types.RecipePrepTask), retVals.Error(1)
}

// UpdateRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) UpdateRecipePrepTask(ctx context.Context, updated *types.RecipePrepTask) error {
	retVals := m.Called(ctx, updated)

	return retVals.Error(0)
}

// ArchiveRecipePrepTask implements the requisite interface.
func (m *RecipePrepTaskDataManagerMock) ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error {
	retVals := m.Called(ctx, recipeID, recipePrepTaskID)

	return retVals.Error(0)
}
