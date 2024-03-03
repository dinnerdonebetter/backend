package workers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockMealPlanFinalizationWorker struct {
	mock.Mock
}

// FinalizeExpiredMealPlans implement our interface.
func (m *MockMealPlanFinalizationWorker) FinalizeExpiredMealPlans(ctx context.Context, payload []byte) (int, error) {
	returnValues := m.Called(ctx, payload)
	return returnValues.Get(0).(int), returnValues.Error(1)
}

// FinalizeExpiredMealPlansWithoutReturningCount implement our interface.
func (m *MockMealPlanFinalizationWorker) FinalizeExpiredMealPlansWithoutReturningCount(ctx context.Context, _ []byte) error {
	return m.Called(ctx).Error(0)
}
