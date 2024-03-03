package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildMealForIntegrationTest(userID string, recipe *types.Recipe) *types.Meal {
	exampleMeal := fakes.BuildFakeMeal()
	exampleMeal.CreatedByUser = userID
	exampleMeal.Components = []*types.MealComponent{
		{
			ComponentType: types.MealComponentTypesMain,
			Recipe:        *recipe,
			RecipeScale:   1,
		},
	}

	return exampleMeal
}

func createMealForTest(t *testing.T, ctx context.Context, exampleMeal *types.Meal, dbc *Querier) *types.Meal {
	t.Helper()

	// create
	if exampleMeal == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		recipe := createRecipeForTest(t, ctx, nil, dbc, true)
		exampleMeal = buildMealForIntegrationTest(user.ID, recipe)
	}
	dbInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

	created, err := dbc.CreateMeal(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	originalComponents := exampleMeal.Components
	exampleMeal.CreatedAt = created.CreatedAt
	exampleMeal.Components = created.Components
	assert.Equal(t, exampleMeal, created)

	meal, err := dbc.GetMeal(ctx, created.ID)
	exampleMeal.CreatedAt = meal.CreatedAt
	exampleMeal.Components = originalComponents

	assert.NoError(t, err)
	assert.Equal(t, meal, exampleMeal)

	return created
}

func TestQuerier_Integration_Meals(t *testing.T) {
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
	recipe := createRecipeForTest(t, ctx, buildRecipeForTestCreation(t, ctx, user.ID, dbc), dbc, false)

	exampleMeal := buildMealForIntegrationTest(user.ID, recipe)
	createdMeals := []*types.Meal{}

	// create
	createdMeals = append(createdMeals, createMealForTest(t, ctx, exampleMeal, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := buildMealForIntegrationTest(user.ID, recipe)
		input.Name = fmt.Sprintf("%s %d", exampleMeal.Name, i)
		createdMeals = append(createdMeals, createMealForTest(t, ctx, input, dbc))
	}

	// fetch as list
	meals, err := dbc.GetMeals(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, meals.Data)
	assert.Equal(t, len(createdMeals), len(meals.Data))

	results, err := dbc.GetMealsWithIDs(ctx, []string{createdMeals[0].ID})
	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	ids, err := dbc.GetMealIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, ids)

	// delete
	for _, meal := range createdMeals {
		assert.NoError(t, dbc.ArchiveMeal(ctx, meal.ID, user.ID))

		var exists bool
		exists, err = dbc.MealExists(ctx, meal.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.Meal
		y, err = dbc.GetMeal(ctx, meal.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_MealExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.MealExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMeal(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.CreateMeal(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with missing meal ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleInput := converters.ConvertMealComponentToMealComponentDatabaseCreationInput(exampleMeal.Components[0])

		err := c.CreateMealComponent(ctx, c.db, "", exampleInput)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.CreateMealComponent(ctx, c.db, exampleMeal.ID, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, exampleMeal.ID, ""))
	})
}

func TestQuerier_MarkMealAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkMealAsIndexed(ctx, ""))
	})
}
