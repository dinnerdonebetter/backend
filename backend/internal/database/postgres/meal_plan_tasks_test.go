package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createMealPlanTaskForTest(t *testing.T, ctx context.Context, exampleMealPlanTask *types.MealPlanTask, dbc *Querier) *types.MealPlanTask {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanTaskToMealPlanTaskDatabaseCreationInput(exampleMealPlanTask)

	created, err := dbc.CreateMealPlanTask(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleMealPlanTask.CreatedAt = created.CreatedAt
	require.Equal(t, exampleMealPlanTask.RecipePrepTask.ID, created.RecipePrepTask.ID)
	exampleMealPlanTask.RecipePrepTask = created.RecipePrepTask
	require.Equal(t, exampleMealPlanTask.MealPlanOption.ID, created.MealPlanOption.ID)
	exampleMealPlanTask.MealPlanOption = created.MealPlanOption
	assert.Equal(t, exampleMealPlanTask, created)

	mealPlanTask, err := dbc.GetMealPlanTask(ctx, created.ID)
	require.NoError(t, err)

	exampleMealPlanTask.CreatedAt = mealPlanTask.CreatedAt
	exampleMealPlanTask.RecipePrepTask = mealPlanTask.RecipePrepTask
	exampleMealPlanTask.MealPlanOption = mealPlanTask.MealPlanOption
	require.Equal(t, exampleMealPlanTask.CreatedAt, mealPlanTask.CreatedAt)
	require.Equal(t, exampleMealPlanTask.LastUpdatedAt, mealPlanTask.LastUpdatedAt)
	require.Equal(t, exampleMealPlanTask.ID, mealPlanTask.ID)
	require.Equal(t, exampleMealPlanTask.Status, mealPlanTask.Status)

	assert.Equal(t, exampleMealPlanTask, mealPlanTask)

	return mealPlanTask
}

func TestQuerier_Integration_MealPlanTasks(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	recipe := createRecipeForTest(t, ctx, nil, dbc, true)
	meal := createMealForTest(t, ctx, buildMealForIntegrationTest(user.ID, recipe), dbc)

	exampleMealPlan := buildMealPlanForIntegrationTest(user.ID, meal)
	exampleMealPlan.BelongsToHousehold = householdID
	mealPlan := createMealPlanForTest(t, ctx, exampleMealPlan, dbc)

	exampleMealPlanTask := fakes.BuildFakeMealPlanTask()
	exampleMealPlanTask.RecipePrepTask = *recipe.PrepTasks[0]
	exampleMealPlanTask.MealPlanOption = *mealPlan.Events[0].Options[0]

	// create
	createdMealPlanTasks := []*types.MealPlanTask{}
	createdMealPlanTasks = append(createdMealPlanTasks, createMealPlanTaskForTest(t, ctx, exampleMealPlanTask, dbc))

	// fetch as list
	mealPlanTasks, err := dbc.GetMealPlanTasksForMealPlan(ctx, mealPlan.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanTasks)
	assert.Equal(t, len(createdMealPlanTasks), len(mealPlanTasks))

	// delete
	for _, mealPlanTask := range createdMealPlanTasks {
		assert.NoError(t, dbc.ChangeMealPlanTaskStatus(ctx, &types.MealPlanTaskStatusChangeRequestInput{
			Status:            pointer.To(types.MealPlanTaskStatusFinished),
			StatusExplanation: t.Name(),
			AssignedToUser:    &user.ID,
			ID:                mealPlanTask.ID,
		}))

		var exists bool
		exists, err = dbc.MealPlanTaskExists(ctx, mealPlanTask.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestQuerier_MealPlanTaskExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanTaskID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.MealPlanTaskExists(ctx, "", exampleMealPlanTaskID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.MealPlanTaskExists(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanTaskID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetMealPlanTask(ctx, exampleMealPlanTaskID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan task ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetMealPlanTask(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_createMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		actual, err := c.createMealPlanTask(ctx, tx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanTask(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.CreateMealPlanTask(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanTasksForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with missing meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetMealPlanTasksForMealPlan(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkMealPlanAsHavingTasksCreated(T *testing.T) {
	T.Parallel()

	T.Run("with empty meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.MarkMealPlanAsHavingTasksCreated(ctx, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkMealPlanAsHavingGroceryListInitialized(T *testing.T) {
	T.Parallel()

	T.Run("with empty meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.MarkMealPlanAsHavingGroceryListInitialized(ctx, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ChangeMealPlanTaskStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ChangeMealPlanTaskStatus(ctx, nil))
	})
}
