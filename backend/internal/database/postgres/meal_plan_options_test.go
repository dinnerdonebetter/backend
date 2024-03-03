package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resenje.org/schulze"
)

func buildMealPlanOptionForIntegrationTest(meal *types.Meal) *types.MealPlanOption {
	exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
	exampleMealPlanOption.Meal = *meal
	exampleMealPlanOption.MealScale = 1
	exampleMealPlanOption.AssignedCook = nil
	exampleMealPlanOption.AssignedDishwasher = nil
	exampleMealPlanOption.Votes = []*types.MealPlanOptionVote{}

	return exampleMealPlanOption
}

func createMealPlanOptionForTest(t *testing.T, ctx context.Context, mealPlanID string, exampleMealPlanOption *types.MealPlanOption, dbc *Querier) *types.MealPlanOption {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput(exampleMealPlanOption)

	created, err := dbc.CreateMealPlanOption(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleMealPlanOption.CreatedAt = created.CreatedAt
	require.Equal(t, exampleMealPlanOption.Meal.ID, created.Meal.ID)
	exampleMealPlanOption.Meal = created.Meal
	require.Equal(t, len(exampleMealPlanOption.Votes), len(created.Votes))
	exampleMealPlanOption.Votes = created.Votes
	assert.Equal(t, exampleMealPlanOption, created)

	mealPlanOption, err := dbc.GetMealPlanOption(ctx, mealPlanID, created.BelongsToMealPlanEvent, created.ID)
	require.NoError(t, err)

	exampleMealPlanOption.CreatedAt = mealPlanOption.CreatedAt
	require.Equal(t, exampleMealPlanOption.Meal.ID, mealPlanOption.Meal.ID)
	exampleMealPlanOption.Meal = mealPlanOption.Meal
	require.Equal(t, len(exampleMealPlanOption.Votes), len(mealPlanOption.Votes))
	exampleMealPlanOption.Votes = mealPlanOption.Votes
	require.Equal(t, exampleMealPlanOption.CreatedAt, mealPlanOption.CreatedAt)
	require.Equal(t, exampleMealPlanOption.LastUpdatedAt, mealPlanOption.LastUpdatedAt)
	require.Equal(t, exampleMealPlanOption.ID, mealPlanOption.ID)

	assert.Equal(t, exampleMealPlanOption, mealPlanOption)

	return mealPlanOption
}

func TestQuerier_Integration_MealPlanOptions(t *testing.T) {
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

	newMeal := createMealForTest(t, ctx, nil, dbc)
	exampleMealPlanOption := buildMealPlanOptionForIntegrationTest(newMeal)
	exampleMealPlanOption.BelongsToMealPlanEvent = mealPlan.Events[0].ID

	// create
	createdMealPlanOptions := []*types.MealPlanOption{
		mealPlan.Events[0].Options[0],
	}
	createdMealPlanOptions = append(createdMealPlanOptions, createMealPlanOptionForTest(t, ctx, mealPlan.ID, exampleMealPlanOption, dbc))

	// fetch as list
	mealPlanOptions, err := dbc.GetMealPlanOptions(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanOptions)
	assert.Equal(t, len(createdMealPlanOptions), len(mealPlanOptions.Data))

	assert.NoError(t, dbc.UpdateMealPlanOption(ctx, createdMealPlanOptions[0]))

	byID, err := dbc.getMealPlanOptionByID(ctx, createdMealPlanOptions[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, createdMealPlanOptions[0].ID, byID.ID)

	_, err = dbc.FinalizeMealPlanOption(ctx, mealPlan.ID, mealPlan.Events[0].ID, createdMealPlanOptions[0].ID, householdID)
	assert.NoError(t, err)

	// delete
	for _, mealPlanOption := range createdMealPlanOptions {
		assert.NoError(t, dbc.ArchiveMealPlanOption(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, mealPlanOption.ID))

		var exists bool
		exists, err = dbc.MealPlanOptionExists(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, mealPlanOption.ID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestQuerier_MealPlanOptionExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, "", "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, "", "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getMealPlanOptionByID(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getMealPlanOptionByID(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptions(ctx, "", "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanOption(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanOption(ctx, nil))
	})
}

func TestQuerier_ArchiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, "", "", exampleMealPlanOption.ID))
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, ""))
	})
}

func Test_determineWinner(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)

		expected := "blah blah blah"
		exampleWinners := []schulze.Result[string]{
			{
				Choice: t.Name(),
				Wins:   1,
			},
			{
				Choice: "",
				Wins:   2,
			},
			{
				Choice: expected,
				Wins:   3,
			},
		}

		actual := c.determineWinner(exampleWinners)

		assert.Equal(t, expected, actual)
	})

	T.Run("with tie", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)

		expectedA := "blah blah blah"
		expectedB := "beeble beeble"
		exampleWinners := []schulze.Result[string]{
			{
				Choice: expectedA,
				Wins:   3,
			},
			{
				Choice: "",
				Wins:   1,
			},
			{
				Choice: expectedB,
				Wins:   3,
			},
		}

		actual := c.determineWinner(exampleWinners)

		assert.True(t, expectedA == actual || expectedB == actual)
	})
}

func Test_decideOptionWinner(T *testing.T) {
	T.Parallel()

	optionA := "eggs benedict"
	optionB := "scrambled eggs"
	optionC := "buttered toast"
	userID1 := fakes.BuildFakeID()
	userID2 := fakes.BuildFakeID()
	userID3 := fakes.BuildFakeID()
	userID4 := fakes.BuildFakeID()

	T.Run("with clear winner", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		expected := optionA
		exampleOptions := []*types.MealPlanOption{
			{
				ID: optionA,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID1,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    1,
						ByUser:                  userID3,
					},
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    2,
						ByUser:                  userID4,
					},
				},
			},
			{
				ID: optionB,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    0,
						ByUser:                  userID3,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    1,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    1,
						ByUser:                  userID4,
					},
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    2,
						ByUser:                  userID1,
					},
				},
			},
			{
				ID: optionC,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    0,
						ByUser:                  userID4,
					},

					{
						BelongsToMealPlanOption: optionC,
						Rank:                    1,
						ByUser:                  userID1,
					},
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    2,
						ByUser:                  userID2,
					},
					{
						BelongsToMealPlanOption: optionC,
						Rank:                    2,
						ByUser:                  userID3,
					},
				},
			},
		}

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.Equal(t, expected, actual)
		assert.False(t, tiebroken)
		assert.True(t, chosen)
	})

	T.Run("with tie", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		exampleOptions := []*types.MealPlanOption{
			{
				ID: optionA,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID1,
					},
				},
			},
			{
				ID: optionB,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    0,
						ByUser:                  userID2,
					},
				},
			},
		}

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.NotEmpty(t, actual)
		assert.True(t, tiebroken)
		assert.True(t, chosen)
	})

	T.Run("without enough votes", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		exampleOptions := []*types.MealPlanOption{
			{
				ID:    optionA,
				Votes: nil,
			},
		}

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.Empty(t, actual)
		assert.False(t, tiebroken)
		assert.False(t, chosen)
	})
}
