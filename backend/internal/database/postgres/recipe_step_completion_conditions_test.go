package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRecipeStepCompletionConditionForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepCompletionCondition *types.RecipeStepCompletionCondition, dbc *Querier) *types.RecipeStepCompletionCondition {
	t.Helper()

	// create
	if exampleRecipeStepCompletionCondition == nil {
		t.Fatal("exampleRecipeStepCompletionCondition must not be nil")
	}

	recipe, getRecipeErr := dbc.GetRecipe(ctx, recipeID)
	require.NoError(t, getRecipeErr)
	require.NotNil(t, recipe)

	recipeStep, getRecipeStepErr := dbc.GetRecipeStep(ctx, recipeID, exampleRecipeStepCompletionCondition.BelongsToRecipeStep)
	require.NoError(t, getRecipeStepErr)
	require.NotNil(t, recipeStep)

	ingredientState, getIngredientStateErr := dbc.GetValidIngredientState(ctx, exampleRecipeStepCompletionCondition.IngredientState.ID)
	require.NoError(t, getIngredientStateErr)
	require.NotNil(t, ingredientState)

	dbInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(exampleRecipeStepCompletionCondition)
	dbInput.Ingredients = []*types.RecipeStepCompletionConditionIngredientDatabaseCreationInput{
		{
			ID:                                     identifiers.New(),
			BelongsToRecipeStepCompletionCondition: dbInput.ID,
			RecipeStepIngredient:                   recipe.Steps[0].Ingredients[0].ID,
		},
	}

	created, err := dbc.CreateRecipeStepCompletionCondition(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepCompletionCondition.CreatedAt = created.CreatedAt
	exampleRecipeStepCompletionCondition.IngredientState = created.IngredientState
	exampleRecipeStepCompletionCondition.Ingredients = created.Ingredients
	assert.Equal(t, exampleRecipeStepCompletionCondition, created)

	recipeStepCompletionCondition, err := dbc.GetRecipeStepCompletionCondition(ctx, recipe.ID, recipeStep.ID, created.ID)
	require.NoError(t, err)
	require.NotNil(t, recipeStepCompletionCondition)

	exampleRecipeStepCompletionCondition.CreatedAt = recipeStepCompletionCondition.CreatedAt
	exampleRecipeStepCompletionCondition.IngredientState = recipeStepCompletionCondition.IngredientState
	exampleRecipeStepCompletionCondition.Ingredients = recipeStepCompletionCondition.Ingredients

	require.Equal(t, exampleRecipeStepCompletionCondition, recipeStepCompletionCondition)

	return recipeStepCompletionCondition
}

func TestQuerier_Integration_RecipeStepCompletionConditions(t *testing.T) {
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

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
	exampleRecipeStep := createdRecipe.Steps[0]

	ingredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)
	exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
	exampleRecipeStepCompletionCondition.IngredientState = *ingredientState
	exampleRecipeStepCompletionCondition.BelongsToRecipeStep = exampleRecipeStep.ID
	exampleRecipeStepCompletionCondition.Ingredients = []*types.RecipeStepCompletionConditionIngredient{
		{
			ID:                                     identifiers.New(),
			BelongsToRecipeStepCompletionCondition: exampleRecipeStepCompletionCondition.ID,
			RecipeStepIngredient:                   exampleRecipeStep.Ingredients[0].ID,
		},
	}

	// create
	createdRecipeStepCompletionConditions := []*types.RecipeStepCompletionCondition{}
	createdRecipeStepCompletionConditions = append(createdRecipeStepCompletionConditions, createRecipeStepCompletionConditionForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepCompletionCondition, dbc))

	// fetch as list
	recipeStepCompletionConditions, err := dbc.GetRecipeStepCompletionConditions(ctx, exampleRecipe.ID, createdRecipeStepCompletionConditions[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepCompletionConditions.Data)
	assert.Equal(t, len(createdRecipeStepCompletionConditions), len(recipeStepCompletionConditions.Data))

	// delete
	for _, recipeStepCompletionCondition := range createdRecipeStepCompletionConditions {
		assert.NoError(t, dbc.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStep.ID, recipeStepCompletionCondition.ID))

		var exists bool
		exists, err = dbc.RecipeStepCompletionConditionExists(ctx, exampleRecipe.ID, recipeStepCompletionCondition.BelongsToRecipeStep, recipeStepCompletionCondition.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepCompletionCondition
		y, err = dbc.GetRecipeStepCompletionCondition(ctx, exampleRecipe.ID, recipeStepCompletionCondition.BelongsToRecipeStep, recipeStepCompletionCondition.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_RecipeStepCompletionConditionExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, "", exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, "", exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, "", exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, "", exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionConditions(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepCompletionCondition(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createRecipeStepCompletionCondition(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepCompletionCondition(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCompletionCondition(ctx, "", exampleRecipeStepCompletionCondition.ID))
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStepID, ""))
	})
}
