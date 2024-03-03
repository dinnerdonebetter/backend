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

func createRecipeStepVesselForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepVessel *types.RecipeStepVessel, dbc *Querier) *types.RecipeStepVessel {
	t.Helper()

	// create
	if exampleRecipeStepVessel == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
		exampleRecipeStep := createdRecipe.Steps[0]

		exampleRecipeStepVessel = fakes.BuildFakeRecipeStepVessel()
		exampleRecipeStepVessel.BelongsToRecipeStep = exampleRecipeStep.ID
	}
	dbInput := converters.ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(exampleRecipeStepVessel)

	created, err := dbc.CreateRecipeStepVessel(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepVessel.CreatedAt = created.CreatedAt
	exampleRecipeStepVessel.Vessel.CreatedAt = created.Vessel.CreatedAt
	exampleRecipeStepVessel.Vessel = created.Vessel
	assert.Equal(t, exampleRecipeStepVessel, created)

	recipeStepVessel, err := dbc.GetRecipeStepVessel(ctx, recipeID, exampleRecipeStepVessel.BelongsToRecipeStep, exampleRecipeStepVessel.ID)

	exampleRecipeStepVessel.CreatedAt = recipeStepVessel.CreatedAt
	exampleRecipeStepVessel.Vessel.CreatedAt = recipeStepVessel.Vessel.CreatedAt
	exampleRecipeStepVessel.Vessel = recipeStepVessel.Vessel

	require.Equal(t, exampleRecipeStepVessel, recipeStepVessel)

	assert.NoError(t, err)
	assert.Equal(t, recipeStepVessel, exampleRecipeStepVessel)

	return created
}

func TestQuerier_Integration_RecipeStepVessels(t *testing.T) {
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

	validVessel := createValidVesselForTest(t, ctx, nil, dbc)
	exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
	exampleRecipeStepVessel.Vessel = validVessel
	exampleRecipeStepVessel.BelongsToRecipeStep = exampleRecipeStep.ID
	createdRecipeStepVessels := []*types.RecipeStepVessel{
		exampleRecipeStep.Vessels[0],
	}

	// create
	createdRecipeStepVessels = append(createdRecipeStepVessels, createRecipeStepVesselForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepVessel, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		validVessel = createValidVesselForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeRecipeStepVessel()
		input.Vessel = validVessel
		input.BelongsToRecipeStep = exampleRecipeStep.ID
		createdRecipeStepVessels = append(createdRecipeStepVessels, createRecipeStepVesselForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeStepVessels, err := dbc.GetRecipeStepVessels(ctx, exampleRecipe.ID, createdRecipeStepVessels[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepVessels.Data)
	assert.Equal(t, len(createdRecipeStepVessels), len(recipeStepVessels.Data))

	// delete
	for _, recipeStepVessel := range createdRecipeStepVessels {
		assert.NoError(t, dbc.ArchiveRecipeStepVessel(ctx, exampleRecipeStep.ID, recipeStepVessel.ID))

		var exists bool
		exists, err = dbc.RecipeStepVesselExists(ctx, exampleRecipe.ID, recipeStepVessel.BelongsToRecipeStep, recipeStepVessel.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepVessel
		y, err = dbc.GetRecipeStepVessel(ctx, exampleRecipe.ID, recipeStepVessel.BelongsToRecipeStep, recipeStepVessel.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_RecipeStepVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, "", exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, exampleRecipeID, "", exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, "", exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, "", exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipeStepVessels(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessels(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessels(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepVessel(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepVessel(ctx, "", exampleRecipeStepVessel.ID))
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepVessel(ctx, exampleRecipeStepID, ""))
	})
}
