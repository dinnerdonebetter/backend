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

func createValidIngredientMeasurementUnitForTest(t *testing.T, ctx context.Context, exampleValidIngredientMeasurementUnit *types.ValidIngredientMeasurementUnit, dbc *Querier) *types.ValidIngredientMeasurementUnit {
	t.Helper()

	// create
	if exampleValidIngredientMeasurementUnit == nil {
		exampleValidIngredient := createValidIngredientForTest(t, ctx, nil, dbc)
		exampleValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
		exampleValidIngredientMeasurementUnit = fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnit.Ingredient = *exampleValidIngredient
		exampleValidIngredientMeasurementUnit.MeasurementUnit = *exampleValidMeasurementUnit
	}

	dbInput := converters.ConvertValidIngredientMeasurementUnitToValidIngredientMeasurementUnitDatabaseCreationInput(exampleValidIngredientMeasurementUnit)

	created, err := dbc.CreateValidIngredientMeasurementUnit(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientMeasurementUnit.CreatedAt = created.CreatedAt
	exampleValidIngredientMeasurementUnit.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleValidIngredientMeasurementUnit.MeasurementUnit.ID}
	exampleValidIngredientMeasurementUnit.Ingredient = types.ValidIngredient{ID: exampleValidIngredientMeasurementUnit.Ingredient.ID}
	assert.Equal(t, exampleValidIngredientMeasurementUnit, created)

	validIngredientMeasurementUnit, err := dbc.GetValidIngredientMeasurementUnit(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validIngredientMeasurementUnit)
	exampleValidIngredientMeasurementUnit.CreatedAt = validIngredientMeasurementUnit.CreatedAt
	exampleValidIngredientMeasurementUnit.MeasurementUnit = validIngredientMeasurementUnit.MeasurementUnit
	exampleValidIngredientMeasurementUnit.Ingredient = validIngredientMeasurementUnit.Ingredient

	assert.NoError(t, err)
	assert.Equal(t, validIngredientMeasurementUnit, exampleValidIngredientMeasurementUnit)

	return created
}

func TestQuerier_Integration_ValidIngredientMeasurementUnits(t *testing.T) {
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
	exampleValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)
	exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
	exampleValidIngredientMeasurementUnit.MeasurementUnit = *exampleValidMeasurementUnit
	exampleValidIngredientMeasurementUnit.Ingredient = *exampleValidIngredient
	createdValidIngredientMeasurementUnits := []*types.ValidIngredientMeasurementUnit{}

	// create
	createdValidIngredientMeasurementUnits = append(createdValidIngredientMeasurementUnits, createValidIngredientMeasurementUnitForTest(t, ctx, exampleValidIngredientMeasurementUnit, dbc))

	// update
	updatedValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
	updatedValidIngredientMeasurementUnit.ID = createdValidIngredientMeasurementUnits[0].ID
	updatedValidIngredientMeasurementUnit.MeasurementUnit = createdValidIngredientMeasurementUnits[0].MeasurementUnit
	updatedValidIngredientMeasurementUnit.Ingredient = createdValidIngredientMeasurementUnits[0].Ingredient
	assert.NoError(t, dbc.UpdateValidIngredientMeasurementUnit(ctx, updatedValidIngredientMeasurementUnit))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidIngredientMeasurementUnit()
		input.MeasurementUnit = createdValidIngredientMeasurementUnits[0].MeasurementUnit
		input.Ingredient = createdValidIngredientMeasurementUnits[0].Ingredient
		createdValidIngredientMeasurementUnits = append(createdValidIngredientMeasurementUnits, createValidIngredientMeasurementUnitForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validIngredientMeasurementUnits, err := dbc.GetValidIngredientMeasurementUnits(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientMeasurementUnits.Data)
	assert.Equal(t, len(createdValidIngredientMeasurementUnits), len(validIngredientMeasurementUnits.Data))

	forIngredient, err := dbc.GetValidIngredientMeasurementUnitsForIngredient(ctx, createdValidIngredientMeasurementUnits[0].Ingredient.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forIngredient.Data)

	forMeasurementUnit, err := dbc.GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx, createdValidIngredientMeasurementUnits[0].MeasurementUnit.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forMeasurementUnit.Data)

	// delete
	for _, validIngredientMeasurementUnit := range createdValidIngredientMeasurementUnits {
		assert.NoError(t, dbc.ArchiveValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnit.ID))

		var exists bool
		exists, err = dbc.ValidIngredientMeasurementUnitExists(ctx, validIngredientMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientMeasurementUnit
		y, err = dbc.GetValidIngredientMeasurementUnit(ctx, validIngredientMeasurementUnit.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidIngredientMeasurementUnitExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientMeasurementUnitExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientMeasurementUnit(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientMeasurementUnit(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientMeasurementUnit(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientMeasurementUnit(ctx, ""))
	})
}
