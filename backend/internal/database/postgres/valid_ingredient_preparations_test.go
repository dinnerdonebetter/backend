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

func createValidIngredientPreparationForTest(t *testing.T, ctx context.Context, exampleValidIngredientPreparation *types.ValidIngredientPreparation, dbc *Querier) *types.ValidIngredientPreparation {
	t.Helper()

	// create
	if exampleValidIngredientPreparation == nil {
		exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidIngredientPreparation = fakes.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
		exampleValidIngredientPreparation.Preparation = *exampleValidPreparation
	}

	dbInput := converters.ConvertValidIngredientPreparationToValidIngredientPreparationDatabaseCreationInput(exampleValidIngredientPreparation)

	created, err := dbc.CreateValidIngredientPreparation(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientPreparation.CreatedAt = created.CreatedAt
	exampleValidIngredientPreparation.Preparation = types.ValidPreparation{ID: exampleValidIngredientPreparation.Preparation.ID}
	exampleValidIngredientPreparation.Ingredient = types.ValidIngredient{ID: exampleValidIngredientPreparation.Ingredient.ID}
	assert.Equal(t, exampleValidIngredientPreparation, created)

	validIngredientPreparation, err := dbc.GetValidIngredientPreparation(ctx, created.ID)
	exampleValidIngredientPreparation.CreatedAt = validIngredientPreparation.CreatedAt
	exampleValidIngredientPreparation.Preparation = validIngredientPreparation.Preparation
	exampleValidIngredientPreparation.Ingredient = validIngredientPreparation.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientPreparation, exampleValidIngredientPreparation)

	return created
}

func TestQuerier_Integration_ValidIngredientPreparations(t *testing.T) {
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

	exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
	exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	exampleValidIngredientPreparation.Preparation = *exampleValidPreparation
	exampleValidIngredientPreparation.Ingredient = *exampleValidIngredient
	createdValidIngredientPreparations := []*types.ValidIngredientPreparation{}

	// create
	createdValidIngredientPreparations = append(createdValidIngredientPreparations, createValidIngredientPreparationForTest(t, ctx, exampleValidIngredientPreparation, dbc))

	// update
	updatedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()
	updatedValidIngredientPreparation.ID = createdValidIngredientPreparations[0].ID
	updatedValidIngredientPreparation.Preparation = createdValidIngredientPreparations[0].Preparation
	updatedValidIngredientPreparation.Ingredient = createdValidIngredientPreparations[0].Ingredient
	assert.NoError(t, dbc.UpdateValidIngredientPreparation(ctx, updatedValidIngredientPreparation))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredientPreparation()
		input.Preparation = createdValidIngredientPreparations[0].Preparation
		input.Ingredient = createdValidIngredientPreparations[0].Ingredient
		createdValidIngredientPreparations = append(createdValidIngredientPreparations, createValidIngredientPreparationForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientPreparations, err := dbc.GetValidIngredientPreparations(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientPreparations.Data)
	assert.Equal(t, len(createdValidIngredientPreparations), len(validIngredientPreparations.Data))

	forIngredient, err := dbc.GetValidIngredientPreparationsForIngredient(ctx, createdValidIngredientPreparations[0].Ingredient.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forIngredient.Data)

	forPreparation, err := dbc.GetValidIngredientPreparationsForPreparation(ctx, createdValidIngredientPreparations[0].Preparation.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forPreparation.Data)

	// delete
	for _, validIngredientPreparation := range createdValidIngredientPreparations {
		assert.NoError(t, dbc.ArchiveValidIngredientPreparation(ctx, validIngredientPreparation.ID))

		var exists bool
		exists, err = dbc.ValidIngredientPreparationExists(ctx, validIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientPreparation
		y, err = dbc.GetValidIngredientPreparation(ctx, validIngredientPreparation.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientPreparationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparationsForIngredient(ctx, exampleValidIngredient.ID, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidIngredientPreparationsForPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientPreparationsForPreparation(ctx, exampleValidPreparation.ID, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientPreparation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientPreparation(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientPreparation(ctx, ""))
	})
}
