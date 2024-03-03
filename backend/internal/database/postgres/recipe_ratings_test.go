package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRecipeRatingForTest(t *testing.T, ctx context.Context, exampleRecipeRating *types.RecipeRating, dbc *Querier) *types.RecipeRating {
	t.Helper()

	// create
	if exampleRecipeRating == nil {
		exampleRecipeRating = fakes.BuildFakeRecipeRating()
	}
	dbInput := converters.ConvertRecipeRatingToRecipeRatingDatabaseCreationInput(exampleRecipeRating)

	created, err := dbc.CreateRecipeRating(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeRating.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleRecipeRating, created)

	recipeRating, err := dbc.GetRecipeRating(ctx, created.ID)
	exampleRecipeRating.CreatedAt = recipeRating.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, recipeRating, exampleRecipeRating)

	return created
}

func TestQuerier_Integration_RecipeRatings(t *testing.T) {
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

	exampleRecipeRating := fakes.BuildFakeRecipeRating()
	exampleRecipeRating.ByUser = user.ID
	exampleRecipeRating.RecipeID = createdRecipe.ID
	createdRecipeRatings := []*types.RecipeRating{}

	// create
	createdRecipeRatings = append(createdRecipeRatings, createRecipeRatingForTest(t, ctx, exampleRecipeRating, dbc))

	// fetch as list
	recipeRatings, err := dbc.GetRecipeRatings(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeRatings.Data)
	assert.Equal(t, len(createdRecipeRatings), len(recipeRatings.Data))

	// delete
	for _, recipeRating := range createdRecipeRatings {
		assert.NoError(t, dbc.ArchiveRecipeRating(ctx, recipeRating.ID))

		var exists bool
		exists, err = dbc.RecipeRatingExists(ctx, recipeRating.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeRating
		y, err = dbc.GetRecipeRating(ctx, recipeRating.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_RecipeRatingExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.RecipeRatingExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeRating(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeRating(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeRating(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeRating(ctx, ""))
	})
}
