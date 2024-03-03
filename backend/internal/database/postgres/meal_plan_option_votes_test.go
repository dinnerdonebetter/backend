package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createMealPlanOptionVoteForTest(t *testing.T, ctx context.Context, mealPlanID, mealPlanEventID string, exampleMealPlanOptionVote *types.MealPlanOptionVote, dbc *Querier) *types.MealPlanOptionVote {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput(exampleMealPlanOptionVote)

	rawCreated, err := dbc.CreateMealPlanOptionVote(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, rawCreated)
	assert.Len(t, rawCreated, 1)
	created := rawCreated[0]

	exampleMealPlanOptionVote.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleMealPlanOptionVote, created)

	mealPlanOptionVote, err := dbc.GetMealPlanOptionVote(ctx, mealPlanID, mealPlanEventID, created.BelongsToMealPlanOption, created.ID)
	require.NoError(t, err)

	exampleMealPlanOptionVote.CreatedAt = mealPlanOptionVote.CreatedAt
	require.Equal(t, exampleMealPlanOptionVote.CreatedAt, mealPlanOptionVote.CreatedAt)
	require.Equal(t, exampleMealPlanOptionVote.LastUpdatedAt, mealPlanOptionVote.LastUpdatedAt)
	require.Equal(t, exampleMealPlanOptionVote.ID, mealPlanOptionVote.ID)

	assert.Equal(t, exampleMealPlanOptionVote, mealPlanOptionVote)

	return mealPlanOptionVote
}

func TestQuerier_Integration_MealPlanOptionVotes(t *testing.T) {
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
	buildMealForIntegrationTest(user.ID, recipe)
	meal := createMealForTest(t, ctx, nil, dbc)

	exampleMealPlan := buildMealPlanForIntegrationTest(user.ID, meal)
	exampleMealPlan.BelongsToHousehold = householdID
	mealPlan := createMealPlanForTest(t, ctx, exampleMealPlan, dbc)
	mealPlanEvent := mealPlan.Events[0]
	mealPlanOption := mealPlanEvent.Options[0]

	exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
	exampleMealPlanOptionVote.ByUser = user.ID
	exampleMealPlanOptionVote.BelongsToMealPlanOption = mealPlanOption.ID

	// create
	createdMealPlanOptionVotes := []*types.MealPlanOptionVote{}
	createdMealPlanOptionVotes = append(createdMealPlanOptionVotes, createMealPlanOptionVoteForTest(t, ctx, mealPlan.ID, mealPlanEvent.ID, exampleMealPlanOptionVote, dbc))

	// fetch as list
	mealPlanOptionVotes, err := dbc.GetMealPlanOptionVotes(ctx, mealPlan.ID, mealPlanEvent.ID, mealPlanOption.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanOptionVotes)
	assert.Equal(t, len(createdMealPlanOptionVotes), len(mealPlanOptionVotes.Data))

	assert.NoError(t, dbc.UpdateMealPlanOptionVote(ctx, createdMealPlanOptionVotes[0]))

	// delete
	for _, mealPlanOptionVote := range createdMealPlanOptionVotes {
		assert.NoError(t, dbc.ArchiveMealPlanOptionVote(ctx, mealPlan.ID, mealPlanEvent.ID, mealPlanOption.ID, mealPlanOptionVote.ID))

		var exists bool
		exists, err = dbc.MealPlanOptionVoteExists(ctx, mealPlan.ID, mealPlanEvent.ID, mealPlanOption.ID, mealPlanOptionVote.ID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestQuerier_MealPlanOptionVoteExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVotes(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanOptionVote(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error creating transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
		exampleMealPlanOptionVote.ID = "1"
		exampleInput := converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput(exampleMealPlanOptionVote)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMealPlanOptionVote(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanOptionVote(ctx, nil))
	})
}

func TestQuerier_ArchiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID))
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, ""))
	})
}
