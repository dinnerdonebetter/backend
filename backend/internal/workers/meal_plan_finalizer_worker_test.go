package workers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestChoresWorker(t *testing.T) *mealPlanFinalizationWorker {
	t.Helper()

	worker := ProvideMealPlanFinalizationWorker(
		logging.NewNoopLogger(),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		tracing.NewNoopTracerProvider(),
	)
	assert.NotNil(t, worker)

	return worker.(*mealPlanFinalizationWorker)
}

func TestProvideChoresWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := ProvideMealPlanFinalizationWorker(
			logging.NewNoopLogger(),
			&database.MockDatabase{},
			&mockpublishers.Publisher{},
			tracing.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
	})
}

func TestChoresWorker_FinalizeExpiredMealPlansWithoutReturningCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleInput := &types.ChoreMessage{
			ChoreType: types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType,
		}
		body, err := json.Marshal(exampleInput)
		require.NoError(t, err)

		exampleMealPlans := fakes.BuildFakeMealPlanList().Data

		dbm := database.NewMockDatabase()
		dbm.MealPlanDataManagerMock.On(
			"GetUnfinalizedMealPlansWithExpiredVotingPeriods",
			testutils.ContextMatcher,
		).Return(exampleMealPlans, nil)

		mqm := &mockpublishers.Publisher{}

		for _, mealPlan := range exampleMealPlans {
			dbm.MealPlanDataManagerMock.On(
				"AttemptToFinalizeMealPlan",
				testutils.ContextMatcher,
				mealPlan.ID,
				mealPlan.BelongsToHousehold,
			).Return(true, nil)
		}

		worker := newTestChoresWorker(t)
		worker.dataManager = dbm
		worker.postUpdatesPublisher = mqm

		assert.NoError(t, worker.FinalizeExpiredMealPlansWithoutReturningCount(ctx, body))

		mock.AssertExpectationsForObjects(t, dbm, mqm)
	})
}
