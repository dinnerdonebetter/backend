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

func createRecipeStepInstrumentForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepInstrument *types.RecipeStepInstrument, dbc *Querier) *types.RecipeStepInstrument {
	t.Helper()

	// create
	if exampleRecipeStepInstrument == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
		exampleRecipeStep := createdRecipe.Steps[0]

		exampleRecipeStepInstrument = fakes.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
	}
	dbInput := converters.ConvertRecipeStepInstrumentToRecipeStepInstrumentDatabaseCreationInput(exampleRecipeStepInstrument)

	created, err := dbc.CreateRecipeStepInstrument(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepInstrument.CreatedAt = created.CreatedAt
	exampleRecipeStepInstrument.Instrument.CreatedAt = created.Instrument.CreatedAt
	exampleRecipeStepInstrument.Instrument = created.Instrument
	assert.Equal(t, exampleRecipeStepInstrument, created)

	recipeStepInstrument, err := dbc.GetRecipeStepInstrument(ctx, recipeID, exampleRecipeStepInstrument.BelongsToRecipeStep, exampleRecipeStepInstrument.ID)

	exampleRecipeStepInstrument.CreatedAt = recipeStepInstrument.CreatedAt
	exampleRecipeStepInstrument.Instrument.CreatedAt = recipeStepInstrument.Instrument.CreatedAt
	exampleRecipeStepInstrument.Instrument = recipeStepInstrument.Instrument

	require.Equal(t, exampleRecipeStepInstrument, recipeStepInstrument)

	assert.NoError(t, err)
	assert.Equal(t, recipeStepInstrument, exampleRecipeStepInstrument)

	return created
}

func TestQuerier_Integration_RecipeStepInstruments(t *testing.T) {
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

	validInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
	exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
	exampleRecipeStepInstrument.Instrument = validInstrument
	exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
	createdRecipeStepInstruments := []*types.RecipeStepInstrument{
		exampleRecipeStep.Instruments[0],
	}

	// create
	createdRecipeStepInstruments = append(createdRecipeStepInstruments, createRecipeStepInstrumentForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepInstrument, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		validInstrument = createValidInstrumentForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeRecipeStepInstrument()
		input.Instrument = validInstrument
		input.BelongsToRecipeStep = exampleRecipeStep.ID
		createdRecipeStepInstruments = append(createdRecipeStepInstruments, createRecipeStepInstrumentForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeStepInstruments, err := dbc.GetRecipeStepInstruments(ctx, exampleRecipe.ID, createdRecipeStepInstruments[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepInstruments.Data)
	assert.Equal(t, len(createdRecipeStepInstruments), len(recipeStepInstruments.Data))

	// delete
	for _, recipeStepInstrument := range createdRecipeStepInstruments {
		assert.NoError(t, dbc.ArchiveRecipeStepInstrument(ctx, exampleRecipeStep.ID, recipeStepInstrument.ID))

		var exists bool
		exists, err = dbc.RecipeStepInstrumentExists(ctx, exampleRecipe.ID, recipeStepInstrument.BelongsToRecipeStep, recipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepInstrument
		y, err = dbc.GetRecipeStepInstrument(ctx, exampleRecipe.ID, recipeStepInstrument.BelongsToRecipeStep, recipeStepInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_RecipeStepInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, "", exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, "", exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, "", exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, "", exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstruments(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepInstrument(ctx, "", exampleRecipeStepInstrument.ID))
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepInstrument(ctx, exampleRecipeStepID, ""))
	})
}
