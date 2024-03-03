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
	"github.com/stretchr/testify/require"
)

func createValidIngredientForTest(t *testing.T, ctx context.Context, exampleValidIngredient *types.ValidIngredient, dbc *Querier) *types.ValidIngredient {
	t.Helper()

	// create
	if exampleValidIngredient == nil {
		exampleValidIngredient = fakes.BuildFakeValidIngredient()
	}
	dbInput := converters.ConvertValidIngredientToValidIngredientDatabaseCreationInput(exampleValidIngredient)

	created, err := dbc.CreateValidIngredient(ctx, dbInput)
	exampleValidIngredient.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidIngredient, created)

	validIngredient, err := dbc.GetValidIngredient(ctx, created.ID)
	exampleValidIngredient.CreatedAt = validIngredient.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validIngredient, exampleValidIngredient)

	return created
}

func TestQuerier_Integration_ValidIngredients(t *testing.T) {
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

	exampleValidIngredient := fakes.BuildFakeValidIngredient()
	createdValidIngredients := []*types.ValidIngredient{}

	// create
	createdValidIngredients = append(createdValidIngredients, createValidIngredientForTest(t, ctx, exampleValidIngredient, dbc))

	// update
	updatedValidIngredient := fakes.BuildFakeValidIngredient()
	updatedValidIngredient.ID = createdValidIngredients[0].ID
	assert.NoError(t, dbc.UpdateValidIngredient(ctx, updatedValidIngredient))
	createdValidIngredients[0] = updatedValidIngredient

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredient()
		input.Name = fmt.Sprintf("%s %d", updatedValidIngredient.Name, i)
		createdValidIngredients = append(createdValidIngredients, createValidIngredientForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredients, err := dbc.GetValidIngredients(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredients.Data)
	assert.Equal(t, len(createdValidIngredients), len(validIngredients.Data))

	// fetch as list of IDs
	validIngredientIDs := []string{}
	for _, validIngredient := range createdValidIngredients {
		validIngredientIDs = append(validIngredientIDs, validIngredient.ID)
	}

	byIDs, err := dbc.GetValidIngredientsWithIDs(ctx, validIngredientIDs)
	assert.NoError(t, err)
	assert.Equal(t, validIngredients.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidIngredients(ctx, updatedValidIngredient.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, validIngredients.Data, byName.Data)

	random, err := dbc.GetRandomValidIngredient(ctx)
	require.NoError(t, err)
	require.NotNil(t, random)

	needToIndex, err := dbc.GetValidIngredientIDsThatNeedSearchIndexing(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, needToIndex)

	validPreparation := fakes.BuildFakeValidPreparation()
	validPreparation.RestrictToIngredients = false
	preparation := createValidPreparationForTest(t, ctx, validPreparation, dbc)
	validIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	validIngredientPreparation.Ingredient = *createdValidIngredients[0]
	validIngredientPreparation.Preparation = *preparation
	ingredientPrepDBInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationDatabaseCreationInput(validIngredientPreparation)
	createdIngredientPreparation, err := dbc.CreateValidIngredientPreparation(ctx, ingredientPrepDBInput)
	require.NoError(t, err)
	require.NotNil(t, createdIngredientPreparation)
	validIngredientPreparations, err := dbc.SearchForValidIngredientsForPreparation(ctx, preparation.ID, updatedValidIngredient.Name[0:2], nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientPreparations.Data)

	validIngredientStateIngredient := fakes.BuildFakeValidIngredientStateIngredient()
	validIngredientStateIngredient.Ingredient = *createdValidIngredients[0]
	ingredientState := createValidIngredientStateForTest(t, ctx, nil, dbc)
	validIngredientStateIngredient.IngredientState = *ingredientState
	ingredientStateIngredientDBInput := converters.ConvertValidIngredientStateIngredientToValidIngredientStateIngredientDatabaseCreationInput(validIngredientStateIngredient)
	createdIngredientStateIngredient, err := dbc.CreateValidIngredientStateIngredient(ctx, ingredientStateIngredientDBInput)
	require.NoError(t, err)
	require.NotNil(t, createdIngredientStateIngredient)

	// delete
	for _, validIngredient := range createdValidIngredients {
		assert.NoError(t, dbc.MarkValidIngredientAsIndexed(ctx, validIngredient.ID))
		assert.NoError(t, dbc.ArchiveValidIngredient(ctx, validIngredient.ID))

		var exists bool
		exists, err = dbc.ValidIngredientExists(ctx, validIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredient
		y, err = dbc.GetValidIngredient(ctx, validIngredient.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredient(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredients(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredients(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredientsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()
		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidIngredientsForPreparation(ctx, exampleValidPreparation.ID, exampleValidIngredient.Name, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredient(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredient(ctx, ""))
	})
}

func TestQuerier_MarkValidIngredientAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidIngredientAsIndexed(ctx, ""))
	})
}
