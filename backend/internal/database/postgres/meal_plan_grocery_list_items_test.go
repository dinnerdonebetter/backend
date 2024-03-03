package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createMealPlanGroceryListItemForTest(t *testing.T, ctx context.Context, exampleMealPlanGroceryListItem *types.MealPlanGroceryListItem, dbc *Querier) *types.MealPlanGroceryListItem {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)

	created, err := dbc.CreateMealPlanGroceryListItem(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleMealPlanGroceryListItem.CreatedAt = created.CreatedAt
	require.Equal(t, exampleMealPlanGroceryListItem.MeasurementUnit.ID, created.MeasurementUnit.ID)
	exampleMealPlanGroceryListItem.MeasurementUnit = created.MeasurementUnit
	require.Equal(t, exampleMealPlanGroceryListItem.PurchasedMeasurementUnit.ID, created.PurchasedMeasurementUnit.ID)
	exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = created.PurchasedMeasurementUnit
	require.Equal(t, exampleMealPlanGroceryListItem.Ingredient.ID, created.Ingredient.ID)
	exampleMealPlanGroceryListItem.Ingredient = created.Ingredient
	assert.Equal(t, exampleMealPlanGroceryListItem, created)

	mealPlanGroceryListItem, err := dbc.GetMealPlanGroceryListItem(ctx, created.BelongsToMealPlan, created.ID)
	require.NoError(t, err)

	exampleMealPlanGroceryListItem.CreatedAt = mealPlanGroceryListItem.CreatedAt
	require.Equal(t, exampleMealPlanGroceryListItem.MeasurementUnit.ID, mealPlanGroceryListItem.MeasurementUnit.ID)
	exampleMealPlanGroceryListItem.MeasurementUnit = mealPlanGroceryListItem.MeasurementUnit
	require.Equal(t, exampleMealPlanGroceryListItem.PurchasedMeasurementUnit.ID, mealPlanGroceryListItem.PurchasedMeasurementUnit.ID)
	exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = mealPlanGroceryListItem.PurchasedMeasurementUnit
	require.Equal(t, exampleMealPlanGroceryListItem.Ingredient.ID, mealPlanGroceryListItem.Ingredient.ID)
	exampleMealPlanGroceryListItem.Ingredient = mealPlanGroceryListItem.Ingredient
	require.Equal(t, exampleMealPlanGroceryListItem.CreatedAt, mealPlanGroceryListItem.CreatedAt)
	require.Equal(t, exampleMealPlanGroceryListItem.LastUpdatedAt, mealPlanGroceryListItem.LastUpdatedAt)
	require.Equal(t, exampleMealPlanGroceryListItem.ID, mealPlanGroceryListItem.ID)

	assert.Equal(t, exampleMealPlanGroceryListItem, mealPlanGroceryListItem)

	return mealPlanGroceryListItem
}

func TestQuerier_Integration_MealPlanGroceryListItems(t *testing.T) {
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

	ingredient := createValidIngredientForTest(t, ctx, nil, dbc)
	measurmentUnit := createValidMeasurementUnitForTest(t, ctx, nil, dbc)

	exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
	exampleMealPlanGroceryListItem.BelongsToMealPlan = mealPlan.ID
	exampleMealPlanGroceryListItem.Ingredient = *ingredient
	exampleMealPlanGroceryListItem.MeasurementUnit = *measurmentUnit
	exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = measurmentUnit

	// create
	createdMealPlanGroceryListItems := []*types.MealPlanGroceryListItem{}
	createdMealPlanGroceryListItems = append(createdMealPlanGroceryListItems, createMealPlanGroceryListItemForTest(t, ctx, exampleMealPlanGroceryListItem, dbc))

	// update
	assert.NoError(t, dbc.UpdateMealPlanGroceryListItem(ctx, createdMealPlanGroceryListItems[0]))

	// fetch as list
	mealPlanGroceryListItems, err := dbc.GetMealPlanGroceryListItemsForMealPlan(ctx, mealPlan.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanGroceryListItems)
	assert.Equal(t, len(createdMealPlanGroceryListItems), len(mealPlanGroceryListItems))

	// delete
	for _, mealPlanGroceryListItem := range createdMealPlanGroceryListItems {
		assert.NoError(t, dbc.ArchiveMealPlanGroceryListItem(ctx, mealPlanGroceryListItem.ID))

		var exists bool
		exists, err = dbc.MealPlanGroceryListItemExists(ctx, mealPlanGroceryListItem.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}

func TestQuerier_MealPlanGroceryListItemExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		c, _ := buildTestClient(t)

		actual, err := c.MealPlanGroceryListItemExists(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_fleshOutMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.fleshOutMealPlanGroceryListItem(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanGroceryListItem(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_createMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		actual, err := c.createMealPlanGroceryListItem(ctx, tx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanGroceryListItem(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanGroceryListItem(ctx, nil))
	})
}

func TestQuerier_ArchiveMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanGroceryListItem(ctx, ""))
	})
}
